package v2

import (
	stdcontext "context"
	"net/http"

	"github.com/gogs/git-module"
	"github.com/go-enry/go-enry/v2"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	reposervice "gogs.io/gogs/internal/services/repo"
	"gogs.io/gogs/internal/template/highlight"
	"gogs.io/gogs/internal/tool"
)

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
			subtree, _ := tree.Subtree(entry.Name())
			_ = walkTreeForStats(subtree, langStats)
			continue
		}
		if enry.IsVendor(entry.Name()) || enry.IsDotFile(entry.Name()) { continue }
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
		res["last_commit"] = lastCommit[0]
	}

	if tool.IsTextFile(p) {
		content := string(p)
		res["content"] = content
		res["highlighted_content"] = highlight.Code(blob.Name(), content)
	}

	c.JSON(http.StatusOK, res)
}

// ... rest of the file (CreateRepo, GetRepoCommits, etc. same as before) ...

// CreateRepo creates a new repository for the authenticated user.
func CreateRepo(c *context.APIContext) {
	var f struct {
		Name        string `json:"name" binding:"required,max=100"`
		Description string `json:"description" binding:"max=255"`
		Private     bool   `json:"private"`
		AutoInit    bool   `json:"auto_init"`
		Gitignores  string `json:"gitignores"`
		License     string `json:"license"`
		Readme      string `json:"readme"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	repo, err := reposervice.Repos.Create(c.Req.Context(), c.User, c.User, reposervice.CreateOptions{
		Name:        f.Name,
		Description: f.Description,
		IsPrivate:   f.Private,
		AutoInit:    f.AutoInit,
		Gitignores:  f.Gitignores,
		License:     f.License,
		Readme:      f.Readme,
	})
	if err != nil {
		c.Error(err, "create repository")
		return
	}

	c.JSON(http.StatusCreated, repo.APIFormat(c.User))
}

type userCommit struct {
	User *database.User `json:"user"`
	*git.Commit
}

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

func GetRepoCommits(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 { page = 1 }
	pageSize := c.QueryInt("pageSize")
	if pageSize <= 0 { pageSize = 20 }
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

func GetRepoBranches(c *context.APIContext) {
	branches, err := c.Repo.Repository.GetBranches()
	if err != nil {
		c.Error(err, "get branches")
		return
	}
	c.JSON(http.StatusOK, branches)
}
