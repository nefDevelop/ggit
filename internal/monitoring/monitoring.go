package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/osutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	startTime = time.Now()
)

// PrometheusHandler returns the prometheus metrics handler.
func PrometheusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.Prometheus.EnableBasicAuth {
			user, pass, ok := c.Request.BasicAuth()
			if !ok || user != conf.Prometheus.BasicAuthUsername || pass != conf.Prometheus.BasicAuthPassword {
				c.Header("WWW-Authenticate", `Basic realm="Prometheus Metrics"`)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		metrics.WritePrometheus(c.Writer, true)
	}
}

// PrometheusMiddleware collects metrics for HTTP requests.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := http.StatusText(c.Writer.Status())
		method := c.Request.Method

		metrics.GetOrCreateCounter(fmt.Sprintf(`ggit_http_requests_total{method=%q,path=%q,status=%q}`, method, path, status)).Inc()
		metrics.GetOrCreateHistogram(fmt.Sprintf(`ggit_http_request_duration_seconds{method=%q,path=%q}`, method, path)).Update(duration)
	}
}

// HealthStatus represents the health of the system.
type HealthStatus struct {
	Status    string                 `json:"status"`
	Uptime    string                 `json:"uptime"`
	Database  string                 `json:"database"`
	Disk      map[string]interface{} `json:"disk"`
	Memory    map[string]interface{} `json:"memory"`
	Survival  bool                   `json:"survival_mode"`
	Timestamp time.Time              `json:"timestamp"`
}

// HealthCheckHandler returns the health status of the system.
func HealthCheckHandler(c *gin.Context) {
	status := "UP"
	dbStatus := "OK"

	// Check Database
	if database.Handle != nil {
		sqlDB, err := database.Handle.GormDB().DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "DOWN"
			status = "DEGRADED"
		}
	} else {
		dbStatus = "NOT_INITIALIZED"
		status = "DEGRADED"
	}

	// Check Disk
	diskInfo := make(map[string]interface{})
	survival := false
	total, free, err := osutil.DiskUsage(conf.Repository.Root)
	if err == nil {
		used := total - free
		usedPercent := float64(used) / float64(total) * 100
		diskInfo["total"] = total
		diskInfo["free"] = free
		diskInfo["used_percent"] = usedPercent
		if usedPercent > 95 {
			status = "CRITICAL"
			survival = true
		}
	} else {
		diskInfo["error"] = err.Error()
	}

	// Check Memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memInfo := map[string]interface{}{
		"alloc":       m.Alloc,
		"total_alloc": m.TotalAlloc,
		"sys":         m.Sys,
		"num_gc":      m.NumGC,
	}

	c.JSON(http.StatusOK, HealthStatus{
		Status:    status,
		Uptime:    time.Since(startTime).String(),
		Database:  dbStatus,
		Disk:      diskInfo,
		Memory:    memInfo,
		Survival:  survival,
		Timestamp: time.Now(),
	})
}

// InitTracing initializes OpenTelemetry tracing.
func InitTracing() (func(context.Context) error, error) {
	// Only enable stdout tracing if explicitly requested via environment variable
	if os.Getenv("GGIT_TRACE_STDOUT") != "true" {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ggit"),
			attribute.String("environment", os.Getenv("GIN_MODE")),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

var diskUsage = osutil.DiskUsage

// SurvivalModeMiddleware blocks write requests if disk space is critical.
func SurvivalModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead && c.Request.Method != http.MethodOptions {
			total, free, err := diskUsage(conf.Repository.Root)
			if err == nil {
				used := total - free
				usedPercent := float64(used) / float64(total) * 100
				if usedPercent > 95 {
					c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
						"error": "Survival Mode: Disk space is critical (>95%). Write operations are blocked.",
					})
					return
				}
			}
		}
		c.Next()
	}
}
