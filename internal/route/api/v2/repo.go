package v2

import (
	stdcontext "context"
	json "github.com/nef734/ggit/internal/jsonx"
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/gitutil"
	"github.com/nef734/ggit/internal/gitx"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/markup"
	reposervice "github.com/nef734/ggit/internal/services/repo"
	"github.com/nef734/ggit/internal/tool"
)

// SearchRepoContents searches for a keyword in the repository using git grep.
func SearchRepoContents(c *context.APIContext) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []any{})
		return
	}

	ref := c.Query("ref")
	if ref == "" {
		ref = c.Repo.BranchName
		if ref == "" {
			ref = c.Repo.Repository.DefaultBranch
		}
	}

	args := []string{"grep", "-I", "--line-number", "--null", "--max-count", "100"}
	if c.Query("case_sensitive") != "true" {
		args = append(args, "-i")
	}
	if c.Query("regex") == "true" {
		args = append(args, "-E")
	} else {
		args = append(args, "-F")
	}

	// 🚨 SECURITY: Prevent Parameter Injection (e.g. "--open-files-in-pager=...")
	// Force the end of command-line flags so Git interprets 'ref' ONLY as a revision.
	// In 'git grep', the syntax is: git grep [options] [-e pattern] [--] [<rev>...]
	if strings.HasPrefix(ref, "-") {
		c.Error(errors.New("invalid ref format"), "malicious ref detected")
		return
	}
	args = append(args, "-e", q, "--", ref)

	if err := tool.AcquireSemaphore(c.Req.Context(), tool.GitGrepSemaphore); err != nil {
		c.Error(err, "acquire grep semaphore")
		return
	}
	defer tool.ReleaseSemaphore(tool.GitGrepSemaphore)

	gitutil.Acquire()
	output, err := git.NewCommand(args...).RunInDir(c.Repo.Repository.RepoPath())
	gitutil.Release()
	if err != nil {
		// git grep returns exit code 1 if no matches found
		c.JSON(http.StatusOK, []any{})
		return
	}

	stdout := string(output)
	lines := strings.Split(stdout, "\n")
	results := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		// git grep --null --line-number format:
		// ref:filename\x00line:content
		parts := strings.SplitN(line, "\x00", 2)
		if len(parts) < 2 {
			continue
		}

		filePart := parts[0]
		contentPart := parts[1]

		// filePart is "ref:filename"
		fileName := strings.TrimPrefix(filePart, ref+":")

		// contentPart is "line:content"
		contentSubParts := strings.SplitN(contentPart, ":", 2)
		lineNum := 0
		content := ""
		if len(contentSubParts) == 2 {
			lineNum, _ = strconv.Atoi(contentSubParts[0])
			content = contentSubParts[1]
		}

		results = append(results, map[string]any{
			"filename": fileName,
			"line":     lineNum,
			"content":  content,
		})
	}

	c.JSON(http.StatusOK, results)
}

