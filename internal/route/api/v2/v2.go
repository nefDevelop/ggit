package v2

import (
	"github.com/gin-gonic/gin"
	
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/ginx"
)

// RegisterRoutes registers all v2 API routes to the given router group.
func RegisterRoutes(r *gin.RouterGroup) {
	r.Use(context.APIContexter())
	
	r.GET("/install/settings", ginx.Handle(GetInstallSettings))

	user := r.Group("/user")
	{
		user.POST("/login", ginx.Handle(Login))
		user.POST("/register", ginx.Handle(Register))
		
		user.GET("/check", ginx.Handle(CheckAuth))
		user.POST("/logout", ginx.Handle(Logout))
		user.POST("/login/two_factor", ginx.Handle(LoginTwoFactor))
		user.POST("/forgot_password", ginx.Handle(ForgotPassword))
		user.POST("/reset_password", ginx.Handle(ResetPassword))
		user.POST("/activate", ginx.Handle(Activate))

		settings := user.Group("/settings")
		{
			settings.GET("/profile", ginx.Handle(GetUserProfileSettings))
			settings.PUT("/profile", ginx.Handle(UpdateUserProfileSettings))
			settings.PUT("/password", ginx.Handle(UpdateUserPassword))

			settings.GET("/emails", ginx.Handle(GetUserEmails))
			settings.POST("/emails", ginx.Handle(AddUserEmail))
			settings.DELETE("/emails", ginx.Handle(DeleteUserEmail))
			settings.PUT("/emails/primary", ginx.Handle(MarkEmailPrimary))

			settings.GET("/ssh_keys", ginx.Handle(GetUserSSHKeys))
			settings.POST("/ssh_keys", ginx.Handle(AddUserSSHKey))
			settings.DELETE("/ssh_keys/:id", ginx.Handle(DeleteUserSSHKey))

			settings.GET("/tokens", ginx.Handle(GetUserAccessTokens))
			settings.POST("/tokens", ginx.Handle(AddUserAccessToken))
			settings.DELETE("/tokens/:id", ginx.Handle(DeleteUserAccessToken))

			settings.POST("/avatar", ginx.Handle(UpdateUserAvatar))
			settings.DELETE("/avatar", ginx.Handle(DeleteUserAvatar))

			settings.GET("/repos", ginx.Handle(GetUserRepositories))
			settings.DELETE("/repos/:id", ginx.Handle(LeaveUserRepository))

			settings.GET("/orgs", ginx.Handle(GetUserOrganizations))
			settings.DELETE("/orgs/:id", ginx.Handle(LeaveUserOrganization))

			settings.DELETE("/account", ginx.Handle(DeleteUserAccount))
		}

		dashboard := user.Group("/dashboard")
		{
			dashboard.GET("/data", ginx.Handle(GetDashboardData))
			dashboard.GET("/feeds", ginx.Handle(GetDashboardFeeds))
		}

		user.POST("/repos", ginx.Handle(CreateRepo))
		user.POST("/repos/migrate", ginx.Handle(MigrateRepo))
	}

	explore := r.Group("/explore")
	{
		explore.GET("/repos", ginx.Handle(ExploreRepos))
		explore.GET("/users", ginx.Handle(ExploreUsers))
		explore.GET("/organizations", ginx.Handle(ExploreOrganizations))
	}

	admin := r.Group("/admin", context.Toggle(&context.ToggleOptions{SignInRequired: true, AdminRequired: true}))
	{
		admin.GET("/stats", ginx.Handle(GetAdminStats))
		admin.GET("/config", ginx.Handle(GetAdminConfig))
		
		admin.GET("/users", ginx.Handle(GetAdminUsers))
		admin.PATCH("/users/:id", ginx.Handle(UpdateAdminUser))
		admin.DELETE("/users/:id", ginx.Handle(DeleteAdminUser))

		admin.GET("/repos", ginx.Handle(GetAdminRepos))
		admin.GET("/orgs", ginx.Handle(GetAdminOrgs))

		admin.GET("/auths", ginx.Handle(GetAdminAuths))
		admin.DELETE("/auths/:id", ginx.Handle(DeleteAdminAuth))

		admin.GET("/notices", ginx.Handle(GetAdminNotices))
		admin.DELETE("/notices/:id", ginx.Handle(DeleteAdminNotice))
		admin.DELETE("/notices", ginx.Handle(ClearAdminNotices))
	}

	repo := r.Group("/:username/:reponame", context.RepoAssignment(), context.RepoRef())
	{
		repo.GET("/view/*path", ginx.Handle(GetRepoView))
		repo.GET("/pulls", ginx.Handle(ListPullRequests))
		repo.GET("/commits/*path", ginx.Handle(GetRepoCommits))
		repo.GET("/watchers", ginx.Handle(GetRepoWatchers))
		repo.GET("/stars", ginx.Handle(GetRepoStars))
		repo.GET("/forks", ginx.Handle(GetRepoForks))

		issueGroup := repo.Group("/issues")
		{
			issueGroup.GET("", ginx.Handle(ListIssues))
			issueGroup.POST("", ginx.Handle(CreateIssue))
			issueGroup.GET("/:index", ginx.Handle(GetIssue))
			issueGroup.PUT("/:index", ginx.Handle(UpdateIssue))
			issueGroup.GET("/:index/comments", ginx.Handle(GetIssueComments))
			issueGroup.POST("/:index/comments", ginx.Handle(CreateIssueComment))
		}

		pullGroup := repo.Group("/pulls")
		{
			pullGroup.GET("", ginx.Handle(ListPullRequests))
			pullGroup.POST("", ginx.Handle(CreatePullRequest))
			pullGroup.GET("/:index", ginx.Handle(GetPullRequest))
			pullGroup.GET("/:index/diff", ginx.Handle(GetPullRequestDiff))
			pullGroup.POST("/:index/merge", ginx.Handle(MergePullRequest))
		}

		wikiGroup := repo.Group("/wiki")
		{
			wikiGroup.GET("", ginx.Handle(ListWikiPages))
			wikiGroup.POST("", ginx.Handle(CreateWikiPage))
			wikiGroup.GET("/:page", ginx.Handle(GetWikiPage))
			wikiGroup.GET("/:page/history", ginx.Handle(GetWikiPageHistory))
			wikiGroup.PUT("/:page", ginx.Handle(UpdateWikiPage))
			wikiGroup.DELETE("/:page", ginx.Handle(DeleteWikiPage))
		}

		releaseGroup := repo.Group("/releases")
		{
			releaseGroup.GET("", ginx.Handle(ListReleases))
			releaseGroup.POST("", ginx.Handle(CreateRelease))
			releaseGroup.GET("/:id", ginx.Handle(GetRelease))
			releaseGroup.PATCH("/:id", ginx.Handle(UpdateRelease))
			releaseGroup.DELETE("/:id", ginx.Handle(DeleteRelease))
		}

		repo.GET("/branches", ginx.Handle(ListRepoBranches))
		repo.DELETE("/branches/:branch", ginx.Handle(DeleteRepoBranch))
		
		repo.GET("/compare/*info", ginx.Handle(GetComparison))
		
		editor := repo.Group("/editor")
		{
			editor.POST("/edit", ginx.Handle(EditRepoFile))
			editor.POST("/delete", ginx.Handle(DeleteRepoFile))
		}

		settingGroup := repo.Group("/settings", context.RequireRepoAdmin())
		{
			settingGroup.GET("", ginx.Handle(GetRepoSettings))
			settingGroup.PATCH("", ginx.Handle(UpdateRepoSettings))
			settingGroup.DELETE("", ginx.Handle(DeleteRepository))

			settingGroup.GET("/collaborators", ginx.Handle(GetRepoCollaborators))
			settingGroup.POST("/collaborators/:collaborator", ginx.Handle(AddRepoCollaborator))
			settingGroup.DELETE("/collaborators/:collaborator", ginx.Handle(DeleteRepoCollaborator))

			settingGroup.GET("/hooks", ginx.Handle(GetRepoHooks))
			settingGroup.DELETE("/hooks/:id", ginx.Handle(DeleteRepoHook))
		}
	}

	orgs := r.Group("/orgs")
	{
		orgs.GET("/:org", ginx.Handle(GetOrg))
		orgs.PATCH("/:org", ginx.Handle(UpdateOrg))
		orgs.DELETE("/:org", ginx.Handle(DeleteOrg))

		orgs.GET("/:org/teams", ginx.Handle(ListOrgTeams))
		orgs.POST("/:org/teams", ginx.Handle(CreateOrgTeam))
		orgs.GET("/teams/:teamID", ginx.Handle(GetTeam))
		orgs.PATCH("/teams/:teamID", ginx.Handle(UpdateOrgTeam))
		orgs.DELETE("/teams/:teamID", ginx.Handle(DeleteOrgTeam))

		orgs.POST("/teams/:teamID/members/:username", ginx.Handle(AddTeamMember))
		orgs.DELETE("/teams/:teamID/members/:username", ginx.Handle(RemoveTeamMember))

		orgs.POST("/teams/:teamID/repos/:reponame", ginx.Handle(AddTeamRepo))
		orgs.DELETE("/teams/:teamID/repos/:repoID", ginx.Handle(RemoveTeamRepo))
	}

	users := r.Group("/users")
	{
		users.GET("/:username", ginx.Handle(GetUserProfile))
		users.GET("/:username/repos", ginx.Handle(GetUserPublicRepositories))
		users.GET("/:username/stars", ginx.Handle(GetUserStars))
		users.GET("/:username/followers", ginx.Handle(GetUserFollowers))
		users.GET("/:username/following", ginx.Handle(GetUserFollowing))
	}
}
