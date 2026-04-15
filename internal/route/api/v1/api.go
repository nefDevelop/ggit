package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/ginx"
)

// repoAssignment extracts information from URL parameters to retrieve the repository,
// and makes sure the context user has at least the read access to the repository.
func repoAssignment() gin.HandlerFunc {
	return func(c *gin.Context) {
		actx := c.MustGet("apiContext").(*context.APIContext)
		username := c.Param("username")
		reponame := c.Param("reponame")

		var err error
		var owner *database.User

		if actx.IsLogged && actx.User.LowerName == strings.ToLower(username) {
			owner = actx.User
		} else {
			owner, err = database.Handle.Users().GetByUsername(c.Request.Context(), username)
			if err != nil {
				actx.NotFoundOrError(err, "get user by name")
				return
			}
		}
		actx.Repo.Owner = owner

		repo, err := database.Handle.Repositories().GetByName(c.Request.Context(), owner.ID, reponame)
		if err != nil {
			actx.NotFoundOrError(err, "get repository by name")
			return
		}

		if actx.IsTokenAuth && actx.User.IsAdmin {
			actx.Repo.AccessMode = database.AccessModeOwner
		} else {
			actx.Repo.AccessMode = database.Handle.Permissions().AccessMode(c.Request.Context(), actx.UserID(), repo.ID,
				database.AccessModeOptions{
					OwnerID: repo.OwnerID,
					Private: repo.IsPrivate,
				},
			)
		}

		if !actx.Repo.HasAccess() {
			actx.NotFound()
			return
		}

		actx.Repo.Repository = repo
		c.Next()
	}
}

func reqToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*context.Context)
		if !gctx.IsTokenAuth {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

func reqAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*context.Context)
		if !gctx.IsLogged || !gctx.User.IsAdmin {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RegisterRoutes(parent *gin.RouterGroup) {
	v1 := parent.Group("/v1")

	// handle routes using m.GET, m.POST, etc. with ginx.Handle
	v1.GET("/users/search", ginx.Handle(searchUsers))

	repos := v1.Group("/repos")
	{
		repos.GET("/search", ginx.Handle(searchRepos))

		repoBody := repos.Group("/:username/:reponame", repoAssignment())
		{
			repoBody.GET("", ginx.Handle(getRepo))
			repoBody.GET("/releases", ginx.Handle(releases))
		}
	}

	user := v1.Group("/user", reqToken())
	{
		user.GET("", ginx.Handle(getAuthenticatedUser))
	}
}
