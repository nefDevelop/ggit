package v2

import (
	"bytes"
	stdcontext "context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-enry/go-enry/v2"
	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/markup"
	reposervice "github.com/nef734/ggit/internal/services/repo"
	"github.com/nef734/ggit/internal/template/highlight"
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
	args = append(args, "-e", q, ref)

	if err := tool.AcquireSemaphore(c.Req.Context(), tool.GitGrepSemaphore); err != nil {
		c.Error(err, "acquire grep semaphore")
		return
	}
	defer tool.ReleaseSemaphore(tool.GitGrepSemaphore)

	output, err := git.NewCommand(args...).RunInDir(c.Repo.Repository.RepoPath())
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
		owner, err = database.Handle.Users().GetByID(c.Req.Context(), f.UID)
		if err != nil {
			if database.IsErrUserNotExist(err) {
				c.StatusError(http.StatusNotFound, "Owner not found.", err)
			} else {
				c.Error(err, "get owner")
			}
			return
		}

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
	if c.Repo.Repository.IsBare {
		c.JSON(http.StatusOK, map[string]any{
			"is_bare":    true,
			"repository": c.Repo.Repository.APIFormat(c.Repo.Owner),
		})
		return
	}

	entry, err := c.Repo.Commit.TreeEntry(c.Repo.TreePath)
	if err != nil {
		c.NotFoundOrError(err, "get tree entry")
		return
	}

	if entry.IsTree() {
		getRepoDirectory(c)
	} else {
		getRepoFile(c, entry)
	}
}

func getRepoDirectory(c *context.APIContext) {
	tree, err := c.Repo.Commit.Subtree(c.Repo.TreePath)
	if err != nil {
		c.Error(err, "get subtree")
		return
	}

	entries, err := tree.Entries()
	if err != nil {
		c.Error(err, "list entries")
		return
	}
	entries.Sort()

	type entryData struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Size int64  `json:"size"`
	}

	apiEntries := make([]entryData, len(entries))
	for i, e := range entries {
		apiEntries[i] = entryData{
			Name: e.Name(),
			Type: string(e.Type()),
			Size: e.Size(),
		}
	}

	res := map[string]any{
		"type":       "directory",
		"entries":    apiEntries,
		"path":       c.Repo.TreePath,
		"branch":     c.Repo.BranchName,
		"repository": c.Repo.Repository.APIFormat(c.Repo.Owner),
	}

	// Last commit info for this path
	lastCommit, _ := c.Repo.GitRepo.Log(c.Repo.Commit.ID.String(), git.LogOptions{Path: c.Repo.TreePath, MaxCount: 1})
	if len(lastCommit) > 0 {
		res["last_commit"] = map[string]any{
			"id":      lastCommit[0].ID.String(),
			"summary": lastCommit[0].Summary(),
			"author": map[string]any{
				"name":  lastCommit[0].Author.Name,
				"email": lastCommit[0].Author.Email,
				"when":  lastCommit[0].Author.When,
			},
		}
	}

	// Readme detection and rendering
	var readmeFile *git.Blob
	for _, entry := range entries {
		if entry.IsTree() || !markup.IsReadmeFile(entry.Name()) {
			continue
		}
		readmeFile = entry.Blob()
		break
	}

	if readmeFile != nil {
		p, err := readmeFile.Bytes()
		if err == nil && tool.IsTextFile(p) {
			var rendered []byte
			treeLink := c.Repo.RepoLink + "/src/" + c.Repo.BranchName + "/" + c.Repo.TreePath
			switch markup.Detect(readmeFile.Name()) {
			case markup.TypeMarkdown:
				rendered = markup.Markdown(p, treeLink, c.Repo.Repository.ComposeMetas())
			case markup.TypeOrgMode:
				rendered = markup.OrgMode(p, treeLink, c.Repo.Repository.ComposeMetas())
			default:
				rendered = bytes.ReplaceAll(p, []byte("\n"), []byte(`<br>`))
			}
			res["readme"] = map[string]any{
				"name":    readmeFile.Name(),
				"content": string(rendered),
			}
		}
	}

	// Stats and Languages (Only on root)
	if c.Repo.TreePath == "" {
		stats := map[string]int64{}
		count, _ := c.Repo.Commit.CommitsCount()
		stats["commits"] = count

		branches, _ := c.Repo.Repository.GetBranches()
		stats["branches"] = int64(len(branches))

		tags, _ := c.Repo.GitRepo.Tags()
		stats["tags"] = int64(len(tags))
		res["stats"] = stats

		// Language distribution
		langStats := make(map[string]int64)
		_ = walkTreeForStats(c.Repo.Commit.Tree, langStats)
		res["languages"] = langStats
	}

	c.JSON(http.StatusOK, res)
}

