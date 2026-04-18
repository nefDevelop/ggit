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
	language := c.Query("lang")
	topics := c.QueryStrings("topic")
	repos, count, err := database.SearchRepositoryByName(&database.SearchRepoOptions{
		Keyword:  keyword,
		Language: language,
		Topics:   topics,
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

	if err = database.RepositoryList(repos).LoadTopics(c.Req.Context()); err != nil {
		c.Error(err, "load topics")
		return
	}

	apiRepos := make([]any, len(repos))
	for i, r := range repos {
		apiRepos[i] = r.APIFormat(r.Owner)
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":    true,
		"repos": apiRepos,
		"total": count,
	})
}

// SearchRepos returns a list of repositories based on search criteria.
func SearchRepos(c *context.APIContext) {
	opts := &database.SearchRepoOptions{
		Keyword:  c.Query("q"),
		OwnerID:  c.QueryInt64("uid"),
		Topics:   c.QueryStrings("topic"),
		Page:     c.QueryInt("page"),
		PageSize: c.QueryInt("limit"),
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 || opts.PageSize > 50 {
		opts.PageSize = 20
	}

	repos, count, err := database.SearchRepositoryByName(opts)
	if err != nil {
		c.Error(err, "search repository")
		return
	}

	if err = database.RepositoryList(repos).LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	if err = database.RepositoryList(repos).LoadTopics(c.Req.Context()); err != nil {
		c.Error(err, "load topics")
		return
	}

	apiRepos := make([]any, len(repos))
	for i, r := range repos {
		apiRepos[i] = r.APIFormat(r.Owner)
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":    true,
		"data":  apiRepos,
		"repos": apiRepos, // For consistency with ExploreRepos
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

	apiUsers := make([]any, len(users))
	for i, u := range users {
		apiUsers[i] = u.APIFormat()
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":    true,
		"users": apiUsers,
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

	apiOrgs := make([]any, len(orgs))
	for i, o := range orgs {
		apiOrgs[i] = o.APIFormat()
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":            true,
		"organizations": apiOrgs,
		"total":         count,
	})
}

// ExploreTopics returns a list of popular topics.
func ExploreTopics(c *context.APIContext) {
	topics, err := database.Handle.Repositories().ListPopularTopics(c.Req.Context(), 20)
	if err != nil {
		c.Error(err, "list popular topics")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":     true,
		"topics": topics,
	})
}