// CreateRepo creates a new repository for the authenticated user.
func CreateRepo(c *context.APIContext) {
	var f struct {
		UID           int64     `json:"uid"`
		Name          string    `json:"name" binding:"required,max=100"`
		Description   string    `json:"description" binding:"max=255"`
		Website       string    `json:"website" binding:"max=100"`
		DefaultBranch string    `json:"default_branch" binding:"max=100"`
		Private       bool      `json:"private"`
		Unlisted      bool      `json:"unlisted"`
		AutoInit      bool      `json:"auto_init"`
		Gitignores    string    `json:"gitignores"`
		License       string    `json:"license"`
		Readme        string    `json:"readme"`
		IsArchived    bool      `json:"is_archived"`
		CreatedAt     time.Time `json:"created_at"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	owner := c.User
	if f.UID != 0 && f.UID != c.User.ID {
		var err error
		du, err := database.Handle.Users().GetByID(c.Req.Context(), f.UID)
		if err != nil {
			if database.IsErrUserNotExist(err) {
				c.StatusError(http.StatusNotFound, "Owner not found.", err)
			} else {
				c.Error(err, "get owner")
			}
			return
		}
		owner = database.FromDomainUser(du)

		// Check permission
		if !c.User.IsAdmin {
			if !owner.IsOrganization() || !owner.IsOwnedBy(c.User.ID) {
				c.Status(http.StatusForbidden)
				return
			}
		}
	}

	repo, err := reposervice.Repos.Create(c.Req.Context(), c.User, owner, reposervice.CreateOptions{
		Name:          f.Name,
		Description:   f.Description,
		Website:       f.Website,
		DefaultBranch: f.DefaultBranch,
		IsPrivate:     f.Private,
		IsUnlisted:    f.Unlisted,
		AutoInit:      f.AutoInit,
		Gitignores:    f.Gitignores,
		License:       f.License,
		Readme:        f.Readme,
		IsArchived:    f.IsArchived,
		CreatedAt:     f.CreatedAt,
	})
	if err != nil {
		c.Error(err, "create repository")
		return
	}

	c.JSON(http.StatusCreated, repo.APIFormat(c.User))
}

// GetRepoView returns the content of a file or directory in the repository.
func GetRepoView(c *context.APIContext) {
	log.Trace("GetRepoView: repo=%v, path=%q", c.Repo.Repository != nil, c.Repo.TreePath)
	if c.Repo == nil || c.Repo.Repository == nil {
		c.Status(http.StatusNotFound)
		return
	}

	if c.Repo.Repository.IsBare {
		c.JSON(http.StatusOK, map[string]any{
			"is_bare":    true,
			"repository": c.Repo.Repository.APIFormat(c.Repo.Owner),
		})
		return
	}

	start := time.Now()
	// For root directory, bypass TreeEntry check
	if c.Repo.TreePath == "" {
		getRepoDirectory(c)
		log.Trace("GetRepoView (root) took %v", time.Since(start))
		return
	}

	if c.Repo.Commit == nil {
		c.Status(http.StatusNoContent)
		return
	}

	entry, err := c.Repo.Commit.TreeEntry(c.Repo.TreePath)
	if err != nil {
		log.Trace("TreeEntry error after %v: %v", time.Since(start), err)
		c.NotFoundOrError(gitx.NewError(err), "get tree entry")
		return
	}

	if entry.IsTree() {
		getRepoDirectory(c)
	} else {
		getRepoFile(c, entry)
	}
	log.Trace("GetRepoView (%s) took %v", c.Repo.TreePath, time.Since(start))
}

func getRepoDirectory(c *context.APIContext) {
	start := time.Now()
	log.Trace("getRepoDirectory: treePath=%q", c.Repo.TreePath)

	if c.Repo.Commit == nil {
		log.Trace("getRepoDirectory: commit is nil, returning no content")
		c.Status(http.StatusNoContent)
		return
	}

	tree, err := c.Repo.Commit.Subtree(c.Repo.TreePath)
	if err != nil {
		log.Trace("getRepoDirectory: Subtree error: %v", err)
		c.NotFoundOrError(gitx.NewError(err), "get subtree")
		return
	}
	log.Trace("getRepoDirectory: Subtree took %v", time.Since(start))

	startEntries := time.Now()
	entries, err := tree.Entries()
	if err != nil {
		log.Trace("getRepoDirectory: Entries error: %v", err)
		c.Error(err, "list entries")
		return
	}
	entries.Sort()
	log.Trace("getRepoDirectory: Entries took %v", time.Since(startEntries))

	type entryData struct {
		Name       string         `json:"name"`
		Type       string         `json:"type"`
		Size       int64          `json:"size"`
		LastCommit map[string]any `json:"last_commit,omitempty"`
	}

	apiEntries := make([]entryData, len(entries))
	entryNames := make([]string, len(entries))
	for i, e := range entries {
		apiEntries[i] = entryData{
			Name: e.Name(),
			Type: string(e.Type()),
			Size: e.Size(),
		}
		entryNames[i] = e.Name()
	}
	log.Trace("Entries in %s: %v", c.Repo.TreePath, entryNames)

	if len(apiEntries) > 0 && c.Repo.Commit != nil {
		// Fetching last commit for each entry can be slow, but it's essential for a good UI.
		// We use a simple goroutine-based approach for small counts to improve speed.
		log.Trace("GetRepoView: Fetching last commits for %d entries in parallel", len(apiEntries))

		type commitResult struct {
			index  int
			commit *git.Commit
		}
		results := make(chan commitResult, len(apiEntries))

		for i := range apiEntries {
			go func(idx int, entryName string) {
				lastCommit, err := c.Repo.GitRepo.Log(c.Repo.Commit.ID.String(), git.LogOptions{
					Path:     path.Join(c.Repo.TreePath, entryName),
					MaxCount: 1,
				})
				if err == nil && len(lastCommit) > 0 {
					results <- commitResult{index: idx, commit: lastCommit[0]}
				} else {
					results <- commitResult{index: idx, commit: nil}
				}
			}(i, apiEntries[i].Name)
		}

		for range apiEntries {
			res := <-results
			if res.commit != nil {
				apiEntries[res.index].LastCommit = map[string]any{
					"id":      res.commit.ID.String(),
					"summary": res.commit.Summary(),
					"author": map[string]any{
						"name":  res.commit.Author.Name,
						"email": res.commit.Author.Email,
						"when":  res.commit.Author.When,
					},
				}
			}
		}
		log.Trace("GetRepoView: Finished fetching parallel commits")
	}

	res := map[string]any{
		"type":       "directory",
		"entries":    apiEntries,
		"path":       c.Repo.TreePath,
		"branch":     c.Repo.BranchName,
		"repository": c.Repo.Repository.APIFormat(c.Repo.Owner),
	}

	if c.Repo.Commit != nil {
		log.Trace("GetRepoView: Fetching last commit for current path: %q", c.Repo.TreePath)
		lastCommit, err := c.Repo.GitRepo.Log(c.Repo.Commit.ID.String(), git.LogOptions{Path: c.Repo.TreePath, MaxCount: 1})
		if err == nil && len(lastCommit) > 0 {
			sha := lastCommit[0].ID.String()
			res["last_commit"] = map[string]any{
				"id":      sha,
				"summary": lastCommit[0].Summary(),
				"status":  string(database.Handle.CommitStatuses().GetCombinedStatusBySHA(c.Req.Context(), c.Repo.Repository.ID, sha)),
				"author": map[string]any{
					"name":  lastCommit[0].Author.Name,
					"email": lastCommit[0].Author.Email,
					"when":  lastCommit[0].Author.When,
				},
			}
		}
	}
	log.Trace("Directory %s has %d entries", c.Repo.TreePath, len(apiEntries))

	// Readme detection
	var readmeFile *git.Blob
	for _, entry := range entries {
		if entry.IsTree() || !markup.IsReadmeFile(entry.Name()) {
			continue
		}

		var err error
		readmeFile, err = tree.Blob(entry.Name())
		if err != nil {
			log.Error("Failed to get README blob: %v", err)
		}
		break
	}

	if readmeFile != nil {
		p, err := readmeFile.Bytes()
		if err == nil && tool.IsTextFile(p) {
			res["readme"] = map[string]any{
				"name":    readmeFile.Name(),
				"content": string(p),
			}
		}
	}

	// Stats and Languages (Only on root)
	if c.Repo.TreePath == "" {
		stats := map[string]int64{}
		stats["commits"] = c.Repo.Repository.NumCommits
		stats["branches"] = int64(c.Repo.Repository.NumBranches)

		tags, _ := c.Repo.GitRepo.Tags()
		stats["tags"] = int64(len(tags))

		var numReleases int64
		_ = database.Handle.DB().Model(&database.Release{}).Where("repo_id = ? AND is_draft = ?", c.Repo.Repository.ID, false).Count(&numReleases).Error
		stats["releases"] = numReleases

		res["stats"] = stats

		// Fetch language stats from database cache
		langStats := make(map[string]int64)
		if c.Repo.Repository.Languages != "" {
			_ = json.Unmarshal([]byte(c.Repo.Repository.Languages), &langStats)
		}
		res["languages"] = langStats
	}

	c.JSON(http.StatusOK, res)
}

func getRepoFile(c *context.APIContext, entry *git.TreeEntry) {
	blob := entry.Blob()
	p, err := blob.Bytes()
	if err != nil {
		c.Error(err, "read blob")
		return
	}

	res := map[string]any{
		"type":         "file",
		"name":         blob.Name(),
		"size":         blob.Size(),
		"is_text_file": tool.IsTextFile(p),
		"path":         c.Repo.TreePath,
		"branch":       c.Repo.BranchName,
		"repository":   c.Repo.Repository.APIFormat(c.Repo.Owner),
	}

	lastCommit, _ := c.Repo.GitRepo.Log(c.Repo.Commit.ID.String(), git.LogOptions{Path: c.Repo.TreePath, MaxCount: 1})
	if len(lastCommit) > 0 {
		sha := lastCommit[0].ID.String()
		res["last_commit"] = map[string]any{
			"id":      sha,
			"summary": lastCommit[0].Summary(),
			"status":  string(database.Handle.CommitStatuses().GetCombinedStatusBySHA(c.Req.Context(), c.Repo.Repository.ID, sha)),
			"author": map[string]any{
				"name":  lastCommit[0].Author.Name,
				"email": lastCommit[0].Author.Email,
				"when":  lastCommit[0].Author.When,
			},
		}
	}

	if tool.IsTextFile(p) {
		res["content"] = string(p)

		if markup.IsReadmeFile(blob.Name()) {
			res["readme_content"] = string(p)
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetRepoBlame returns the blame information for a file.
func GetRepoBlame(c *context.APIContext) {
	if c.Repo.TreePath == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	gitutil.Acquire()
	defer gitutil.Release()

	// Use --porcelain for machine-readable output
	args := []string{"blame", "--porcelain"}
	if c.Repo.CommitID != "" {
		args = append(args, c.Repo.CommitID)
	}
	args = append(args, "--", c.Repo.TreePath)

	stdout, err := git.NewCommand(args...).RunInDir(c.Repo.Repository.RepoPath())
	if err != nil {
		c.Error(err, "git blame")
		return
	}

	lines := strings.Split(string(stdout), "\n")
	type hunk struct {
		Commit string
		Author string
		Time   int64
	}
	hunks := make(map[string]*hunk)
	results := make([]map[string]any, 0, 100)

	for i := 0; i < len(lines); {
		line := lines[i]
		if line == "" {
			i++
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 4 {
			i++
			continue
		}

		sha := parts[0]
		// Each hunk starts with a header line: <sha> <orig_line> <final_line> <num_lines>
		h, ok := hunks[sha]
		if !ok {
			h = &hunk{Commit: sha}
			hunks[sha] = h
			// Parse metadata following the header until we hit a tab-prefixed line (the content)
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "\t") {
				meta := lines[i]
				if strings.HasPrefix(meta, "author ") {
					h.Author = strings.TrimPrefix(meta, "author ")
				} else if strings.HasPrefix(meta, "author-time ") {
					t, _ := strconv.ParseInt(strings.TrimPrefix(meta, "author-time "), 10, 64)
					h.Time = t
				}
				i++
			}
		} else {
			// Skip metadata if we already have it for this SHA
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "\t") {
				i++
			}
		}

		if i < len(lines) && strings.HasPrefix(lines[i], "\t") {
			results = append(results, map[string]any{
				"commit": h.Commit,
				"author": h.Author,
				"time":   h.Time,
				"line":   strings.TrimPrefix(lines[i], "\t"),
			})
			i++
		}
	}

	c.JSON(http.StatusOK, results)
}

// matchUsersWithCommitEmails matches existing users using commit author emails and includes commit statuses.
func matchUsersWithCommitEmails(ctx stdcontext.Context, repoID int64, oldCommits []*git.Commit) []map[string]any {
	emailToUsers := make(map[string]*database.User)
	newCommits := make([]map[string]any, len(oldCommits))
	usersStore := database.Handle.Users()
	statusStore := database.Handle.CommitStatuses()
	for i := range oldCommits {
		var u *database.User
		if v, ok := emailToUsers[oldCommits[i].Author.Email]; !ok {
			du, _ := usersStore.GetByEmail(ctx, oldCommits[i].Author.Email)
			u = database.FromDomainUser(du)
			emailToUsers[oldCommits[i].Author.Email] = u
		} else {
			u = v
		}

		sha := oldCommits[i].ID.String()
		newCommits[i] = map[string]any{
			"user":    u,
			"id":      sha,
			"message": oldCommits[i].Message,
			"author": map[string]any{
				"name":  oldCommits[i].Author.Name,
				"email": oldCommits[i].Author.Email,
				"when":  oldCommits[i].Author.When,
			},
			"committer": map[string]any{
				"name":  oldCommits[i].Committer.Name,
				"email": oldCommits[i].Committer.Email,
				"when":  oldCommits[i].Committer.When,
			},
			"status": string(statusStore.GetCombinedStatusBySHA(ctx, repoID, sha)),
		}
	}
	return newCommits
}

// GetRepoCommits returns a list of commits for the repository.
func GetRepoCommits(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize := c.QueryInt("pageSize")
	if pageSize <= 0 {
		pageSize = 20
	}

	if c.Repo.Commit == nil {
		c.JSON(http.StatusOK, map[string]any{
			"commits":  []any{},
			"page":     page,
			"pageSize": pageSize,
		})
		return
	}

	commits, err := c.Repo.Commit.CommitsByPage(page, pageSize, git.CommitsByPageOptions{Path: c.Repo.TreePath})
	if err != nil {
		c.Error(err, "paging commits")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"commits":  matchUsersWithCommitEmails(c.Req.Context(), c.Repo.Repository.ID, commits),
		"page":     page,
		"pageSize": pageSize,
	})
}

func GetCommitStatuses(c *context.APIContext) {
	sha := c.Params(":sha")
	statuses, err := database.Handle.CommitStatuses().GetLatestBySHA(c.Req.Context(), c.Repo.Repository.ID, sha)
	if err != nil {
		c.Error(err, "get latest statuses")
		return
	}

	state := database.Handle.CommitStatuses().GetCombinedStatusBySHA(c.Req.Context(), c.Repo.Repository.ID, sha)
	c.JSON(http.StatusOK, ToCombinedCommitStatusResponse(c.Repo.Repository, sha, state, statuses))
}

func CreateCommitStatus(c *context.APIContext) {
	sha := c.Params(":sha")
	var form struct {
		State       database.CommitStatusState `json:"state" binding:"required"`
		TargetURL   string                     `json:"target_url"`
		Description string                     `json:"description"`
		Context     string                     `json:"context"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.Error(err, "parse json")
		return
	}

	// Validate state
	switch form.State {
	case database.CommitStatusPending, database.CommitStatusSuccess,
		database.CommitStatusError, database.CommitStatusFailure,
		database.CommitStatusWarning:
		// valid
	default:
		c.StatusError(http.StatusUnprocessableEntity, "Invalid commit status state.", nil)
		return
	}

	if form.Context == "" {
		form.Context = "default"
	}

	status := &database.CommitStatus{
		RepoID:      c.Repo.Repository.ID,
		SHA:         sha,
		State:       form.State,
		TargetURL:   form.TargetURL,
		Description: form.Description,
		Context:     form.Context,
		CreatorID:   c.User.ID,
	}

	if err := database.Handle.CommitStatuses().Create(c.Req.Context(), status); err != nil {
		c.Error(err, "create status")
		return
	}

	c.JSON(http.StatusCreated, ToCommitStatusResponse(status))
}
