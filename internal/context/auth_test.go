package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestToggle_AdminBypassActivation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conf.Security.InstallLock = true
	conf.Auth.RequireEmailConfirmation = true

	t.Run("Admin bypasses activation check", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("context", &Context{
				Context:  c,
				IsLogged: true,
				User: &database.User{
					ID:       1,
					IsActive: false,
					IsAdmin:  true,
				},
			})
		})
		r.GET("/test", Toggle(&ToggleOptions{SignInRequired: true}), func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("Non-admin blocked by activation check", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("context", &Context{
				Context:  c,
				Data:     make(map[string]any),
				IsLogged: true,
				User: &database.User{
					ID:       2,
					IsActive: false,
					IsAdmin:  false,
				},
				Req: c.Request,
			})
		})
		r.GET("/test", Toggle(&ToggleOptions{SignInRequired: true}), func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		// It should redirect to /user/auth/activate (or return 403 if it was API)
		// Since /test is not an API path (doesn't start with /api/), it tries to render HTML.
		// In our lobotomized HTML, it might return 200 with error if status was not set.
		// Wait, Toggle calls gctx.Success("user/auth/activate") which is gctx.HTML(http.StatusOK, "user/auth/activate")
		assert.Equal(t, http.StatusOK, w.Code) 
		// But it shouldn't reach the final handler
		assert.NotEqual(t, "OK", w.Body.String())
	})
	
	t.Run("Non-admin blocked by activation check (API)", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("context", &Context{
				Context:  c,
				Data:     make(map[string]any),
				IsLogged: true,
				User: &database.User{
					ID:       2,
					IsActive: false,
					IsAdmin:  false,
				},
				Req: c.Request,
			})
		})
		r.GET("/api/test", Toggle(&ToggleOptions{SignInRequired: true}), func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Your account is not activated.")
	})
}
