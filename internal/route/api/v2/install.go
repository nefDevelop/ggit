package v2

import (
	"net/http"
	"os"
	"strings"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/form"
	"gogs.io/gogs/internal/osx"
)

// GetInstallSettings returns the default settings for the installation page.
func GetInstallSettings(c *context.APIContext) {
	if conf.Security.InstallLock {
		c.StatusError(http.StatusForbidden, "Install lock is enabled.", nil)
		return
	}

	runUser := conf.App.RunUser
	if conf.IsWindowsRuntime() && runUser == "git" {
		runUser = osx.CurrentUsername()
	}

	c.JSON(http.StatusOK, map[string]any{
		"db_types": []string{"MySQL", "PostgreSQL", "SQLite3"},
		"default_settings": map[string]any{
			"db_type":          "PostgreSQL",
			"db_host":          conf.Database.Host,
			"db_user":          conf.Database.User,
			"db_name":          conf.Database.Name,
			"db_schema":        conf.Database.Schema,
			"db_path":          conf.Database.Path,
			"app_name":         conf.App.BrandName,
			"repo_root_path":   conf.Repository.Root,
			"run_user":         runUser,
			"domain":           conf.Server.Domain,
			"ssh_port":         conf.SSH.Port,
			"use_builtin_ssh":  conf.SSH.StartBuiltinServer,
			"http_port":        conf.Server.HTTPPort,
			"app_url":          conf.Server.ExternalURL,
			"log_root_path":    conf.Log.RootPath,
			"default_branch":   conf.Repository.DefaultBranch,
			"smtp_host":        conf.Email.Host,
			"smtp_from":        conf.Email.From,
			"smtp_user":        conf.Email.User,
			"register_confirm": conf.Auth.RequireEmailConfirmation,
			"mail_notify":      conf.User.EnableEmailNotification,
			"offline_mode":     conf.Server.OfflineMode,
			"disable_gravatar": conf.Picture.DisableGravatar,
		},
	})
}

// PerformInstall processes the installation request.
func PerformInstall(c *context.APIContext) {
	if conf.Security.InstallLock {
		c.StatusError(http.StatusForbidden, "Install lock is enabled.", nil)
		return
	}

	var f form.Install
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// This is tricky because we need to reuse the legacy logic which is tied to web context.
	// For now, I'll implement a clean version of the installation logic here or bridge it.
	// Actually, I'll just bridge it to the same internal functions if possible.
	// But since the legacy logic is in 'route/install.go', I might need to move some parts.
	
	// For Phase 11, I'll implement the key parts here and call the existing database init.
	c.JSON(http.StatusNotImplemented, map[string]string{"error": "Installation via API is under construction. Please use the legacy interface for now or wait for full implementation."})
}
