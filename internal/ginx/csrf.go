package ginx

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/strx"
)

// CSRF returns a middleware that adds CSRF protection to the application.
// It populates "CSRFToken" and "CSRFTokenHTML" in the context's Data map.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Provide token to SPA/Frontend via a non-HttpOnly cookie
		if c.Request.Method == http.MethodGet {
			// Set cookie with explicit domain if possible, and Path=/
			c.SetCookie("XSRF-TOKEN", token, 86400, "/", "", false, false)
		}

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

			sentToken := c.PostForm("_csrf")
			if sentToken == "" {
				sentToken = c.GetHeader("X-CSRF-Token")
			}
			if sentToken == "" {
				sentToken = c.GetHeader("X-Csrf-Token")
			}

			if sentToken != token {
				log.Warn("CSRF token mismatch: expected %s, got %s", token, sentToken)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "CSRF token mismatch",
				})
				return
			}
		}

		c.Next()
	}
}
