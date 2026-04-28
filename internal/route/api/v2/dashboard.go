package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	log "unknwon.dev/clog/v2"
)

// DashboardData represents the consolidated data for the user dashboard.
type DashboardData struct {
	Repositories     []*apiv2types.RepositoryResponse `json:"repositories"`
	Collaborative    []*apiv2types.RepositoryResponse `json:"collaborative"`
	OrganizationList []*apiv2types.UserResponse       `json:"organizations"`
}

// GetDashboardData returns the initial data needed for the dashboard.
func GetDashboardData(c *context.APIContext) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Panic in GetDashboardData: %v", err)
			c.StatusError(http.StatusInternalServerError, "Internal Server Error", nil)
		}
	}()

	ctxUser := c.User
	if ctxUser == nil {
		c.StatusError(http.StatusUnauthorized, "User must be logged in", nil)
		return
	}

	// 1. Get Personal Repositories
	log.Trace("GetDashboardData: Fetching personal repositories for user %d", ctxUser.ID)
	repos, err := database.GetUserRepositories(
		&database.UserRepoOptions{
			UserID:   ctxUser.ID,
			Private:  true,
			Page:     1,
			PageSize: conf.UI.User.RepoPagingNum,
		},
	)
	if err != nil {
		c.Error(err, "get user repositories")
		return
	}
	apiRepos := make([]*apiv2types.RepositoryResponse, 0, len(repos))
	for _, r := range repos {
		if r == nil {
			continue
		}
		apiRepos = append(apiRepos, r.APIFormat(ctxUser))
	}

	// 2. Get Collaborative Repositories
	log.Trace("GetDashboardData: Fetching collaborative repositories for user %d", ctxUser.ID)
	collaborative, err := database.Handle.Repositories().GetByCollaboratorID(c.Req.Context(), ctxUser.ID, conf.UI.User.RepoPagingNum, "updated_unix DESC")
	if err != nil {
		c.Error(err, "get collaborative repositories")
		return
	}
	log.Trace("GetDashboardData: Loading attributes for %d collaborative repos", len(collaborative))
	_ = database.RepositoryList(collaborative).LoadAttributes()
	
	apiCollaborative := make([]*apiv2types.RepositoryResponse, 0, len(collaborative))
	for _, r := range collaborative {
		if r == nil {
			continue
		}
		apiCollaborative = append(apiCollaborative, r.APIFormat(r.MustOwner()))
	}

	// 4. Get Organizations
	log.Trace("GetDashboardData: Fetching organizations for user %d", ctxUser.ID)
	orgs, err := database.Handle.Organizations().List(
		c.Req.Context(),
		database.ListOrgsOptions{
			MemberID:              ctxUser.ID,
			IncludePrivateMembers: true,
		},
	)
	if err != nil {
		c.Error(err, "list organizations")
		return
	}
	apiOrgs := make([]*apiv2types.UserResponse, 0, len(orgs))
	for _, o := range orgs {
		if o == nil {
			continue
		}
		apiOrgs = append(apiOrgs, o.APIFormat())
	}

	log.Trace("GetDashboardData: Returning consolidated data")
	c.JSON(http.StatusOK, DashboardData{
		Repositories:     apiRepos,
		Collaborative:    apiCollaborative,
		OrganizationList: apiOrgs,
	})
}

// GetDashboardFeeds returns the activity feed for the dashboard.
func GetDashboardFeeds(c *context.APIContext) {
	ctxUser := c.User
	if ctxUser == nil {
		c.StatusError(http.StatusUnauthorized, "User must be logged in", nil)
		return
	}

	if ctxUser.DisableNewsFeed {
		c.JSON(http.StatusOK, map[string]any{
			"feeds": []any{},
		})
		return
	}
	afterID := c.QueryInt64("after_id")

	actions, err := database.Handle.Actions().ListByUser(c.Req.Context(), ctxUser.ID, c.User.ID, afterID, false)
	if err != nil {
		c.Error(err, "list actions")
		return
	}

	// Prepare avatars for the feed
	unameAvatars := make(map[string]string)
	for _, act := range actions {
		if _, ok := unameAvatars[act.ActUserName]; !ok {
			u, err := database.Handle.Users().GetByUsername(c.Req.Context(), act.ActUserName)
			if err == nil {
				unameAvatars[act.ActUserName] = u.AvatarURLPath()
			}
		}
		act.ActAvatar = unameAvatars[act.ActUserName]
	}

	c.JSON(http.StatusOK, map[string]any{
		"feeds": actions,
	})
}
