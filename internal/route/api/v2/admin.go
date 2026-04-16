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

// UpdateAdminUser updates a user's information administratively.
func UpdateAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	var f struct {
		FullName        *string `json:"full_name"`
		Website         *string `json:"website"`
		Location        *string `json:"location"`
		MaxRepoCreation *int    `json:"max_repo_creation"`
		IsActivated     *bool   `json:"activated"`
		IsAdmin         *bool   `json:"admin"`
		ProhibitLogin   *bool   `json:"prohibit_login"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	opts := database.UpdateUserOptions{
		FullName:        f.FullName,
		Website:         f.Website,
		Location:        f.Location,
		MaxRepoCreation: f.MaxRepoCreation,
		IsActivated:      f.IsActivated,
		IsAdmin:          f.IsAdmin,
		ProhibitLogin:    f.ProhibitLogin,
	}

	if err := database.Handle.Users().Update(c.Req.Context(), userID, opts); err != nil {
		c.Error(err, "update user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteAdminUser deletes a user account.
func DeleteAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	if err := database.Handle.Users().DeleteByID(c.Req.Context(), userID, false); err != nil {
		c.Error(err, "delete user")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetAdminAuths returns a list of all login sources.
func GetAdminAuths(c *context.APIContext) {
	sources, err := database.Handle.LoginSources().List(c.Req.Context(), database.ListLoginSourceOptions{})
	if err != nil {
		c.Error(err, "list login sources")
		return
	}
	c.JSON(http.StatusOK, sources)
}

// DeleteAdminAuth deletes a login source.
func DeleteAdminAuth(c *context.APIContext) {
	authID := c.ParamsInt64("id")
	if err := database.Handle.LoginSources().DeleteByID(c.Req.Context(), authID); err != nil {
		c.Error(err, "delete login source")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetAdminNotices returns a list of system notices.
func GetAdminNotices(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	notices, err := database.Handle.Notices().List(c.Req.Context(), page, 50)
	if err != nil {
		c.Error(err, "list notices")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"notices": notices,
		"total":   database.Handle.Notices().Count(c.Req.Context()),
	})
}

// DeleteAdminNotice deletes a system notice.
func DeleteAdminNotice(c *context.APIContext) {
	noticeID := c.ParamsInt64("id")
	if err := database.Handle.Notices().DeleteByIDs(c.Req.Context(), noticeID); err != nil {
		c.Error(err, "delete notice")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// ClearAdminNotices deletes all system notices.
func ClearAdminNotices(c *context.APIContext) {
	if err := database.Handle.Notices().DeleteAll(c.Req.Context()); err != nil {
		c.Error(err, "clear notices")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
