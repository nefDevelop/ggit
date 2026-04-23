package v2

import (
	"net/http"
	"strings"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/form"
)

// GetRepoSettings returns the repository settings.
func GetRepoSettings(c *context.APIContext) {
	repo := c.Repo.Repository
	res := map[string]any{
		"id":            repo.ID,
		"name":          repo.Name,
		"description":   repo.Description,
		"website":       repo.Website,
		"private":       repo.IsPrivate,
		"unlisted":      repo.IsUnlisted,
		"enable_wiki":   repo.EnableWiki,
		"enable_issues": repo.EnableIssues,
		"enable_pulls":  repo.EnablePulls,
	}

	c.JSON(http.StatusOK, res)
}

// UpdateRepoSettings updates the repository settings.
func UpdateRepoSettings(c *context.APIContext) {
	var f form.RepoSetting
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	repo := c.Repo.Repository
	visibilityChanged := repo.IsPrivate != f.Private || repo.IsUnlisted != f.Unlisted
	repo.Name = f.RepoName
	repo.LowerName = strings.ToLower(f.RepoName)
	repo.Description = f.Description
	repo.Website = f.Website
	repo.IsPrivate = f.Private
	repo.IsUnlisted = f.Unlisted

	if err := database.UpdateRepository(repo, visibilityChanged); err != nil {
		c.Error(err, "update repository")
		return
	}

	c.JSON(http.StatusOK, repo)
}

// GetRepoCollaborators returns a list of collaborators for the repository.
func GetRepoCollaborators(c *context.APIContext) {
	users, err := c.Repo.Repository.GetCollaborators()
	if err != nil {
		c.Error(err, "get collaborators")
		return
	}
	c.JSON(http.StatusOK, users)
}

// AddRepoCollaborator adds a new collaborator to the repository.
func AddRepoCollaborator(c *context.APIContext) {
	userName := c.Params("collaborator")
	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	if err = c.Repo.Repository.AddCollaborator(u); err != nil {
		c.Error(err, "add collaborator")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// DeleteRepoCollaborator deletes a collaborator from the repository.
func DeleteRepoCollaborator(c *context.APIContext) {
	userName := c.Params("collaborator")
	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	if err = c.Repo.Repository.DeleteCollaboration(u.ID); err != nil {
		c.Error(err, "delete collaboration")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetRepoHooks returns a list of webhooks for the repository.
func GetRepoHooks(c *context.APIContext) {
	ws, err := database.GetWebhooksByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get webhooks")
		return
	}
	c.JSON(http.StatusOK, ws)
}

// DeleteRepoHook deletes a webhook from the repository.
func DeleteRepoHook(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if err := database.DeleteWebhookOfRepoByID(c.Repo.Repository.ID, id); err != nil {
		c.Error(err, "delete webhook")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// DeleteRepository deletes the repository.
func DeleteRepository(c *context.APIContext) {
	var f struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if f.Name != c.Repo.Repository.Name {
		c.StatusError(http.StatusUnprocessableEntity, "Confirmation name does not match.", nil)
		return
	}

	if err := database.DeleteRepository(c.User.ID, c.Repo.Repository.ID); err != nil {
		c.Error(err, "delete repository")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
