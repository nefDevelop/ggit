package main

import (
	stdctx "context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/ginx"
	"github.com/nef734/ggit/internal/monitoring"
	"github.com/nef734/ggit/internal/route"
	apiv2 "github.com/nef734/ggit/internal/route/api/v2"
	"github.com/nef734/ggit/internal/route/lfs"
	"github.com/nef734/ggit/internal/route/repo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var webCommand = cli.Command{
	Name:  "web",
	Usage: "Start Ggit web server",
	Flags: []cli.Flag{
		stringFlag("port, p", "3005", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "", "Custom configuration file path"),
	},
	Action: runWeb,
}

func runWeb(_ stdctx.Context, cmd *cli.Command) error {
	err := route.GlobalInit(configFromLineage(cmd))
	if err != nil {
		log.Fatal("Failed to initialize application: %v", err)
	}

	if cmd.IsSet("port") {
		conf.Server.HTTPPort = cmd.String("port")
	}

	r := ginx.New()

	// Monitoring & Tracing
	if conf.Prometheus.Enabled {
		r.Use(monitoring.PrometheusMiddleware())
		r.GET("/metrics", monitoring.PrometheusHandler())
	}
	r.GET("/health", monitoring.HealthCheckHandler)
	r.Use(monitoring.SurvivalModeMiddleware())

	// OpenTelemetry Tracing
	shutdown, err := monitoring.InitTracing()
	if err != nil {
		log.Error("Failed to initialize tracing: %v", err)
	} else {
		defer shutdown(stdctx.Background())
		r.Use(otelgin.Middleware("ggit"))
	}

	// Git HTTP
	repoHTTP := r.Group("/:username/:reponame", repo.HTTPContexter(repo.NewStore()))
	{
		// LFS
		lfs.RegisterRoutes(repoHTTP.Group("/info/lfs"))
		lfs.RegisterRoutes(repoHTTP.Group("/.git/info/lfs"))

		repoHTTP.POST("/git-upload-pack", repo.HTTP)
		repoHTTP.POST("/git-receive-pack", repo.HTTP)
		repoHTTP.GET("/info/refs", repo.HTTP)
		repoHTTP.GET("/HEAD", repo.HTTP)
		repoHTTP.GET("/objects/*path", repo.HTTP)

		gitSuffix := repoHTTP.Group("/.git")
		{
			gitSuffix.POST("/git-upload-pack", repo.HTTP)
			gitSuffix.POST("/git-receive-pack", repo.HTTP)
			gitSuffix.GET("/info/refs", repo.HTTP)
			gitSuffix.GET("/HEAD", repo.HTTP)
			gitSuffix.GET("/objects/*path", repo.HTTP)
		}
	}

	// API v2
	v2Group := r.Group("/api/v2", context.APIContexter())
	apiv2.RegisterRoutes(v2Group)

	// Swagger UI
	registerSwagger(r.Engine)

	// SPA Routing
	r.GET("/", context.Toggle(&context.ToggleOptions{SignInRequired: false}), ginx.ServeSPA)
	r.NoRoute(context.Toggle(&context.ToggleOptions{SignInRequired: false}), ginx.ServeSPA)

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

	if protocol == "https" {
		return server.ListenAndServeTLS(conf.Server.CertFile, conf.Server.KeyFile)
	}
	return server.ListenAndServe()
}