func walkTreeForStats(tree *git.Tree, langStats map[string]int64) error {
	entries, _ := tree.Entries()
	for _, entry := range entries {
		if entry.IsTree() {
			subtree, err := tree.Subtree(entry.Name())
			if err == nil {
				_ = walkTreeForStats(subtree, langStats)
			}
			continue
		}
		if enry.IsVendor(entry.Name()) || enry.IsDotFile(entry.Name()) {
			continue
		}
		lang, _ := enry.GetLanguageByExtension(entry.Name())
		if lang != "" {
			langStats[lang] += entry.Size()
		}
	}
	return nil
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
		res["last_commit"] = map[string]any{
			"id":      lastCommit[0].ID.String(),
			"summary": lastCommit[0].Summary(),
			"author": map[string]any{
				"name":  lastCommit[0].Author.Name,
				"email": lastCommit[0].Author.Email,
				"when":  lastCommit[0].Author.When,
			},
		}
	}

	if tool.IsTextFile(p) {
		content := string(p)
		res["content"] = content
		res["highlighted_content"] = highlight.Code(blob.Name(), content)

		if markup.IsReadmeFile(blob.Name()) {
			treeLink := c.Repo.RepoLink + "/src/" + c.Repo.BranchName + "/" + c.Repo.TreePath
			var rendered []byte
			switch markup.Detect(blob.Name()) {
			case markup.TypeMarkdown:
				rendered = markup.Markdown(p, treeLink, c.Repo.Repository.ComposeMetas())
			case markup.TypeOrgMode:
				rendered = markup.OrgMode(p, treeLink, c.Repo.Repository.ComposeMetas())
			default:
				rendered = bytes.ReplaceAll(p, []byte("\n"), []byte(`<br>`))
			}
			res["readme_rendered"] = string(rendered)
		}
	}

	c.JSON(http.StatusOK, res)
}

type userCommit struct {
	User *database.User `json:"user"`
	*git.Commit
}

// matchUsersWithCommitEmails matches existing users using commit author emails.
func matchUsersWithCommitEmails(ctx stdcontext.Context, oldCommits []*git.Commit) []*userCommit {
	emailToUsers := make(map[string]*database.User)
	newCommits := make([]*userCommit, len(oldCommits))
	usersStore := database.Handle.Users()
	for i := range oldCommits {
		var u *database.User
		if v, ok := emailToUsers[oldCommits[i].Author.Email]; !ok {
			u, _ = usersStore.GetByEmail(ctx, oldCommits[i].Author.Email)
			emailToUsers[oldCommits[i].Author.Email] = u
		} else {
			u = v
		}

		newCommits[i] = &userCommit{
			User:   u,
			Commit: oldCommits[i],
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

	commits, err := c.Repo.Commit.CommitsByPage(page, pageSize, git.CommitsByPageOptions{Path: c.Repo.TreePath})
	if err != nil {
		c.Error(err, "paging commits")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"commits":  matchUsersWithCommitEmails(c.Req.Context(), commits),
		"page":     page,
		"pageSize": pageSize,
	})
}

