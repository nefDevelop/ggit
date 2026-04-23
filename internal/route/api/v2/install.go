package v2

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/form"
	"gogs.io/gogs/internal/osx"
	"gogs.io/gogs/internal/route"
	"gogs.io/gogs/internal/strx"
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
		"db_types": []string{"PostgreSQL", "SQLite3"},
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

	if errs := f.Validate(c, nil); errs.Len() > 0 {
		c.JSON(http.StatusBadRequest, map[string]any{
			"message": errs[0].Classification,
			"errors":  errs,
		})
		return
	}

	// Basic checks
	if f.AppUrl != "" && f.AppUrl[len(f.AppUrl)-1] != '/' {
		f.AppUrl += "/"
	}

	// Database mapping
	dbTypes := map[string]string{
		"PostgreSQL": "postgres",
		"Postgres":   "postgres",
		"postgres":   "postgres",
		"SQLite3":    "sqlite3",
		"sqlite3":    "sqlite3",
		"sqlite":     "sqlite3",
	}
	conf.Database.Type = dbTypes[f.DbType]

	if conf.Database.Type == "" {
		c.StatusError(http.StatusBadRequest, "Unrecognized database type: "+f.DbType, nil)
		return
	}
	conf.Database.Host = f.DbHost
	conf.Database.User = f.DbUser
	conf.Database.Password = f.DbPasswd
	conf.Database.Name = f.DbName
	conf.Database.Schema = f.DbSchema
	conf.Database.SSLMode = f.SSLMode
	conf.Database.Path = f.DbPath

	if conf.Database.Type == "sqlite3" && conf.Database.Path == "" {
		c.StatusError(http.StatusBadRequest, "Database path cannot be empty for SQLite3.", nil)
		return
	}

	// Test database connection
	if err := database.NewTestEngine(); err != nil {
		c.StatusError(http.StatusBadRequest, "Invalid database settings: "+err.Error(), nil)
		return
	}

	// Ensure paths exist
	_ = os.MkdirAll(f.RepoRootPath, os.ModePerm)
	_ = os.MkdirAll(f.LogRootPath, os.ModePerm)

	// Save settings to app.ini
	cfg := ini.Empty()
	if osx.IsFile(conf.CustomConf) {
		_ = cfg.Append(conf.CustomConf)
	}

	cfg.Section("database").Key("TYPE").SetValue(conf.Database.Type)
	cfg.Section("database").Key("HOST").SetValue(conf.Database.Host)
	cfg.Section("database").Key("NAME").SetValue(conf.Database.Name)
	cfg.Section("database").Key("SCHEMA").SetValue(conf.Database.Schema)
	cfg.Section("database").Key("USER").SetValue(conf.Database.User)
	cfg.Section("database").Key("PASSWORD").SetValue(conf.Database.Password)
	cfg.Section("database").Key("SSL_MODE").SetValue(conf.Database.SSLMode)
	cfg.Section("database").Key("PATH").SetValue(conf.Database.Path)

	cfg.Section("").Key("BRAND_NAME").SetValue(f.AppName)
	cfg.Section("repository").Key("ROOT").SetValue(f.RepoRootPath)
	cfg.Section("").Key("RUN_USER").SetValue(f.RunUser)
	cfg.Section("server").Key("DOMAIN").SetValue(f.Domain)
	cfg.Section("server").Key("HTTP_PORT").SetValue(f.HTTPPort)
	cfg.Section("server").Key("EXTERNAL_URL").SetValue(f.AppUrl)

	if f.SSHPort == 0 {
		cfg.Section("server").Key("DISABLE_SSH").SetValue("true")
	} else {
		cfg.Section("server").Key("DISABLE_SSH").SetValue("false")
		cfg.Section("server").Key("SSH_PORT").SetValue(strconv.Itoa(f.SSHPort))
		cfg.Section("server").Key("START_SSH_SERVER").SetValue(strconv.FormatBool(f.UseBuiltinSSHServer))
	}

	cfg.Section("security").Key("INSTALL_LOCK").SetValue("true")
	secretKey, _ := strx.RandomChars(15)
	cfg.Section("security").Key("SECRET_KEY").SetValue(secretKey)
	cfg.Section("").Key("RUN_MODE").SetValue("prod")

	_ = os.MkdirAll(filepath.Dir(conf.CustomConf), os.ModePerm)
	if err := cfg.SaveTo(conf.CustomConf); err != nil {
		c.StatusError(http.StatusInternalServerError, "Failed to save configuration: "+err.Error(), nil)
		return
	}

	// Global initialization
	if err := route.GlobalInit(conf.CustomConf); err != nil {
		c.StatusError(http.StatusInternalServerError, "Failed to initialize application: "+err.Error(), nil)
		return
	}

	c.NoContent()
}
