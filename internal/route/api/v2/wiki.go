package v2

import (
	"net/http"
	"strings"

	"github.com/gogs/git-module"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/form"
	"gogs.io/gogs/internal/markup"
)

// ListWikiPages returns a list of wiki pages for the repository.
func ListWikiPages(c *context.APIContext) {
	if !c.Repo.Repository.HasWiki() {
		c.JSON(http.StatusOK, []PageMeta{})
		return
	}

	wikiPath := c.Repo.Repository.WikiPath()
	wikiRepo, err := git.Open(wikiPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	branch := database.WikiBranch(wikiPath)
	commit, err := wikiRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	entries, err := commit.Entries(git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.Error(err, "list entries")
		return
	}

	pages := make([]PageMeta, 0, len(entries))
	for i := range entries {
		if entries[i].Type() == git.ObjectBlob && strings.HasSuffix(entries[i].Name(), ".md") {
			commits, err := wikiRepo.Log(git.RefsHeads+branch, git.LogOptions{Path: entries[i].Name()})
			if err != nil {
				c.Error(err, "get commits by path")
				return
			}
			name := strings.TrimSuffix(entries[i].Name(), ".md")
			pages = append(pages, PageMeta{
				Name:    name,
				URL:     database.ToWikiPageURL(name),
				Updated: commits[0].Author.When.Unix(),
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

// GetWikiPage returns the content of a single wiki page.
func GetWikiPage(c *context.APIContext) {
	pageURL := c.Params("page")
	if pageURL == "" {
		pageURL = "Home"
	}

	wikiPath := c.Repo.Repository.WikiPath()
	wikiRepo, err := git.Open(wikiPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	branch := database.WikiBranch(wikiPath)
	commit, err := wikiRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	pageName := database.ToWikiPageName(pageURL)
	blob, err := commit.Blob(pageName + ".md")
	if err != nil {
		c.NotFoundOrError(err, "get blob")
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
		"content_html": string(markup.Render(markup.TypeMarkdown, p, c.Repo.RepoLink, nil)),
	})
}

// CreateWikiPage creates a new wiki page.
func CreateWikiPage(c *context.APIContext) {
	var f form.NewWiki
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := c.Repo.Repository.AddWikiPage(c.User, f.Title, f.Content, f.Message); err != nil {
		if database.IsErrWikiAlreadyExist(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Wiki page already exists.", err)
		} else {
			c.Error(err, "add wiki page")
		}
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
		"url":    database.ToWikiPageURL(database.ToWikiPageName(f.Title)),
	})
}

// UpdateWikiPage updates an existing wiki page.
func UpdateWikiPage(c *context.APIContext) {
	pageURL := c.Params("page")
	oldTitle := database.ToWikiPageName(pageURL)

	var f form.NewWiki
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := c.Repo.Repository.EditWikiPage(c.User, oldTitle, f.Title, f.Content, f.Message); err != nil {
		c.Error(err, "edit wiki page")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"url":    database.ToWikiPageURL(database.ToWikiPageName(f.Title)),
	})
}

// DeleteWikiPage deletes a wiki page.
func DeleteWikiPage(c *context.APIContext) {
	pageURL := c.Params("page")
	if pageURL == "" {
		pageURL = "Home"
	}

	pageName := database.ToWikiPageName(pageURL)
	if err := c.Repo.Repository.DeleteWikiPage(c.User, pageName); err != nil {
		c.Error(err, "delete wiki page")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}
// GetWikiPageHistory returns the commit history of a single wiki page.
func GetWikiPageHistory(c *context.APIContext) {
	pageURL := c.Params("page")
	if pageURL == "" {
		pageURL = "Home"
	}

	wikiPath := c.Repo.Repository.WikiPath()
	wikiRepo, err := git.Open(wikiPath)
	if err != nil {
		c.Error(err, "open repository")
		return
	}

	branch := database.WikiBranch(wikiPath)
	pageName := database.ToWikiPageName(pageURL)
	fileName := pageName + ".md"

	commits, err := wikiRepo.Log(git.RefsHeads+branch, git.LogOptions{Path: fileName})
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
