package context

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/errx"
)

type APIContext struct {
	*Context
}

func (c *APIContext) Error(status int, msg string, err error) {
	if err != nil {
		log.ErrorDepth(4, "%s: %v", msg, err)
	}
	c.JSON(status, map[string]string{"message": msg})
}

func (c *APIContext) NotFoundOrError(err error, msg string) {
	if errx.IsNotFound(err) {
		c.Error(http.StatusNotFound, msg, err)
		return
	}
	c.Error(http.StatusInternalServerError, msg, err)
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
