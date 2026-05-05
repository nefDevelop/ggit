package context

import (
	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/database"
)

// ParamsUser is the wrapper type of the target user defined by URL parameter, namely ':username'.
type ParamsUser struct {
	*database.User
}

// InjectParamsUser returns a handler that retrieves target user based on URL parameter ':username',
// and injects it as *ParamsUser.
func InjectParamsUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		u, err := database.Handle.Users().GetByUsername(c.Request.Context(), c.Param("username"))
		if err != nil {
			gctx.NotFoundOrError(err, "get user by name")
			return
		}
		c.Set("paramsUser", &ParamsUser{database.FromDomainUser(u)})
	}
}
