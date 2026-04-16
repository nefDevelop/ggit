package v2

import (
	"github.com/gin-gonic/gin"
	
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/ginx"
)

// RegisterRoutes registers all v2 API routes to the given router group.
func RegisterRoutes(r *gin.RouterGroup) {
	r.Use(context.APIContexter())
	
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
		admin.GET("/repos", ginx.Handle(GetAdminRepos))
		admin.GET("/orgs", ginx.Handle(GetAdminOrgs))
	}

	repo := r.Group("/:username/:reponame", context.RepoAssignment(), context.RepoRef())
	{
		repo.GET("/view/*path", ginx.Handle(GetRepoView))
		repo.GET("/commits/*path", ginx.Handle(GetRepoCommits))
		repo.GET("/branches", ginx.Handle(GetRepoBranches))

		issueGroup := repo.Group("/issues")
		{
			issueGroup.GET("", ginx.Handle(ListIssues))
			issueGroup.POST("", ginx.Handle(CreateIssue))
			issueGroup.GET("/:index", ginx.Handle(GetIssue))
			issueGroup.PUT("/:index", ginx.Handle(UpdateIssue))
		}

		pullGroup := repo.Group("/pulls")
		{
			pullGroup.GET("", ginx.Handle(ListPullRequests))
			pullGroup.GET("/:index", ginx.Handle(GetPullRequest))
			pullGroup.POST("/:index/merge", ginx.Handle(MergePullRequest))
		}

		wikiGroup := repo.Group("/wiki")
		{
			wikiGroup.GET("", ginx.Handle(ListWikiPages))
			wikiGroup.POST("", ginx.Handle(CreateWikiPage))
			wikiGroup.GET("/:page", ginx.Handle(GetWikiPage))
			wikiGroup.PUT("/:page", ginx.Handle(UpdateWikiPage))
			wikiGroup.DELETE("/:page", ginx.Handle(DeleteWikiPage))
		}

		releaseGroup := repo.Group("/releases")
		{
			releaseGroup.GET("", ginx.Handle(ListReleases))
			releaseGroup.POST("", ginx.Handle(CreateRelease))
			releaseGroup.GET("/:id", ginx.Handle(GetRelease))
			releaseGroup.PUT("/:id", ginx.Handle(UpdateRelease))
			releaseGroup.DELETE("/:id", ginx.Handle(DeleteRelease))
		}

		settingGroup := repo.Group("/settings", context.RequireRepoAdmin())
		{
			settingGroup.GET("", ginx.Handle(GetRepoSettings))
			settingGroup.PUT("", ginx.Handle(UpdateRepoSettings))

			settingGroup.GET("/collaborators", ginx.Handle(GetRepoCollaborators))
			settingGroup.POST("/collaborators/:collaborator", ginx.Handle(AddRepoCollaborator))
			settingGroup.DELETE("/collaborators/:collaborator", ginx.Handle(DeleteRepoCollaborator))

			settingGroup.GET("/hooks", ginx.Handle(GetRepoHooks))
			settingGroup.DELETE("/hooks/:id", ginx.Handle(DeleteRepoHook))
		}
	}
}
