package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	"github.com/stretchr/testify/assert"
)

// runAuthz executes mw against a synthetic request carrying gctx and reports the
// resulting status code together with whether the wrapped handler was reached.
//
// The probe path is deliberately under /api/v2 because that is the branch
// Context.NotFound uses to answer with JSON and, crucially, to abort. If the
// middleware stops calling Abort the handler would still run, which is exactly
// the regression these tests exist to catch.
func runAuthz(t *testing.T, mw gin.HandlerFunc, gctx *Context) (code int, reached bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	router := gin.New()

	router.Use(func(c *gin.Context) {
		gctx.Context = c
		gctx.Req = c.Request
		gctx.Resp = c.Writer
		c.Set("context", gctx)
		c.Next()
	})
	router.GET("/api/v2/probe", mw, func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v2/probe", nil))
	return recorder.Code, reached
}

func repoContextFor(mode domain.AccessMode, logged bool, siteAdmin bool) *Context {
	gctx := &Context{
		Context:  nil,
		IsLogged: logged,
		Repo:     &Repository{AccessMode: mode},
	}
	if logged {
		gctx.User = &database.User{ID: 2, IsAdmin: siteAdmin}
	}
	return gctx
}

func TestRequireRepoWriter(t *testing.T) {
	tests := []struct {
		name      string
		mode      domain.AccessMode
		logged    bool
		siteAdmin bool
		wantAllow bool
	}{
		{name: "anonymous is rejected", mode: domain.AccessModeNone, logged: false, wantAllow: false},
		{name: "no access is rejected", mode: domain.AccessModeNone, logged: true, wantAllow: false},
		{name: "read access is rejected", mode: domain.AccessModeRead, logged: true, wantAllow: false},
		{name: "write access is allowed", mode: domain.AccessModeWrite, logged: true, wantAllow: true},
		{name: "admin access is allowed", mode: domain.AccessModeAdmin, logged: true, wantAllow: true},
		{name: "owner is allowed", mode: domain.AccessModeOwner, logged: true, wantAllow: true},
		{name: "site admin without repo access is allowed", mode: domain.AccessModeNone, logged: true, siteAdmin: true, wantAllow: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gctx := repoContextFor(tt.mode, tt.logged, tt.siteAdmin)
			code, reached := runAuthz(t, RequireRepoWriter(), gctx)

			if tt.wantAllow {
				assert.True(t, reached, "the handler must run when access is granted")
				assert.Equal(t, http.StatusOK, code)
				return
			}

			assert.False(t, reached, "the handler must not run when access is denied")
			assert.Equal(t, http.StatusNotFound, code, "denial must answer 404 to avoid disclosing the repository")
		})
	}
}

func TestRequireRepoAdmin(t *testing.T) {
	tests := []struct {
		name      string
		mode      domain.AccessMode
		logged    bool
		siteAdmin bool
		wantAllow bool
	}{
		{name: "anonymous is rejected", mode: domain.AccessModeOwner, logged: false, wantAllow: false},
		{name: "read access is rejected", mode: domain.AccessModeRead, logged: true, wantAllow: false},
		{name: "write access is rejected", mode: domain.AccessModeWrite, logged: true, wantAllow: false},
		{name: "admin access is allowed", mode: domain.AccessModeAdmin, logged: true, wantAllow: true},
		{name: "owner is allowed", mode: domain.AccessModeOwner, logged: true, wantAllow: true},
		{name: "site admin without repo access is allowed", mode: domain.AccessModeNone, logged: true, siteAdmin: true, wantAllow: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gctx := repoContextFor(tt.mode, tt.logged, tt.siteAdmin)
			code, reached := runAuthz(t, RequireRepoAdmin(), gctx)

			if tt.wantAllow {
				assert.True(t, reached, "the handler must run when access is granted")
				assert.Equal(t, http.StatusOK, code)
				return
			}

			assert.False(t, reached, "the handler must not run when access is denied")
			assert.Equal(t, http.StatusNotFound, code)
		})
	}
}

// TestDeniedRequestsDoNotReachTheHandler is a regression guard for the whole
// class of bug these middlewares exist to prevent: a denial that forgets to
// abort lets the handler run anyway, which is what allowed unauthenticated
// callers to delete repositories, branches and organizations.
func TestDeniedRequestsDoNotReachTheHandler(t *testing.T) {
	middlewares := map[string]gin.HandlerFunc{
		"RequireRepoWriter": RequireRepoWriter(),
		"RequireRepoAdmin":  RequireRepoAdmin(),
	}

	for name, mw := range middlewares {
		t.Run(name, func(t *testing.T) {
			_, reached := runAuthz(t, mw, repoContextFor(domain.AccessModeNone, false, false))
			assert.False(t, reached, "an anonymous request must be stopped by "+name)
		})
	}
}
