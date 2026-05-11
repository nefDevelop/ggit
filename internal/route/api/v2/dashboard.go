package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
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

	// 1. Get Personal, Collaborative and Watched Repositories
	repos, err := database.GetUserAndCollaborativeRepositories(ctxUser.ID)
	if err != nil {
		c.Error(err, "get user and collaborative repositories")
		return
	}

	// 2. Get Organizations
	orgs, err := database.GetOrgsByUserID(ctxUser.ID, true)
	if err != nil {
		log.Error("GetDashboardData: Failed to get organizations: %v", err)
	}

	// 3. Get Organization Repositories (explicitly)
	var orgRepos []*database.Repository
	for _, org := range orgs {
		ores, _, err := database.SearchRepositoryByName(&database.SearchRepoOptions{
			OwnerID:  org.ID,
			UserID:   ctxUser.ID,
			PageSize: 100,
		})
		if err != nil {
			log.Error("GetDashboardData: Failed to get repos for org %s: %v", org.Name, err)
			continue
		}
		orgRepos = append(orgRepos, ores...)
	}

	// 4. Merge and De-duplicate
	repoMap := make(map[int64]*database.Repository)
	for _, r := range repos {
		repoMap[r.ID] = r
	}
	for _, r := range orgRepos {
		repoMap[r.ID] = r
	}

	allRepos := make([]*database.Repository, 0, len(repoMap))
	for _, r := range repoMap {
		allRepos = append(allRepos, r)
	}

	// Sort by updated_unix DESC
	database.RepositoryList(allRepos).SortByUpdatedDesc()

	// Load attributes for all merged repos
	if err := database.RepositoryList(allRepos).LoadAttributes(); err != nil {
		log.Error("GetDashboardData: Failed to load attributes: %v", err)
	}

	apiRepos := make([]any, 0, len(allRepos))
	for _, r := range allRepos {
		apiRepos = append(apiRepos, r.APIFormat(r.Owner))
	}

	apiOrgs := make([]any, 0, len(orgs))
	for _, o := range orgs {
		apiOrgs = append(apiOrgs, o.APIFormat())
	}

	log.Info("GetDashboardData SUCCESS: User=%d, TotalRepos=%d", ctxUser.ID, len(apiRepos))

	c.JSON(http.StatusOK, map[string]any{
		"repositories":  apiRepos,
		"organizations": apiOrgs,
		"ok":            true,
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
