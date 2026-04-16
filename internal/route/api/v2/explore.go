package v2

import (
	"net/http"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
)

// ExploreRepos returns a list of repositories for the explore page.
func ExploreRepos(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	keyword := c.Query("q")
	repos, count, err := database.SearchRepositoryByName(&database.SearchRepoOptions{
		Keyword:  keyword,
		UserID:   c.UserID(),
		OrderBy:  "updated_unix DESC",
		Page:     page,
		PageSize: conf.UI.ExplorePagingNum,
	})
	if err != nil {
		c.Error(err, "search repository by name")
		return
	}

	if err = database.RepositoryList(repos).LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"repos": repos,
		"total": count,
	})
}

// ExploreUsers returns a list of users for the explore page.
func ExploreUsers(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	keyword := c.Query("q")
	var (
		users []*database.User
		count int64
		err   error
	)
	if keyword == "" {
		users, err = database.Handle.Users().List(c.Req.Context(), page, conf.UI.ExplorePagingNum)
		if err != nil {
			c.Error(err, "list users")
			return
		}
		count = database.Handle.Users().Count(c.Req.Context())
	} else {
		users, count, err = database.Handle.Users().SearchByName(c.Req.Context(), keyword, page, conf.UI.ExplorePagingNum, "updated_unix DESC")
		if err != nil {
			c.Error(err, "search users by name")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"users": users,
		"total": count,
	})
}

// ExploreOrganizations returns a list of organizations for the explore page.
func ExploreOrganizations(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	keyword := c.Query("q")
	var (
		orgs  []*database.User
		count int64
		err   error
	)
	if keyword == "" {
		orgs, err = database.Organizations(page, conf.UI.ExplorePagingNum)
		if err != nil {
			c.Error(err, "list organizations")
			return
		}
		count = database.CountOrganizations()
	} else {
		orgs, count, err = database.Handle.Organizations().SearchByName(c.Req.Context(), keyword, page, conf.UI.ExplorePagingNum, "updated_unix DESC")
		if err != nil {
			c.Error(err, "search organizations by name")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"organizations": orgs,
		"total":         count,
	})
}
