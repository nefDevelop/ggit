package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// GetUserProfile returns the public profile of a user.
func GetUserProfile(c *context.APIContext) {
	username := c.Params("username")
	du, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}
	user := database.FromDomainUser(du)


	c.JSON(http.StatusOK, user.APIFormat())
}

// GetUserPublicRepositories returns a list of public repositories for the given user.
func GetUserPublicRepositories(c *context.APIContext) {
	username := c.Params("username")
	du, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}
	user := database.FromDomainUser(du)


	repos, err := database.GetUserRepositories(&database.UserRepoOptions{
		UserID:   user.ID,
		Private:  c.IsLogged && (c.User.ID == user.ID || c.User.IsAdmin),
		PageSize: 50,
	})
	if err != nil {
		c.Error(err, "get user repositories")
		return
	}

	apiRepos := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, repo := range repos {
		_ = repo.LoadAttributes()
		apiRepos[i] = repo.APIFormat(repo.MustOwner())
	}

	c.JSON(http.StatusOK, apiRepos)
}

// GetUserStars returns a list of repositories starred by the given user.
func GetUserStars(c *context.APIContext) {
	username := c.Params("username")
	du, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}
	user := database.FromDomainUser(du)


	repos, err := database.Handle.Repositories().ListStarred(c.Req.Context(), user.ID, 1, 50)
	if err != nil {
		c.Error(err, "get starred repositories")
		return
	}

	apiRepos := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, repo := range repos {
		_ = repo.LoadAttributes()
		apiRepos[i] = repo.APIFormat(repo.MustOwner())
	}

	c.JSON(http.StatusOK, apiRepos)
}

// GetUserPublicOrganizations returns a list of public organizations for the given user.
func GetUserPublicOrganizations(c *context.APIContext) {
	username := c.Params("username")
	du, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}
	user := database.FromDomainUser(du)


	orgs, err := database.Handle.Organizations().List(c.Req.Context(), database.ListOrgsOptions{
		MemberID:              user.ID,
		IncludePrivateMembers: c.IsLogged && (c.User.ID == user.ID || c.User.IsAdmin),
	})
	if err != nil {
		c.Error(err, "list user organizations")
		return
	}

	apiOrgs := make([]*apiv2types.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		apiOrgs[i] = org.APIFormatOrg()
	}

	c.JSON(http.StatusOK, apiOrgs)
}
