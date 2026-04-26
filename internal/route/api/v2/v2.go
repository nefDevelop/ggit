package v2

import (
	"github.com/gin-gonic/gin"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/ginx"
)

// @title Ggit API v2
// @version 2.0
// @description REST API for Ggit, modernized to Go/Gin and Svelte.
// @host localhost:3000
// @BasePath /api/v2

// RegisterRoutes registers all v2 API routes to the given router group.
func RegisterRoutes(r *gin.RouterGroup) {
	r.Use(context.APIContexter())

	r.GET("/install/settings", ginx.Handle(GetInstallSettings))
	r.POST("/install", ginx.Handle(PerformInstall))
	r.GET("/i18n/:lang", ginx.Handle(GetI18nStrings))
	r.GET("/repo/templates", ginx.Handle(GetRepoTemplates))

	user := r.Group("/user")
	{
		user.GET("/registration_settings", ginx.Handle(GetRegistrationSettings))
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

			settings.GET("/security/two_factor", ginx.Handle(GetTwoFactorStatus))
			settings.POST("/security/two_factor/enable", ginx.Handle(EnableTwoFactor))
			settings.POST("/security/two_factor/confirm", ginx.Handle(ConfirmTwoFactor))
			settings.DELETE("/security/two_factor", ginx.Handle(DisableTwoFactor))
			settings.GET("/security/two_factor/recovery_codes", ginx.Handle(GetTwoFactorRecoveryCodes))
			settings.POST("/security/two_factor/recovery_codes", ginx.Handle(RegenerateTwoFactorRecoveryCodes))

			settings.GET("/ssh_keys", ginx.Handle(GetUserSSHKeys))
			settings.POST("/ssh_keys", ginx.Handle(AddUserSSHKey))
			settings.DELETE("/ssh_keys/:id", ginx.Handle(DeleteUserSSHKey))

			settings.GET("/tokens", ginx.Handle(GetUserAccessTokens))
			settings.POST("/tokens", ginx.Handle(AddUserAccessToken))
			settings.DELETE("/tokens/:id", ginx.Handle(DeleteUserAccessToken))

			settings.POST("/avatar", ginx.Handle(UpdateUserAvatar))
			settings.POST("/avatar/magic", ginx.Handle(UpdateUserAvatar))
			settings.DELETE("/avatar", ginx.Handle(DeleteUserAvatar))

			settings.GET("/repos", ginx.Handle(GetUserRepositories))
			settings.DELETE("/repos/:id", ginx.Handle(LeaveUserRepository))

			settings.GET("/orgs", ginx.Handle(GetUserOrganizations))
			settings.POST("/orgs", ginx.Handle(CreateUserOrganization))
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
		explore.GET("/repos/search", ginx.Handle(SearchRepos))
		explore.GET("/search", ginx.Handle(SearchRepos))
		explore.GET("/users", ginx.Handle(ExploreUsers))
		explore.GET("/organizations", ginx.Handle(ExploreOrganizations))
		explore.GET("/topics", ginx.Handle(ExploreTopics))
	}

	admin := r.Group("/admin", context.Toggle(&context.ToggleOptions{SignInRequired: true, AdminRequired: true}))
	{
		admin.GET("/stats", ginx.Handle(GetAdminStats))
		admin.POST("/operations", ginx.Handle(RunAdminOperation))
		admin.GET("/config", ginx.Handle(GetAdminConfig))
		admin.PUT("/config", ginx.Handle(UpdateAdminConfig))
		admin.POST("/config/test_email", ginx.Handle(SendAdminTestEmail))
		admin.GET("/monitor", ginx.Handle(GetAdminMonitor))

		admin.GET("/users", ginx.Handle(GetAdminUsers))
		admin.POST("/users", ginx.Handle(CreateAdminUser))
		admin.GET("/users/:id", ginx.Handle(GetAdminUser))
		admin.PATCH("/users/:id", ginx.Handle(UpdateAdminUser))
		admin.DELETE("/users/:id", ginx.Handle(DeleteAdminUser))

		admin.GET("/repos", ginx.Handle(GetAdminRepos))
		admin.DELETE("/repos/:id", ginx.Handle(DeleteAdminRepo))
		admin.GET("/orgs", ginx.Handle(GetAdminOrgs))
		admin.POST("/orgs", ginx.Handle(CreateAdminOrg))
		admin.GET("/auths", ginx.Handle(GetAdminAuths))
		admin.POST("/auths", ginx.Handle(CreateAdminAuth))
		admin.GET("/auths/:id", ginx.Handle(GetAdminAuth))
		admin.PATCH("/auths/:id", ginx.Handle(UpdateAdminAuth))
		admin.DELETE("/auths/:id", ginx.Handle(DeleteAdminAuth))

		admin.GET("/notices", ginx.Handle(GetAdminNotices))
		admin.DELETE("/notices/:id", ginx.Handle(DeleteAdminNotice))
		admin.DELETE("/notices", ginx.Handle(ClearAdminNotices))

		admin.GET("/audit-log", ginx.Handle(GetAdminAuditLog))
	}

	repo := r.Group("/:username/:reponame", context.RepoAssignment())
	{
		repo.GET("/search", context.RepoRef(), ginx.Handle(SearchRepoContents))
		repo.GET("/view/*path", context.RepoRef(), ginx.Handle(GetRepoView))
		repo.GET("/commit/:sha", context.RepoRef(), ginx.Handle(GetSingleCommit))
		repo.GET("/commits/*path", context.RepoRef(), ginx.Handle(GetRepoCommits))
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

			issueGroup.GET("/:index/labels", ginx.Handle(ListIssueLabels))
			issueGroup.POST("/:index/labels", ginx.Handle(AddIssueLabels))
			issueGroup.PUT("/:index/labels", ginx.Handle(ReplaceIssueLabels))
			issueGroup.DELETE("/:index/labels", ginx.Handle(ClearIssueLabels))
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
			wikiGroup.GET("", ginx.Handle(GetWikiPage))
			wikiGroup.GET("/*page", ginx.Handle(GetWikiPage))
			wikiGroup.POST("", ginx.Handle(CreateWikiPage))
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
		repo.POST("/branches", ginx.Handle(CreateRepoBranch))
		repo.DELETE("/branches/:branch", ginx.Handle(DeleteRepoBranch))

		repo.GET("/archive/*name", context.RepoRef(), ginx.Handle(GetRepoArchive))
		repo.GET("/raw/:ref/*path", ginx.Handle(GetRepoRawFile))
		repo.GET("/media/:ref/*path", ginx.Handle(GetRepoMediaFile))

		repo.GET("/compare/*info", context.RepoRef(), ginx.Handle(GetComparison))

		repo.GET("/labels", ginx.Handle(ListRepoLabels))
		repo.POST("/labels", ginx.Handle(CreateRepoLabel))
		repo.GET("/labels/:id", ginx.Handle(GetRepoLabel))
		repo.PATCH("/labels/:id", ginx.Handle(UpdateRepoLabel))
		repo.DELETE("/labels/:id", ginx.Handle(DeleteRepoLabel))

		repo.GET("/milestones", ginx.Handle(ListRepoMilestones))
		repo.POST("/milestones", ginx.Handle(CreateRepoMilestone))
		repo.GET("/milestones/:id", ginx.Handle(GetRepoMilestone))
		repo.PATCH("/milestones/:id", ginx.Handle(UpdateRepoMilestone))
		repo.DELETE("/milestones/:id", ginx.Handle(DeleteRepoMilestone))

		editor := repo.Group("/editor")
		{
			editor.POST("/edit", context.RepoRef(), ginx.Handle(EditRepoFile))
			editor.POST("/delete", context.RepoRef(), ginx.Handle(DeleteRepoFile))
		}

		settingGroup := repo.Group("/settings", context.RequireRepoAdmin())
		{
			settingGroup.GET("", ginx.Handle(GetRepoSettings))
			settingGroup.PATCH("", ginx.Handle(UpdateRepoSettings))
			settingGroup.POST("/transfer", ginx.Handle(TransferRepoOwnership))
			settingGroup.DELETE("", ginx.Handle(DeleteRepository))

			settingGroup.GET("/collaborators", ginx.Handle(GetRepoCollaborators))
			settingGroup.POST("/collaborators/:collaborator", ginx.Handle(AddRepoCollaborator))
			settingGroup.DELETE("/collaborators/:collaborator", ginx.Handle(DeleteRepoCollaborator))

			settingGroup.GET("/hooks", ginx.Handle(GetRepoHooks))
			settingGroup.DELETE("/hooks/:id", ginx.Handle(DeleteRepoHook))

			settingGroup.GET("/keys", ginx.Handle(ListRepoDeployKeys))
			settingGroup.POST("/keys", ginx.Handle(AddRepoDeployKey))
			settingGroup.DELETE("/keys/:id", ginx.Handle(DeleteRepoDeployKey))

			settingGroup.GET("/remotes", ginx.Handle(ListRepoRemotes))
			settingGroup.POST("/remotes", ginx.Handle(AddRepoRemote))
			settingGroup.DELETE("/remotes/:name", ginx.Handle(DeleteRepoRemote))
		}
	}


	orgs := r.Group("/orgs")
	{
		orgs.GET("/:org", ginx.Handle(GetOrg))
		orgs.PATCH("/:org", ginx.Handle(UpdateOrg))
		orgs.DELETE("/:org", ginx.Handle(DeleteOrg))
		orgs.GET("/:org/members", ginx.Handle(ListOrgMembers))
		orgs.DELETE("/:org/members/:username", ginx.Handle(RemoveOrgMember))

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
	}
}
