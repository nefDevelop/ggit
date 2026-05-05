package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// GetRepoWatchers returns a list of users watching the repository.
func GetRepoWatchers(c *context.APIContext) {
	users, err := database.Handle.Repositories().ListWatches(c.Req.Context(), c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "list watches")
		return
	}

	apiUsers := make([]any, len(users))
	for i, u := range users {
		apiUsers[i] = database.FromDomainUser(u).APIFormat()
	}

	c.JSON(http.StatusOK, apiUsers)
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
		apiUsers[i] = database.FromDomainUser(u).APIFormat()
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
		_ = database.FromDomainRepository(fork).LoadAttributes()
	}

	c.JSON(http.StatusOK, forks)
}
