package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
)

// ExploreRepos returns a list of repositories for the explore page.
// @Summary Explore Repositories
// @Description Returns a list of public repositories with pagination and filtering.
// @Tags explore
// @Produce json
// @Param page query int false "Page number"
// @Param q query string false "Search keyword"
// @Param lang query string false "Language filter"
// @Param topic query []string false "Topic filter"
// @Success 200 {object} map[string]any
// @Router /explore/repos [get]
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
// @Summary Search Repositories
// @Description Searches for repositories globally based on various criteria.
// @Tags explore
// @Produce json
// @Param q query string false "Search keyword"
// @Param uid query int64 false "Owner user ID"
// @Param topic query []string false "Topic filter"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]any
// @Router /explore/repos/search [get]
// @Router /explore/search [get]
// @Router /repos/search [get]
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
		var dus []*domain.User
		dus, err = database.Handle.Users().List(c.Req.Context(), page, conf.UI.ExplorePagingNum)
		if err != nil {
			c.Error(err, "list users")
			return
		}
		users = make([]*database.User, len(dus))
		for i, du := range dus {
			users[i] = database.FromDomainUser(du)
		}
		count = database.Handle.Users().Count(c.Req.Context())
	} else {
		dus, cnt, err := database.Handle.Users().SearchByName(c.Req.Context(), keyword, page, conf.UI.ExplorePagingNum, "updated_unix DESC")
		if err != nil {
			c.Error(err, "search users by name")
			return
		}
		users = make([]*database.User, len(dus))
		for i, du := range dus {
			users[i] = database.FromDomainUser(du)
		}
		count = cnt
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
		var dus []*domain.User
		dus, count, err = database.Handle.Organizations().SearchByName(c.Req.Context(), keyword, page, conf.UI.ExplorePagingNum, "updated_unix DESC")
		if err != nil {
			c.Error(err, "search organizations by name")
			return
		}
		orgs = make([]*database.User, len(dus))
		for i, du := range dus {
			orgs[i] = database.FromDomainUser(du)
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
	if !database.HasEngine {
		c.JSON(http.StatusOK, map[string]any{"ok": true, "topics": []any{}})
		return
	}

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

// ExploreLanguages returns a list of unique languages from public repositories.
func ExploreLanguages(c *context.APIContext) {
	langs, err := database.ListPublicLanguages()
	if err != nil {
		c.Error(err, "list public languages")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"ok":        true,
		"languages": langs,
	})
}
