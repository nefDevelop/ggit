package ginx

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/strx"
	"strings"
)

// CSRF returns a middleware that adds CSRF protection to the application.
// It populates "CSRFToken" and "CSRFTokenHTML" in the context's Data map.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Security.EnableCSRF {
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
				log.Warn("CSRF 403: Token mismatch. Expected: %s, Got: %s, Path: %s", token, sentToken, c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "CSRF token mismatch",
				})
				return
			}
		}

		c.Next()
	}
}
