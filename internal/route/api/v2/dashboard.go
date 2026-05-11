package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	log "github.com/nef734/ggit/internal/log"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
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

	// 1. Get All Accessible Repositories (Personal + Org + Collab)
	log.Trace("GetDashboardData: Fetching all accessible repositories for user %d", ctxUser.ID)
	repos, _, err := database.SearchRepositoryByName(
		&database.SearchRepoOptions{
			UserID:   ctxUser.ID,
			Page:     1,
			PageSize: conf.UI.User.RepoPagingNum,
			OrderBy:  "updated_unix DESC",
		},
	)
	if err != nil {
		c.Error(err, "search user repositories")
		return
	}
	
	// Load attributes for all found repos
	if err := database.RepositoryList(repos).LoadAttributes(); err != nil {
		log.Error("GetDashboardData: Failed to load attributes: %v", err)
	}

	apiRepos := make([]*apiv2types.RepositoryResponse, 0, len(repos))
	for _, r := range repos {
		if r == nil {
			continue
		}
		apiRepos = append(apiRepos, r.APIFormat(r.MustOwner()))
	}

	// 2. Get Collaborative Repositories
	log.Trace("GetDashboardData: Fetching collaborative repositories for user %d", ctxUser.ID)
	collaborative, err := database.Handle.Repositories().GetByCollaboratorID(c.Req.Context(), ctxUser.ID, conf.UI.User.RepoPagingNum, "updated_unix DESC")
	if err != nil {
		c.Error(err, "get collaborative repositories")
		return
	}
	log.Trace("GetDashboardData: Loading attributes for %d collaborative repos", len(collaborative))
	drepos := make([]*database.Repository, len(collaborative))
	for i, r := range collaborative {
		drepos[i] = database.FromDomainRepository(r)
	}
	_ = database.RepositoryList(drepos).LoadAttributes()

	apiCollaborative := make([]*apiv2types.RepositoryResponse, 0, len(collaborative))
	for i, r := range collaborative {
		if r == nil {
			continue
		}
		dr := drepos[i]
		apiCollaborative = append(apiCollaborative, dr.APIFormat(dr.MustOwner()))
	}

	// 4. Get Organizations
	log.Trace("GetDashboardData: Fetching organizations for user %d", ctxUser.ID)
	orgs, err := database.Handle.Organizations().List(
		c.Req.Context(),
		domain.ListOrgsOptions{
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
		apiOrgs = append(apiOrgs, database.FromDomainUser(o).APIFormat())
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
				unameAvatars[act.ActUserName] = database.FromDomainUser(u).AvatarURLPath()
			}
		}
		act.ActAvatar = unameAvatars[act.ActUserName]
	}

	c.JSON(http.StatusOK, map[string]any{
		"feeds": actions,
	})
}

// GetUnifiedInbox returns a list of assigned and relevant issues/PRs for the user.
func GetUnifiedInbox(c *context.APIContext) {
	ctxUser := c.User

	// 1. Get Assigned Issues and PRs
	var assigned []*database.Issue
	err := database.Handle.DB().Where("assignee_id = ? AND is_closed = ?", ctxUser.ID, false).
		Order("updated_unix DESC").Limit(50).Find(&assigned).Error
	if err != nil {
		c.Error(err, "find assigned items")
		return
	}

	// Load attributes (Repo, Poster, etc.)
	for _, is := range assigned {
		_ = is.LoadAttributes()
	}

	type InboxItem struct {
		ID           int64  `json:"id"`
		RepoName     string `json:"repo_name"`
		RepoOwner    string `json:"repo_owner"`
		Index        int64  `json:"index"`
		Title        string `json:"title"`
		IsPull       bool   `json:"is_pull"`
		Status       string `json:"status"`
		UpdatedUnix  int64  `json:"updated_at"`
		PosterName   string `json:"poster_name"`
		PosterAvatar string `json:"poster_avatar"`
	}

	results := make([]InboxItem, 0, len(assigned))
	for _, is := range assigned {
		repo := is.Repo
		ownerName := ""
		if repo != nil {
			ownerName = repo.MustOwner().Name
		}

		results = append(results, InboxItem{
			ID:           is.ID,
			RepoName:     repo.Name,
			RepoOwner:    ownerName,
			Index:        is.Index,
			Title:        is.Title,
			IsPull:       is.IsPull,
			Status:       "open",
			UpdatedUnix:  is.UpdatedUnix,
			PosterName:   is.Poster.Name,
			PosterAvatar: is.Poster.AvatarURLPath(),
		})
	}

	c.JSON(http.StatusOK, results)
}
