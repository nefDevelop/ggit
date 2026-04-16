package v2

import (
	stdcontext "context"
	"net/http"

	"github.com/gogs/git-module"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	reposervice "gogs.io/gogs/internal/services/repo"
	"gogs.io/gogs/internal/template/highlight"
	"gogs.io/gogs/internal/tool"
)

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

	c.JSON(http.StatusCreated, repo.APIFormatLegacy(nil))
}

// GetRepoView returns the content of a file or directory in the repository.
func GetRepoView(c *context.APIContext) {
	if c.Repo.Repository.IsBare || c.Repo.Commit == nil {
		c.JSON(http.StatusOK, map[string]any{
			"repository": c.Repo.Repository.APIFormatLegacy(nil),
			"is_bare":    true,
		})
		return
	}

	entry, err := c.Repo.Commit.TreeEntry(c.Repo.TreePath)
	if err != nil {
		// If path is empty, it's the root directory
		if c.Repo.TreePath == "" {
			getRepoDirectory(c)
			return
		}
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

	c.JSON(http.StatusOK, map[string]any{
		"repository": c.Repo.Repository.APIFormatLegacy(nil),
		"type":       "directory",
		"entries":    apiEntries,
		"path":       c.Repo.TreePath,
		"branch":     c.Repo.BranchName,
	})
}

func getRepoFile(c *context.APIContext, entry *git.TreeEntry) {
	blob := entry.Blob()
	p, err := blob.Bytes()
	if err != nil {
		c.Error(err, "read blob")
		return
	}

	isTextFile := tool.IsTextFile(p)

	res := map[string]any{
		"type":         "file",
		"name":         blob.Name(),
		"size":         blob.Size(),
		"is_text_file": isTextFile,
		"path":         c.Repo.TreePath,
		"branch":       c.Repo.BranchName,
	}

	if isTextFile {
		content := string(p)
		res["content"] = content
		res["highlighted_content"] = highlight.Code(blob.Name(), content)
	}

	res["repository"] = c.Repo.Repository.APIFormatLegacy(nil)
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

// GetRepoBranches returns a list of branches for the repository.
func GetRepoBranches(c *context.APIContext) {
	branches, err := c.Repo.Repository.GetBranches()
	if err != nil {
		c.Error(err, "get branches")
		return
	}

	c.JSON(http.StatusOK, branches)
}
