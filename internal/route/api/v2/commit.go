package v2

import (
	"net/http"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/gitx"
)

// GetSingleCommit returns the details and diff of a single commit.
func GetSingleCommit(c *context.APIContext) {
	sha := c.Params("sha")
	if len(sha) == 0 {
		c.ErrorStatus(http.StatusBadRequest, nil)
		return
	}

	commit, err := c.Repo.GitRepo.CatFileCommit(sha)
	if err != nil {
		c.NotFoundOrError(err, "get commit")
		return
	}

	// For diff, we need the parent commit if it exists.
	// We can use gitx.RepoDiff but we need the base commit.
	var base string
	if commit.ParentCount() > 0 {
		parent, err := commit.Parent(0)
		if err == nil {
			base = parent.ID.String()
		}
	}

	diff, err := gitx.RepoDiff(c.Repo.GitRepo,
		commit.ID.String(), 500, 5000, 500,
		git.DiffOptions{Base: base},
	)
	if err != nil {
		c.Error(err, "get diff")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"commit": map[string]any{
			"id":      commit.ID.String(),
			"summary": commit.Summary(),
			"message": commit.Message,
			"author": map[string]any{
				"name":  commit.Author.Name,
				"email": commit.Author.Email,
				"when":  commit.Author.When,
			},
			"committer": map[string]any{
				"name":  commit.Committer.Name,
				"email": commit.Committer.Email,
				"when":  commit.Committer.When,
			},
		},
		"diff": diff,
	})
}
