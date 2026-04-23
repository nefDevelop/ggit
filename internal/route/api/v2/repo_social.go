package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// GetRepoWatchers returns a list of users watching the repository.
func GetRepoWatchers(c *context.APIContext) {
	watches, err := database.Handle.Repositories().ListWatches(c.Req.Context(), c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "list watches")
		return
	}

	users := make([]any, 0, len(watches))
	for _, w := range watches {
		u, err := database.Handle.Users().GetByID(c.Req.Context(), w.UserID)
		if err != nil {
			continue
		}
		users = append(users, u.APIFormat())
	}

	c.JSON(http.StatusOK, users)
}

// GetRepoStars returns a list of users who starred the repository.
func GetRepoStars(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	stars, err := database.Handle.Repositories().ListStars(c.Req.Context(), c.Repo.Repository.ID, page, 20)
	if err != nil {
		c.Error(err, "list stars")
		return
	}

	apiUsers := make([]any, len(stars))
	for i, u := range stars {
		apiUsers[i] = u.APIFormat()
	}

	c.JSON(http.StatusOK, apiUsers)
}

// GetRepoForks returns a list of repositories that are forks of this repository.
func GetRepoForks(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	forks, err := database.Handle.Repositories().ListForks(c.Req.Context(), c.Repo.Repository.ID, page, 20)
	if err != nil {
		c.Error(err, "list forks")
		return
	}

	for _, fork := range forks {
		_ = fork.LoadAttributes()
	}

	c.JSON(http.StatusOK, forks)
}
