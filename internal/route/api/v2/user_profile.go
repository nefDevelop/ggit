package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// GetUserProfile returns the public profile of a user.
func GetUserProfile(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by username")
		return
	}

	c.JSON(http.StatusOK, user.APIFormat())
}

// GetUserPublicRepositories returns a list of public repositories for the given user.
func GetUserPublicRepositories(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by username")
		return
	}

	repos, err := database.GetUserRepositories(&database.UserRepoOptions{
		UserID:   user.ID,
		Private:  c.IsLogged && (c.User.ID == user.ID || c.User.IsAdmin),
		PageSize: 50,
	})
	if err != nil {
		c.Error(err, "get user repositories")
		return
	}

	for _, repo := range repos {
		_ = repo.LoadAttributes()
	}

	c.JSON(http.StatusOK, repos)
}

// GetUserStars returns a list of repositories starred by the given user.
func GetUserStars(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by username")
		return
	}

	repos, err := database.Handle.Repositories().ListStarred(c.Req.Context(), user.ID, 1, 50)
	if err != nil {
		c.Error(err, "get starred repositories")
		return
	}

	for _, repo := range repos {
		_ = repo.LoadAttributes()
	}

	c.JSON(http.StatusOK, repos)
}
