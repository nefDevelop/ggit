package v2

import (
	"net/http"
	"strings"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
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

	repos, err := database.GetUserRepositories(user.ID, true)
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

	repos, err := database.GetStarredRepos(user.ID)
	if err != nil {
		c.Error(err, "get starred repositories")
		return
	}

	for _, repo := range repos {
		_ = repo.LoadAttributes()
	}

	c.JSON(http.StatusOK, repos)
}

// GetUserFollowers returns a list of followers for the given user.
func GetUserFollowers(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by username")
		return
	}

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	followers, err := database.Handle.Users().ListFollowers(c.Req.Context(), user.ID, page, 20)
	if err != nil {
		c.Error(err, "list followers")
		return
	}

	apiFollowers := make([]any, len(followers))
	for i, f := range followers {
		apiFollowers[i] = f.APIFormat()
	}

	c.JSON(http.StatusOK, apiFollowers)
}

// GetUserFollowing returns a list of users followed by the given user.
func GetUserFollowing(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by username")
		return
	}

	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	following, err := database.Handle.Users().ListFollowings(c.Req.Context(), user.ID, page, 20)
	if err != nil {
		c.Error(err, "list followings")
		return
	}

	apiFollowing := make([]any, len(following))
	for i, f := range following {
		apiFollowing[i] = f.APIFormat()
	}

	c.JSON(http.StatusOK, apiFollowing)
}
