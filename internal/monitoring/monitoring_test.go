package monitoring

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSurvivalModeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Disk space OK", func(t *testing.T) {
		old := diskUsage
		diskUsage = func(path string) (uint64, uint64, error) {
			return 1000, 500, nil // 50% used
		}
		defer func() { diskUsage = old }()

		r := gin.New()
		r.Use(SurvivalModeMiddleware())
		r.POST("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Disk space Critical - Block POST", func(t *testing.T) {
		old := diskUsage
		diskUsage = func(path string) (uint64, uint64, error) {
			return 1000, 20, nil // 98% used
		}
		defer func() { diskUsage = old }()

		r := gin.New()
		r.Use(SurvivalModeMiddleware())
		r.POST("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "Survival Mode")
	})

	t.Run("Disk space Critical - Allow GET", func(t *testing.T) {
		old := diskUsage
		diskUsage = func(path string) (uint64, uint64, error) {
			return 1000, 20, nil // 98% used
		}
		defer func() { diskUsage = old }()

		r := gin.New()
		r.Use(SurvivalModeMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
