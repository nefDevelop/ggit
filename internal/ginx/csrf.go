package ginx

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/strx"
)

// isProtocolPath reports whether the path belongs to the Git or LFS smart-HTTP
// protocols. Those endpoints authenticate with their own mechanisms (HTTP Basic
// or an LFS signature) and are driven by the git/lfs command line clients, which
// neither send cookies nor can be made to carry the XSRF token. Exempting them
// is required for push and LFS upload to keep working once CSRF is enabled.
//
// Matching is done on path segments rather than prefixes because the routes are
// nested under a user and repository: /:username/:reponame/git-receive-pack.
func isProtocolPath(path string) bool {
	for _, segment := range []string{
		"/git-upload-pack",
		"/git-receive-pack",
		"/objects/batch",
		"/objects/basic/",
	} {
		if strings.Contains(path, segment) {
			return true
		}
	}
	// Las rutas LFS cuelgan de /info/lfs y POST /objects/batch responde 200, asi
	// que tambien se marca el prefijo por si cambia el arbol de rutas.
	return strings.Contains(path, "/info/lfs")
}

// CSRF returns a middleware that adds CSRF protection to the application.
// It populates "CSRFToken" and "CSRFTokenHTML" in the context's Data map.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Security.EnableCSRF {
			c.Next()
			return
		}
		if isProtocolPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		sess := sessions.Default(c)
		token, ok := sess.Get("csrf").(string)
		if !ok || token == "" {
			var err error
			token, err = strx.RandomChars(30)
			if err != nil {
				log.Error("Failed to generate random chars for CSRF: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			sess.Set("csrf", token)
			_ = sess.Save()
		}

		// Inject into Ggit context data if it exists
		gctx, exists := c.Get("context")
		if exists {
			ctx := gctx.(*context.Context)
			ctx.Data["CSRFToken"] = token
			ctx.Data["CSRFTokenHTML"] = template.HTML(fmt.Sprintf(`<input type="hidden" name="_csrf" value="%s">`, token))
		}

		// Always provide token to SPA/Frontend via a non-HttpOnly cookie to prevent out-of-sync.
		// We use the configured subpath to ensure cookie is only sent to the app.
		cookiePath := conf.Server.Subpath
		if cookiePath == "" {
			cookiePath = "/"
		}
		c.SetCookie("XSRF-TOKEN", token, 86400, cookiePath, "", false, false)

		// Validate token for mutable requests
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodDelete ||
			c.Request.Method == http.MethodPatch {

			// Bypass CSRF for API requests using token or basic auth
			if ctx, ok := gctx.(*context.Context); ok && (ctx.IsTokenAuth || ctx.IsBasicAuth) {
				c.Next()
				return
			}

			log.Trace("CSRF Check Path: %s", c.Request.URL.Path)

			sentToken := c.GetHeader("X-CSRF-Token")
			if sentToken == "" {
				sentToken = c.GetHeader("X-Csrf-Token")
			}
			if sentToken == "" {
				sentToken = c.PostForm("_csrf")
			}

			matched := false
			for _, t := range strings.Split(sentToken, ",") {
				if strings.TrimSpace(t) == token {
					matched = true
					break
				}
			}

			if !matched {
				// No se registra el token esperado: es el valor valido de la sesion de
				// la victima, y un CSRF fallido lo dejaria escrito en el log.
				log.Warn("CSRF 403: Token mismatch. Path: %s", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "CSRF token mismatch",
				})
				return
			}
		}

		c.Next()
	}
}
