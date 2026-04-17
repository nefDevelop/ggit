package main

import (
	stdctx "context"
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
	log "unknwon.dev/clog/v2"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gogs.io/gogs/docs"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/ginx"
	"gogs.io/gogs/internal/route"
	apiv1 "gogs.io/gogs/internal/route/api/v1"
	apiv2 "gogs.io/gogs/internal/route/api/v2"
	"gogs.io/gogs/internal/route/lfs"
	"gogs.io/gogs/internal/route/repo"
)

var webCommand = cli.Command{
	Name:  "web",
	Usage: "Start Gogs web server",
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "", "Custom configuration file path"),
	},
	Action: runWeb,
}

func runWeb(_ stdctx.Context, cmd *cli.Command) error {
	err := route.GlobalInit(configFromLineage(cmd))
	if err != nil {
		log.Fatal("Failed to initialize application: %v", err)
	}

	r := ginx.New()

	// Git HTTP
	repoHTTP := r.Group("/:username/:reponame", repo.HTTPContexter(repo.NewStore()))
	{
		repoHTTP.POST("/git-upload-pack", repo.HTTP)
		repoHTTP.POST("/git-receive-pack", repo.HTTP)
		repoHTTP.GET("/info/refs", repo.HTTP)
		repoHTTP.GET("/HEAD", repo.HTTP)
		repoHTTP.GET("/objects/*path", repo.HTTP)

		gitSuffix := repoHTTP.Group("/.git")
		{
			gitSuffix.Any("/*action", repo.HTTP)
		}
	}

	// LFS
	lfs.RegisterRoutes(r.Group("/:username/:reponame/info/lfs"))
	lfs.RegisterRoutes(r.Group("/:username/:reponame/.git/info/lfs"))

	// API v1
	v1Group := r.Group("/api", context.APIContexter())
	apiv1.RegisterRoutes(v1Group)

	// API v2
	v2Group := r.Group("/api/v2", context.APIContexter())
	apiv2.RegisterRoutes(v2Group)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// SPA Routing
	r.GET("/", ginx.ServeSPA)
	r.NoRoute(ginx.ServeSPA)

	// Utils

	r.GET("/robots.txt", func(c *gin.Context) {
		if conf.HasRobotsTxt {
			c.File(filepath.Join(conf.CustomDir(), "robots.txt"))
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	// Final server run
	var (
		listenAddr = fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
		protocol   = conf.Server.Protocol
	)

	log.Info("Listen: %s://%s%s", protocol, listenAddr, conf.Server.Subpath)

	server := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	if protocol == "fcgi" {
		f, err := net.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatal("Failed to listen: %v", err)
		}
		return fcgi.Serve(f, r)
	}

	if protocol == "https" {
		return server.ListenAndServeTLS(conf.Server.CertFile, conf.Server.KeyFile)
		}
	return server.ListenAndServe()
}

func configFromLineageDummy(cmd *cli.Command) string {
	// Not used since we use the one from cmd.go
	return ""
}
