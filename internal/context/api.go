package context

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/errx"
	log "github.com/nef734/ggit/internal/log"
)

type APIContext struct {
	*Context
}

func (c *APIContext) Error(err error, msg string) {
	c.StatusError(http.StatusInternalServerError, msg, err)
}

func (c *APIContext) Errorf(err error, format string, args ...any) {
	c.StatusError(http.StatusInternalServerError, fmt.Sprintf(format, args...), err)
}

func (c *APIContext) StatusError(status int, msg string, err error) {
	if err != nil {
		log.ErrorDepth(4, "%s: %v", msg, err)
	}
	c.JSON(status, map[string]string{"message": msg})
	c.Abort()
}

func (c *APIContext) ErrorStatus(status int, err error) {
	msg := "Unknown error"
	if err != nil {
		msg = err.Error()
	}
	c.StatusError(status, msg, err)
}

func (c *APIContext) NoContent() {
	c.Status(http.StatusNoContent)
}

func (c *APIContext) NotFoundOrError(err error, msg string) {
	if errx.IsNotFound(err) {
		c.StatusError(http.StatusNotFound, msg, err)
		return
	}
	c.StatusError(http.StatusInternalServerError, msg, err)
}

func APIContexter() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx, exists := c.Get("context")
		if !exists {
			log.Error("APIContexter: context not found")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		actx := &APIContext{Context: gctx.(*Context)}
		c.Set("apiContext", actx)
		c.Next()
	}
}

func (c *APIContext) QueryInt(name string) int {
	return c.Context.QueryInt(name)
}
func (c *APIContext) QueryInt64(name string) int64 {
	return c.Context.QueryInt64(name)
}
func (c *APIContext) ParamsInt64(name string) int64 {
	return c.Context.ParamsInt64(name)
}

func (c *APIContext) SetLinkHeader(total, pageSize int) {
	// Paging link header
	// Link: <url>; rel="next", <url>; rel="last"
	// FIXME: implement actual logic if needed
}
