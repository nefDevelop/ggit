package v2

import (
	"net/http"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
)

// GetAdminStats returns global statistics.
func GetAdminStats(c *context.APIContext) {
	stats := database.GetStatistic(c.Req.Context())
	c.JSON(http.StatusOK, stats)
}

// GetAdminConfig returns the current application configuration.
func GetAdminConfig(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"app":        conf.App,
		"server":     conf.Server,
		"ssh":        conf.SSH,
		"repository": conf.Repository,
		"database":   conf.Database,
		"security":   conf.Security,
		"email":      conf.Email,
		"auth":       conf.Auth,
		"user":       conf.User,
		"session":    conf.Session,
		"cache":      conf.Cache,
		"attachment": conf.Attachment,
		"release":    conf.Release,
		"picture":    conf.Picture,
		"http":       conf.HTTP,
		"mirror":     conf.Mirror,
		"webhook":    conf.Webhook,
		"git":        conf.Git,
		"lfs":        conf.LFS,
	})
}

// GetAdminUsers returns a list of all users.
func GetAdminUsers(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	users, err := database.Handle.Users().List(c.Req.Context(), page, 50)
	if err != nil {
		c.Error(err, "list users")
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetAdminRepos returns a list of all repositories.
func GetAdminRepos(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	repos, err := database.Repositories(page, 50)
	if err != nil {
		c.Error(err, "list repositories")
		return
	}

	c.JSON(http.StatusOK, repos)
}

// GetAdminOrgs returns a list of all organizations.
func GetAdminOrgs(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	orgs, err := database.Organizations(page, 50)
	if err != nil {
		c.Error(err, "list organizations")
		return
	}

	c.JSON(http.StatusOK, orgs)
}
