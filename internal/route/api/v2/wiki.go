package v2

import (
	"net/http"
	"strings"

	"github.com/gogs/git-module"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/markup"
)

// ListWikiPages returns a list of wiki pages (Markdown files in /docs).
func ListWikiPages(c *context.APIContext) {
	branch := c.Repo.Repository.DefaultBranch
	commit, err := c.Repo.GitRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	entries, err := commit.Entries(git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.Error(err, "list entries")
		return
	}

	// Look for "docs" or ".docs" directory
	var docsEntry *git.TreeEntry
	for i := range entries {
		if entries[i].Type() == git.ObjectTree && (entries[i].Name() == "docs" || entries[i].Name() == ".docs") {
			docsEntry = entries[i]
			break
		}
	}

	if docsEntry == nil {
		c.JSON(http.StatusOK, []PageMeta{})
		return
	}

	// In git-module, Tree entries are accessed through a commit or sub-tree
	docsTree, err := c.Repo.GitRepo.LsTree(docsEntry.ID().String(), git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.Error(err, "get docs tree")
		return
	}

	entries, err = docsTree.Entries()
	if err != nil {
		c.Error(err, "get docs entries")
		return
	}

	pages := make([]PageMeta, 0, len(entries))
	for i := range entries {
		if entries[i].Type() == git.ObjectBlob && strings.HasSuffix(entries[i].Name(), ".md") {
			name := strings.TrimSuffix(entries[i].Name(), ".md")
			// Get last commit for this file
			commits, err := c.Repo.GitRepo.Log(git.RefsHeads+branch, git.LogOptions{Path: docsEntry.Name() + "/" + entries[i].Name()})
			updated := int64(0)
			if err == nil && len(commits) > 0 {
				updated = commits[0].Author.When.Unix()
			}

			pages = append(pages, PageMeta{
				Name:    name,
				URL:     database.ToWikiPageURL(name),
				Updated: updated,
			})
		}
	}

	c.JSON(http.StatusOK, pages)
}

type PageMeta struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Updated int64  `json:"updated"`
}

// GetWikiPage returns the content of a single doc page.
func GetWikiPage(c *context.APIContext) {
	pageURL := c.Params("page")
	if pageURL == "" {
		pageURL = "README"
	}
	pageName := database.ToWikiPageName(pageURL)

	branch := c.Repo.Repository.DefaultBranch
	commit, err := c.Repo.GitRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	// Try to find the file in /docs/pageName.md or /.docs/pageName.md
	possiblePaths := []string{"docs/" + pageName + ".md", ".docs/" + pageName + ".md"}
	var blob *git.Blob
	var finalPath string

	for _, p := range possiblePaths {
		entry, err := commit.TreeEntry(p)
		if err == nil {
			blob = entry.Blob()
			finalPath = p
			break
		}
	}

	if blob == nil {
		c.NotFound()
		return
	}

	p, err := blob.Bytes()
	if err != nil {
		c.Error(err, "read blob")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"name":         pageName,
		"content":      string(p),
		"content_html": string(markup.Render(markup.TypeMarkdown, p, c.Repo.RepoLink+"/"+finalPath, nil)),
	})
}

// CreateWikiPage is disabled in Docs-as-Code mode.
func CreateWikiPage(c *context.APIContext) {
	c.StatusError(http.StatusNotImplemented, "Direct creation is disabled. Please push files to /docs folder instead.", nil)
}

// UpdateWikiPage is disabled in Docs-as-Code mode.
func UpdateWikiPage(c *context.APIContext) {
	c.StatusError(http.StatusNotImplemented, "Direct editing is disabled. Please push files to /docs folder instead.", nil)
}

// DeleteWikiPage is disabled in Docs-as-Code mode.
func DeleteWikiPage(c *context.APIContext) {
	c.StatusError(http.StatusNotImplemented, "Direct deletion is disabled. Please push files to /docs folder instead.", nil)
}

// GetWikiPageHistory returns the commit history of a single doc page.
func GetWikiPageHistory(c *context.APIContext) {
	pageURL := c.Params("page")
	pageName := database.ToWikiPageName(pageURL)
	branch := c.Repo.Repository.DefaultBranch

	// Check both possible paths
	possiblePaths := []string{"docs/" + pageName + ".md", ".docs/" + pageName + ".md"}
	var finalPath string
	for _, p := range possiblePaths {
		if _, err := c.Repo.GitRepo.Log(git.RefsHeads+branch, git.LogOptions{Path: p}); err == nil {
			finalPath = p
			break
		}
	}

	if finalPath == "" {
		c.NotFound()
		return
	}

	commits, err := c.Repo.GitRepo.Log(git.RefsHeads+branch, git.LogOptions{Path: finalPath})
	if err != nil {
		c.Error(err, "get commits by path")
		return
	}

	result := make([]map[string]any, len(commits))
	for i, commit := range commits {
		result[i] = map[string]any{
			"id":      commit.ID.String(),
			"author":  commit.Author.Name,
			"email":   commit.Author.Email,
			"message": commit.Summary(),
			"created": commit.Author.When.Unix(),
		}
	}

	c.JSON(http.StatusOK, result)
}
