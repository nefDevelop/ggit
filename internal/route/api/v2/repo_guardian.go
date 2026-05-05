package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
)

// GetRepoGuardianSettings returns Secret Guardian settings for the repository.
func GetRepoGuardianSettings(c *context.APIContext) {
	repo := c.Repo.Repository
	patterns, err := database.Handle.SecretPatterns().ListByRepo(c.Req.Context(), repo.ID)
	if err != nil {
		c.Error(err, "list patterns")
		return
	}

	c.JSONSuccess(map[string]any{
		"enabled":  repo.EnableSecretGuardian,
		"strict":   repo.SecretGuardianStrict,
		"patterns": patterns,
	})
}

// UpdateRepoGuardianSettings updates Secret Guardian settings.
func UpdateRepoGuardianSettings(c *context.APIContext) {
	var f struct {
		Enabled *bool `json:"enabled"`
		Strict  *bool `json:"strict"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	repo := c.Repo.Repository
	updates := make(map[string]any)
	if f.Enabled != nil {
		updates["enable_secret_guardian"] = *f.Enabled
	}
	if f.Strict != nil {
		updates["secret_guardian_strict"] = *f.Strict
	}

	if err := database.Handle.DB().Model(repo).Updates(updates).Error; err != nil {
		c.Error(err, "update guardian settings")
		return
	}

	c.NoContent()
}

// AddRepoGuardianPattern adds a custom pattern to the repository.
func AddRepoGuardianPattern(c *context.APIContext) {
	var f struct {
		Name        string `json:"name" binding:"required"`
		Regex       string `json:"regex" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	pattern := &domain.SecretPattern{
		Name:        f.Name,
		Regex:       f.Regex,
		Description: f.Description,
		RepoID:      c.Repo.Repository.ID,
		IsGlobal:    false,
	}

	if err := database.Handle.SecretPatterns().Create(c.Req.Context(), pattern); err != nil {
		c.Error(err, "create pattern")
		return
	}

	c.JSONSuccess(pattern)
}

// DeleteRepoGuardianPattern deletes a custom pattern.
func DeleteRepoGuardianPattern(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if err := database.Handle.SecretPatterns().DeleteByID(c.Req.Context(), id); err != nil {
		c.Error(err, "delete pattern")
		return
	}
	c.NoContent()
}

// GetAdminGuardianPatterns returns global patterns.
func GetAdminGuardianPatterns(c *context.APIContext) {
	patterns, err := database.Handle.SecretPatterns().ListGlobal(c.Req.Context())
	if err != nil {
		c.Error(err, "list global patterns")
		return
	}
	c.JSONSuccess(patterns)
}

// AddAdminGuardianPattern adds a global pattern.
func AddAdminGuardianPattern(c *context.APIContext) {
	var f struct {
		Name        string `json:"name" binding:"required"`
		Regex       string `json:"regex" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	pattern := &domain.SecretPattern{
		Name:        f.Name,
		Regex:       f.Regex,
		Description: f.Description,
		IsGlobal:    true,
	}

	if err := database.Handle.SecretPatterns().Create(c.Req.Context(), pattern); err != nil {
		c.Error(err, "create global pattern")
		return
	}

	c.JSONSuccess(pattern)
}

// DeleteAdminGuardianPattern deletes a global pattern.
func DeleteAdminGuardianPattern(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if err := database.Handle.SecretPatterns().DeleteByID(c.Req.Context(), id); err != nil {
		c.Error(err, "delete global pattern")
		return
	}
	c.NoContent()
}

