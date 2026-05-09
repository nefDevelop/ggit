package ginx

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/stretchr/testify/assert"
)

func TestCSRF(t *testing.T) {
	t.Setenv("GGIT_SKIP_RUN_USER_CHECK", "true")
	gin.SetMode(gin.TestMode)
	conf.Security.SecretKey = "test-secret"
	conf.Security.EnableCSRF = true

	r := gin.New()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("gogs_session", store))
	r.Use(context.Contexter(context.NewStore()))
	r.Use(CSRF())

	r.GET("/test", func(c *gin.Context) {
		gctx := c.MustGet("context").(*context.Context)
		c.String(http.StatusOK, gctx.Data["CSRFToken"].(string))
	})

	r.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// 1. Get token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	token := w.Body.String()
	assert.NotEmpty(t, token)
	cookieHeader := w.Header().Get("Set-Cookie")

	// 2. Post without token should fail
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	req.Header.Set("Cookie", cookieHeader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Post with correct token should succeed
	w = httptest.NewRecorder()
	form := url.Values{}
	form.Add("_csrf", token)
	req, _ = http.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookieHeader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 4. Post with X-Csrf-Token header should succeed
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	req.Header.Set("X-Csrf-Token", token)
	req.Header.Set("Cookie", cookieHeader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
