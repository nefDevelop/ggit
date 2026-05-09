package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/authx"
	"github.com/nef734/ggit/internal/conf"
)

func MetricsFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Prometheus.Enabled {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if !conf.Prometheus.EnableBasicAuth {
			return
		}

		username, password := authx.DecodeBasic(c.Request.Header)
		if username != conf.Prometheus.BasicAuthUsername || password != conf.Prometheus.BasicAuthPassword {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
	}
}
