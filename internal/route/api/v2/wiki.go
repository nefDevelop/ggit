package v2

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/markup"
)

type WikiPage struct {
	Name     string      `json:"name"`
	Title    string      `json:"title"`
	Order    int         `json:"order"`
	URL      string      `json:"url"`
	IsFolder bool        `json:"is_folder"`
	Children []*WikiPage `json:"children,omitempty"`
}

var wikiSidebarCache sync.Map // key: repo_id:commit_sha

// ListWikiPages returns a tree of wiki pages (Markdown files in the configured wiki folder).
func ListWikiPages(c *context.APIContext) {
	branch := c.Repo.Repository.DefaultBranch
	commit, err := c.Repo.GitRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	wikiFolder := c.Repo.Repository.WikiFolderName
	if wikiFolder == "" {
		wikiFolder = "docs"
	}

	cacheKey := fmt.Sprintf("%d:%s", c.Repo.Repository.ID, commit.ID.String())
	if val, ok := wikiSidebarCache.Load(cacheKey); ok {
		c.JSON(http.StatusOK, val)
		return
	}

	// Root scan
	tree, err := commit.Subtree(wikiFolder)
	if err != nil {
		// Fallback to .docs if docs doesn't exist and folder is default
		if wikiFolder == "docs" {
			tree, err = commit.Subtree(".docs")
		}
	}

	if err != nil || tree == nil {
		c.JSON(http.StatusOK, []*WikiPage{})
		return
	}

	pages := buildWikiTree(c, tree, "")
	
	// Sort top level
	sortWikiPages(pages)

	wikiSidebarCache.Store(cacheKey, pages)
	c.JSON(http.StatusOK, pages)
}

func buildWikiTree(c *context.APIContext, tree *git.Tree, relPath string) []*WikiPage {
	entries, err := tree.Entries()
	if err != nil {
		return nil
	}

	pages := make([]*WikiPage, 0, len(entries))
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue // ignore hidden
		}

		if entry.IsTree() {
			subtree, err := tree.Subtree(entry.Name())
			if err != nil {
				continue
			}
			subPages := buildWikiTree(c, subtree, path.Join(relPath, entry.Name()))
			if len(subPages) > 0 {
				sortWikiPages(subPages)
				pages = append(pages, &WikiPage{
					Name:     entry.Name(),
					Title:    entry.Name(), // Folder title can also come from a _index.md later
					IsFolder: true,
					Children: subPages,
				})
			}
			continue
		}

		if strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			
			// Read frontmatter for title and order
			title := name
			order := 1000 // default high order
			
			blob, err := entry.Blob().Bytes()
			if err == nil {
				meta, _ := markup.ExtractMetadata(blob)
				if meta != nil {
					if t, ok := meta["title"].(string); ok {
						title = t
					}
					if o, ok := meta["order"].(int); ok {
						order = o
					} else if o, ok := meta["order"].(float64); ok {
						order = int(o)
					}
				}
			}

			pages = append(pages, &WikiPage{
				Name:  name,
				Title: title,
				Order: order,
				URL:   database.ToWikiPageURL(path.Join(relPath, name)),
			})
		}
	}
	return pages
}

func sortWikiPages(pages []*WikiPage) {
	sort.Slice(pages, func(i, j int) bool {
		if pages[i].Order != pages[j].Order {
			return pages[i].Order < pages[j].Order
		}
		return strings.Compare(strings.ToLower(pages[i].Title), strings.ToLower(pages[j].Title)) < 0
	})
}

// GetWikiPage returns the content of a single doc page.
func GetWikiPage(c *context.APIContext) {
	pagePath := c.Params("page")
	if pagePath == "" {
		pagePath = "README"
	}
	
	branch := c.Repo.Repository.DefaultBranch
	commit, err := c.Repo.GitRepo.BranchCommit(branch)
	if err != nil {
		c.Error(err, "get branch commit")
		return
	}

	wikiFolder := c.Repo.Repository.WikiFolderName
	if wikiFolder == "" {
		wikiFolder = "docs"
	}

	// Try multiple paths
	possiblePaths := []string{
		path.Join(wikiFolder, pagePath+".md"),
	}
	if wikiFolder == "docs" {
		possiblePaths = append(possiblePaths, path.Join(".docs", pagePath+".md"))
	}
	
	// Also try the name directly if it ends with .md or is a folder index
	// but the UI usually sends the path without .md

	var blob *git.Blob
	var finalPath string

	for _, p := range possiblePaths {
		entry, err := commit.TreeEntry(p)
		if err == nil && !entry.IsTree() {
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
		"name":         path.Base(pagePath),
		"path":         pagePath,
		"content":      string(p),
		"content_html": string(markup.Render(markup.TypeMarkdown, p, c.Repo.RepoLink+"/src/"+branch+"/"+path.Dir(finalPath), nil)),
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
