package route

import (
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gogs/git-module"
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/cron"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/email"
	"github.com/nef734/ggit/internal/markup"
	"github.com/nef734/ggit/internal/ssh"
	"github.com/nef734/ggit/internal/template/highlight"
	"github.com/nef734/ggit/internal/services"
)

func checkRunMode() {
	if conf.IsProdMode() {
		gin.SetMode(gin.ReleaseMode)
		git.SetOutput(nil)
	} else {
		git.SetOutput(os.Stdout)
	}
	log.Info("Run mode: %s", strings.Title(gin.Mode()))
}

// GlobalInit is for global configuration reload-able.
func GlobalInit(customConf string) error {
	err := conf.Init(customConf)
	if err != nil {
		return errors.Wrap(err, "init configuration")
	}

	conf.InitLogging(false)
	log.Info("%s %s", conf.App.BrandName, conf.App.Version)
	log.Trace("Work directory: %s", conf.WorkDir())
	log.Trace("Custom path: %s", conf.CustomDir())
	log.Trace("Custom config: %s", conf.CustomConf)
	log.Trace("Log path: %s", conf.Log.RootPath)
	log.Trace("Build time: %s", conf.BuildTime)
	log.Trace("Build commit: %s", conf.BuildCommit)

	if conf.Email.Enabled {
		log.Trace("Email service is enabled")
	}

	email.NewContext()

	if conf.Security.InstallLock {
		highlight.NewContext()
		markup.NewSanitizer()
		err := database.NewEngine()
		if err != nil {
			log.Fatal("Failed to initialize ORM engine: %v", err)
		}
		database.HasEngine = true

		services.Init()

		database.LoadRepoConfig()
		database.NewRepoContext()

		// Booting long running goroutines.
		cron.NewContext()
		database.InitDeliverHooks()
		database.InitTestPullRequests()
	}
	if conf.HasMinWinSvc {
		log.Info("Builtin Windows Service is supported")
	}
	if conf.Server.LoadAssetsFromDisk {
		log.Trace("Assets are loaded from disk")
	}
	checkRunMode()

	if !conf.Security.InstallLock {
		return nil
	}

	if conf.SSH.StartBuiltinServer {
		ssh.Listen(conf.SSH, conf.Server.AppDataPath)
		log.Info("SSH server started on %s:%v", conf.SSH.ListenHost, conf.SSH.ListenPort)
		log.Trace("SSH server cipher list: %v", conf.SSH.ServerCiphers)
		log.Trace("SSH server MAC list: %v", conf.SSH.ServerMACs)
		log.Trace("SSH server algorithms: %v", conf.SSH.ServerAlgorithms)
	}

	if conf.SSH.RewriteAuthorizedKeysAtStart {
		if err := database.RewriteAuthorizedKeys(); err != nil {
			log.Warn("Failed to rewrite authorized_keys file: %v", err)
		}
	}

	return nil
}
