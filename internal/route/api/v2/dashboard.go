package v2

import (
	"net/http"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	apiv2types "gogs.io/gogs/internal/route/api/v2/types"
)

// DashboardData represents the consolidated data for the user dashboard.
type DashboardData struct {
	Repositories     []*apiv2types.RepositoryResponse `json:"repositories"`
	Collaborative    []*apiv2types.RepositoryResponse `json:"collaborative"`
	OrganizationList []*apiv2types.UserResponse       `json:"organizations"`
}

// GetDashboardData returns the initial data needed for the dashboard.
func GetDashboardData(c *context.APIContext) {
	ctxUser := c.User // For now, we only support personal dashboard in v2

	// 1. Get Personal Repositories
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
	apiRepos := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, r := range repos {
		apiRepos[i] = r.APIFormat(ctxUser)
	}

	// 2. Get Collaborative Repositories
	collaborative, err := database.Handle.Repositories().GetByCollaboratorID(c.Req.Context(), ctxUser.ID, conf.UI.User.RepoPagingNum, "updated_unix DESC")
	if err != nil {
		c.Error(err, "get collaborative repositories")
		return
	}
	_ = database.RepositoryList(collaborative).LoadAttributes()
	apiCollaborative := make([]*apiv2types.RepositoryResponse, len(collaborative))
	for i, r := range collaborative {
		apiCollaborative[i] = r.APIFormat(r.MustOwner())
	}

	// 4. Get Organizations
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
	apiOrgs := make([]*apiv2types.UserResponse, len(orgs))
	for i, o := range orgs {
		apiOrgs[i] = o.APIFormat()
	}

	c.JSON(http.StatusOK, DashboardData{
		Repositories:     apiRepos,
		Collaborative:    apiCollaborative,
		OrganizationList: apiOrgs,
	})
}

// GetDashboardFeeds returns the activity feed for the dashboard.
func GetDashboardFeeds(c *context.APIContext) {
	if c.User.DisableNewsFeed {
		c.JSON(http.StatusOK, map[string]any{
			"feeds": []any{},
		})
		return
	}

	ctxUser := c.User
	afterID := c.QueryInt64("after_id")

	actions, err := database.Handle.Actions().ListByUser(c.Req.Context(), ctxUser.ID, c.User.ID, afterID, false)
	if err != nil {
		c.Error(err, "list actions")
		return
	}

	// Filter for technical events: Pushes (ActionCommitRepo = 5) and Merges (ActionMergePullRequest = 11)
	technicalActions := make([]*database.Action, 0, len(actions))
	for _, act := range actions {
		if act.OpType == database.ActionCommitRepo || act.OpType == database.ActionMergePullRequest {
			technicalActions = append(technicalActions, act)
		}
	}
	actions = technicalActions

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
