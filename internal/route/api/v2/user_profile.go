package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
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

	// One batch for the whole page: the fifty entries of a profile share the same
	// owner, so loading them individually meant the same user was fetched fifty
	// times.
	if err := database.RepositoryList(repos).LoadAttributes(); err != nil {
		c.Error(err, "load repository attributes")
		return
	}

	apiRepos := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, repo := range repos {
		if repo.Owner == nil {
			c.StatusError(http.StatusInternalServerError, "repository owner not found", nil)
			return
		}
		apiRepos[i] = repo.APIFormat(repo.Owner)
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
		drepo := database.FromDomainRepository(repo)
		if drepo == nil || drepo.Owner == nil {
			// The store resolves the owner while listing, so a nil here means the
			// data is inconsistent. MustOwner would panic and take the response down.
			c.StatusError(http.StatusInternalServerError, "repository owner not found", nil)
			return
		}
		apiRepos[i] = drepo.APIFormat(drepo.MustOwner())
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

	orgs, err := database.Handle.Organizations().List(c.Req.Context(), domain.ListOrgsOptions{
		MemberID:              user.ID,
		IncludePrivateMembers: c.IsLogged && (c.User.ID == user.ID || c.User.IsAdmin),
	})
	if err != nil {
		c.Error(err, "list user organizations")
		return
	}

	apiOrgs := make([]*apiv2types.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		apiOrgs[i] = database.FromDomainUser(org).APIFormatOrg()
	}

	c.JSON(http.StatusOK, apiOrgs)
}
