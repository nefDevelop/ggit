package v2

import (
	"net/http"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/gitx"
)

// GetSingleCommit returns the details and diff of a single commit.
// @Summary Get Single Commit
// @Description Returns the details and full diff of a specific commit by its SHA.
// @Tags repository
// @Produce json
// @Param username path string true "Username"
// @Param reponame path string true "Repository name"
// @Param sha path string true "Commit SHA"
// @Success 200 {object} map[string]any
// @Router /:username/:reponame/commit/{sha} [get]
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
	var opts git.DiffOptions
	if commit.ParentsCount() > 0 {
		parent, err := commit.Parent(0)
		if err == nil {
			opts.Base = parent.ID.String()
		}
	}

	diff, err := gitx.RepoDiff(c.Repo.GitRepo,
		commit.ID.String(), 500, 5000, 500,
		opts,
	)
	if err != nil {
		c.Error(err, "get diff")
		return
	}

	statuses, _ := database.Handle.CommitStatuses().GetLatestBySHA(c.Req.Context(), c.Repo.Repository.ID, commit.ID.String())
	combined := database.Handle.CommitStatuses().GetCombinedStatusBySHA(c.Req.Context(), c.Repo.Repository.ID, commit.ID.String())

	commitUser, _ := database.Handle.Users().GetByEmail(c.Req.Context(), commit.Author.Email)
	authorAvatarURL := resolveUserAvatarURL(database.FromDomainUser(commitUser))

	c.JSON(http.StatusOK, map[string]any{
		"commit": map[string]any{
			"id":      commit.ID.String(),
			"summary": commit.Summary(),
			"message": commit.Message,
			"author": map[string]any{
				"name":       commit.Author.Name,
				"email":      commit.Author.Email,
				"when":       commit.Author.When,
				"avatar_url": authorAvatarURL,
			},
			"committer": map[string]any{
				"name":  commit.Committer.Name,
				"email": commit.Committer.Email,
				"when":  commit.Committer.When,
			},
			"status":   string(combined),
			"statuses": statuses,
		},
		"diff": diff,
	})
}
