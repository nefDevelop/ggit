package main

import (
	stdctx "context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/ginx"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/monitoring"
	"github.com/nef734/ggit/internal/route"
	apiv2 "github.com/nef734/ggit/internal/route/api/v2"
	"github.com/nef734/ggit/internal/route/lfs"
	"github.com/nef734/ggit/internal/route/repo"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// shutdownTimeout es el máximo que se espera a las peticiones en vuelo al
// recibir SIGTERM. Un clone o push grande puede superarlo y se corta.
const shutdownTimeout = 15 * time.Second

var webCommand = cli.Command{
	Name:  "web",
	Usage: "Start Ggit web server",
	Flags: []cli.Flag{
		stringFlag("port, p", "3005", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "", "Custom configuration file path"),
	},
	Action: runWeb,
}

func runWeb(ctx stdctx.Context, cmd *cli.Command) error {
	err := route.GlobalInit(configFromLineage(cmd))
	if err != nil {
		log.Fatal("Failed to initialize application: %v", err)
	}

	if cmd.IsSet("port") {
		conf.Server.HTTPPort = cmd.String("port")
	}

	r := ginx.New()
	r.Use(ginx.MaxBodySize(ginx.DefaultMaxBodySize))

	// Monitoring & Tracing
	if conf.Prometheus.Enabled {
		r.Use(monitoring.PrometheusMiddleware())
		r.GET("/metrics", monitoring.PrometheusHandler())
	}
	r.GET("/health", monitoring.HealthCheckHandler)
	r.Use(monitoring.SurvivalModeMiddleware())

	// pprof
	monitoring.RegisterPprof(r.Engine)

	// OpenTelemetry Tracing
	shutdown, err := monitoring.InitTracing()
	if err != nil {
		log.Error("Failed to initialize tracing: %v", err)
	} else {
		defer shutdown(stdctx.Background())
		r.Use(otelgin.Middleware("ggit"))
	}

	// API v2
	v2Group := r.Group("/api/v2", context.APIContexter())
	apiv2.RegisterRoutes(v2Group)

	// Repository Download Routes (Non-API)
	repoDownloads := r.Group("/:username/:reponame", context.Toggle(&context.ToggleOptions{SignInRequired: false}), context.RepoAssignment())
	{
		repoDownloads.GET("/archive/*name", context.RepoRef(), ginx.Handle(apiv2.GetRepoArchive))
		repoDownloads.GET("/raw/:ref/*path", ginx.Handle(apiv2.GetRepoRawFile))
		repoDownloads.GET("/media/:ref/*path", ginx.Handle(apiv2.GetRepoMediaFile))
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

	// Swagger UI
	registerSwagger(r.Engine)

	// SPA Routing
	r.GET("/", context.Toggle(&context.ToggleOptions{SignInRequired: false}), ginx.ServeSPA)
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			c.Abort()
			return
		}
		context.Toggle(&context.ToggleOptions{SignInRequired: false})(c)
		ginx.ServeSPA(c)
	})

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
		// ReadHeaderTimeout corta los ataques Slowloris: sin el, un cliente puede
		// enviar cabeceras a un byte por segundo y retener la conexion de forma
		// indefinida. Solo limita la recepcion de cabeceras, no del cuerpo.
		ReadHeaderTimeout: 30 * time.Second,
		// IdleTimeout recupera las conexiones keep-alive abandonadas. Sin el, los
		// descriptores se acumulan hasta agotar el limite del proceso.
		IdleTimeout: 120 * time.Second,
		// MaxHeaderBytes acota la memoria por conexion para las cabeceras.
		MaxHeaderBytes: 1 << 20, // 1 MiB
		// ReadTimeout y WriteTimeout se dejan deliberadamente sin fijar: este mismo
		// servidor atiende el protocolo Git HTTP, donde un clone o un push de un
		// repositorio grande mantiene la conexion abierta durante minutos. Un
		// WriteTimeout corto abortaria transferencias legitimas.
	}

	serve := server.ListenAndServe
	if protocol == "https" {
		serve = func() error {
			return server.ListenAndServeTLS(conf.Server.CertFile, conf.Server.KeyFile)
		}
	}

	// Sin este cierre ordenado, un SIGTERM mataba el proceso sin terminar las
	// peticiones en vuelo ni dar ocasión de parar procesos hijos.
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- serve()
	}()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-sigCtx.Done():
	}

	log.Info("Shutting down: waiting up to %s for in-flight requests", shutdownTimeout)
	shutdownCtx, cancel := stdctx.WithTimeout(stdctx.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
