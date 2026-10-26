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
		user.GET("/login/settings", ginx.Handle(GetLoginSettings))
		user.POST("/login", ginx.Handle(Login))
		user.GET("/login/oidc/:id", ginx.Handle(LoginOIDC))
		user.GET("/login/oidc/:id/callback", ginx.Handle(LoginOIDCCallback))
		user.POST("/register", ginx.Handle(Register))

		user.GET("/check", ginx.Handle(CheckAuth))
		user.POST("/logout", ginx.Handle(Logout))
		user.POST("/login/two_factor", ginx.Handle(LoginTwoFactor))
		user.POST("/forgot_password", ginx.Handle(ForgotPassword))
		user.POST("/reset_password", ginx.Handle(ResetPassword))
		user.POST("/activate", ginx.Handle(Activate))

		settings := user.Group("/settings", context.Toggle(&context.ToggleOptions{SignInRequired: true}))
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

			settings.GET("/sessions", ginx.Handle(GetUserLoginSessions))
			settings.DELETE("/sessions", ginx.Handle(DeleteOtherUserLoginSessions))
			settings.DELETE("/sessions/:id", ginx.Handle(DeleteUserLoginSession))

			settings.GET("/blocks", ginx.Handle(GetBlockedUsers))
			settings.POST("/blocks/:username", ginx.Handle(BlockUser))
			settings.DELETE("/blocks/:username", ginx.Handle(UnblockUser))

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

		user.GET("/repos", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(GetUserRepositories))
		user.POST("/repos", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(CreateRepo))
		user.POST("/repos/migrate", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(MigrateRepo))

		dashboard := user.Group("/dashboard", context.Toggle(&context.ToggleOptions{SignInRequired: true}))
		{
			dashboard.GET("/data", ginx.Handle(GetDashboardData))
			dashboard.GET("/feeds", ginx.Handle(GetDashboardFeeds))
			dashboard.GET("/inbox", ginx.Handle(GetUnifiedInbox))
		}
	}

	// Global attachment download
	r.GET("/attachments/:uuid", context.APIContexter(), ginx.Handle(GetAttachment))

	explore := r.Group("/explore", context.Toggle(&context.ToggleOptions{SignInRequired: false}))
	{
		explore.GET("/repos", ginx.Handle(ExploreRepos))
		explore.GET("/repos/search", ginx.Handle(SearchRepos))
		explore.GET("/search", ginx.Handle(SearchRepos))
		explore.GET("/users", ginx.Handle(ExploreUsers))
		explore.GET("/organizations", ginx.Handle(ExploreOrganizations))
		explore.GET("/topics", ginx.Handle(ExploreTopics))
		explore.GET("/languages", ginx.Handle(ExploreLanguages))
	}

	admin := r.Group("/admin", context.Toggle(&context.ToggleOptions{SignInRequired: true, AdminRequired: true}))
	{
		admin.GET("/stats", ginx.Handle(GetAdminStats))
		admin.POST("/operations", ginx.Handle(RunAdminOperation))
		admin.GET("/config", ginx.Handle(GetAdminConfig))
		admin.PUT("/config", ginx.Handle(UpdateAdminConfig))
		admin.POST("/config/test_email", ginx.Handle(SendAdminTestEmail))
		admin.GET("/health-snapshot", ginx.Handle(GetAdminMonitor))
		admin.POST("/monitor/cron", ginx.Handle(RunAdminCronTask))
		admin.POST("/monitor/processes/:pid", ginx.Handle(KillAdminProcess))

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
		admin.GET("/repos/unmanaged", ginx.Handle(GetAdminUnmanagedRepos))
		admin.POST("/repos/adopt", ginx.Handle(AdoptAdminRepo))
		admin.POST("/repos/:id/check-integrity", ginx.Handle(CheckAdminRepoIntegrity))
		admin.GET("/backup/export", ginx.Handle(ExportAdminBackup))
		admin.GET("/janitor", ginx.Handle(GetAdminJanitorStats))
		admin.POST("/janitor/scan-all", ginx.Handle(ScanAllJanitor))

		admin.GET("/pikoci/status", ginx.Handle(GetAdminPikoCIStatus))
		admin.GET("/pikoci/repos", ginx.Handle(GetAdminPikoCIRepos))
		admin.GET("/pikoci/pipeline-template", ginx.Handle(GetAdminPikoCIPipelineTemplate))
		admin.PUT("/pikoci/repos/:repoID", ginx.Handle(UpdateAdminPikoCIRepo))
		admin.POST("/pikoci/repos/:repoID/sync", ginx.Handle(SyncAdminPikoCIRepo))
		admin.DELETE("/pikoci/repos/:repoID", ginx.Handle(DeleteAdminPikoCIRepo))

		admin.GET("/security/guardian/patterns", ginx.Handle(GetAdminGuardianPatterns))
		admin.POST("/security/guardian/patterns", ginx.Handle(AddAdminGuardianPattern))
		admin.DELETE("/security/guardian/patterns/:id", ginx.Handle(DeleteAdminGuardianPattern))
	}

	repo := r.Group("/:username/:reponame", context.Toggle(&context.ToggleOptions{SignInRequired: false}), context.RepoAssignment())
	{
		repo.GET("/search", context.RepoRef(), ginx.Handle(SearchRepoContents))
		repo.GET("/view", context.RepoRef(), ginx.Handle(GetRepoView))
		repo.GET("/view/*path", context.RepoRef(), ginx.Handle(GetRepoView))
		repo.GET("/contents", context.RepoRef(), ginx.Handle(GetRepoView))
		repo.GET("/contents/*path", context.RepoRef(), ginx.Handle(GetRepoView))
		repo.GET("/blame/*path", context.RepoRef(), ginx.Handle(GetRepoBlame))
		repo.GET("/commit/:sha", ginx.Handle(GetSingleCommit))
		repo.GET("/commits/*path", context.RepoRef(), ginx.Handle(GetRepoCommits))
		repo.GET("/watchers", ginx.Handle(GetRepoWatchers))
		repo.GET("/stars", ginx.Handle(GetRepoStars))
		repo.GET("/forks", ginx.Handle(GetRepoForks))
		repo.GET("/statuses/:sha", ginx.Handle(GetCommitStatuses))
		repo.POST("/statuses/:sha", context.RequireRepoWriter(), ginx.Handle(CreateCommitStatus))

		issueGroup := repo.Group("/issues")
		{
			issueGroup.GET("", ginx.Handle(ListIssues))
			issueGroup.POST("", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(CreateIssue))
			issueGroup.GET("/:index", ginx.Handle(GetIssue))
			issueGroup.PUT("/:index", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(UpdateIssue))
			issueGroup.GET("/:index/comments", ginx.Handle(GetIssueComments))
			issueGroup.POST("/:index/comments", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(CreateIssueComment))

			issueGroup.GET("/:index/labels", ginx.Handle(ListIssueLabels))
			issueGroup.POST("/:index/labels", context.RequireRepoWriter(), ginx.Handle(AddIssueLabels))
			issueGroup.PUT("/:index/labels", context.RequireRepoWriter(), ginx.Handle(ReplaceIssueLabels))
			issueGroup.DELETE("/:index/labels", context.RequireRepoWriter(), ginx.Handle(ClearIssueLabels))
		}

		pullGroup := repo.Group("/pulls")
		{
			pullGroup.GET("", ginx.Handle(ListPullRequests))
			pullGroup.POST("", context.RequireRepoWriter(), ginx.Handle(CreatePullRequest))
			pullGroup.GET("/:index", ginx.Handle(GetPullRequest))
			pullGroup.GET("/:index/diff", ginx.Handle(GetPullRequestDiff))
			pullGroup.POST("/:index/merge", context.RequireRepoWriter(), ginx.Handle(MergePullRequest))
			pullGroup.GET("/:index/protection", ginx.Handle(GetPullRequestProtection))
			pullGroup.GET("/:index/reviews", ginx.Handle(ListPullReviews))
			pullGroup.POST("/:index/reviews", context.RequireRepoWriter(), ginx.Handle(CreatePullReview))
			pullGroup.POST("/:index/reviews/:id/comments", context.RequireRepoWriter(), ginx.Handle(AddPullReviewComment))
			pullGroup.POST("/comments/:id/resolve", context.RequireRepoWriter(), ginx.Handle(ResolvePullReviewComment))
		}

		wikiGroup := repo.Group("/wiki")
		{
			wikiGroup.GET("", ginx.Handle(GetWikiPage))
			wikiGroup.GET("/*page", ginx.Handle(GetWikiPage))
			wikiGroup.POST("", context.RequireRepoWriter(), ginx.Handle(CreateWikiPage))
			wikiGroup.PUT("/:page", context.RequireRepoWriter(), ginx.Handle(UpdateWikiPage))
			wikiGroup.DELETE("/:page", context.RequireRepoWriter(), ginx.Handle(DeleteWikiPage))
		}

		releaseGroup := repo.Group("/releases")
		{
			releaseGroup.GET("", ginx.Handle(ListReleases))
			releaseGroup.POST("", context.RequireRepoWriter(), ginx.Handle(CreateRelease))
			releaseGroup.GET("/:id", ginx.Handle(GetRelease))
			releaseGroup.PATCH("/:id", context.RequireRepoWriter(), ginx.Handle(UpdateRelease))
			releaseGroup.DELETE("/:id", context.RequireRepoWriter(), ginx.Handle(DeleteRelease))
			releaseGroup.POST("/attachments", context.RequireRepoWriter(), ginx.Handle(UploadReleaseAttachment))
		}

		repo.GET("/branches", ginx.Handle(ListRepoBranches))
		repo.POST("/branches", context.RequireRepoWriter(), ginx.Handle(CreateRepoBranch))
		repo.DELETE("/branches/:branch", context.RequireRepoWriter(), ginx.Handle(DeleteRepoBranch))
		repo.GET("/branches/:branch/protection", ginx.Handle(GetBranchProtection))
		repo.PUT("/branches/:branch/protection", context.RequireRepoAdmin(), ginx.Handle(UpdateBranchProtection))
		repo.DELETE("/branches/:branch/protection", context.RequireRepoAdmin(), ginx.Handle(DeleteBranchProtection))

		repo.GET("/tags", ginx.Handle(ListRepoTags))
		repo.POST("/tags", context.RequireRepoWriter(), ginx.Handle(CreateRepoTag))
		repo.DELETE("/tags/:tag", context.RequireRepoWriter(), ginx.Handle(DeleteRepoTag))

		repo.GET("/archive/*name", context.RepoRef(), ginx.Handle(GetRepoArchive))
		repo.GET("/raw/:ref/*path", ginx.Handle(GetRepoRawFile))
		repo.GET("/media/:ref/*path", ginx.Handle(GetRepoMediaFile))

		repo.GET("/compare/*info", context.RepoRef(), ginx.Handle(GetComparison))

		repo.GET("/labels", ginx.Handle(ListRepoLabels))
		repo.POST("/labels", context.RequireRepoAdmin(), ginx.Handle(CreateRepoLabel))
		repo.GET("/labels/:id", ginx.Handle(GetRepoLabel))
		repo.PATCH("/labels/:id", context.RequireRepoAdmin(), ginx.Handle(UpdateRepoLabel))
		repo.DELETE("/labels/:id", context.RequireRepoAdmin(), ginx.Handle(DeleteRepoLabel))

		repo.GET("/milestones", ginx.Handle(ListRepoMilestones))
		repo.POST("/milestones", context.RequireRepoWriter(), ginx.Handle(CreateRepoMilestone))
		repo.GET("/milestones/:id", ginx.Handle(GetRepoMilestone))
		repo.PATCH("/milestones/:id", context.RequireRepoWriter(), ginx.Handle(UpdateRepoMilestone))
		repo.DELETE("/milestones/:id", context.RequireRepoWriter(), ginx.Handle(DeleteRepoMilestone))

		editor := repo.Group("/editor")
		{
			editor.POST("/edit", context.RequireRepoWriter(), context.RepoRef(), ginx.Handle(EditRepoFile))
			editor.POST("/delete", context.RequireRepoWriter(), context.RepoRef(), ginx.Handle(DeleteRepoFile))
		}

		settingGroup := repo.Group("/settings", context.RequireRepoAdmin())
		{
			settingGroup.GET("", ginx.Handle(GetRepoSettings))
			settingGroup.PATCH("", ginx.Handle(UpdateRepoSettings))
			settingGroup.POST("/avatar", ginx.Handle(UpdateRepoAvatar))
			settingGroup.DELETE("/avatar", ginx.Handle(DeleteRepoAvatar))
			settingGroup.POST("/transfer", ginx.Handle(TransferRepoOwnership))
			settingGroup.DELETE("", ginx.Handle(DeleteRepository))

			settingGroup.GET("/collaborators", ginx.Handle(GetRepoCollaborators))
			settingGroup.POST("/collaborators/:collaborator", ginx.Handle(AddRepoCollaborator))
			settingGroup.PATCH("/collaborators/:collaborator", ginx.Handle(UpdateRepoCollaborator))
			settingGroup.DELETE("/collaborators/:collaborator", ginx.Handle(DeleteRepoCollaborator))

			settingGroup.GET("/hooks", ginx.Handle(GetRepoHooks))
			settingGroup.POST("/hooks", ginx.Handle(CreateRepoHook))
			settingGroup.GET("/hooks/:id/history", ginx.Handle(GetRepoHookHistory))
			settingGroup.PATCH("/hooks/:id", ginx.Handle(UpdateRepoHook))
			settingGroup.POST("/hooks/:id/redeliver/:uuid", ginx.Handle(RedeliverRepoHookTask))
			settingGroup.POST("/hooks/:id/test", ginx.Handle(TestRepoHook))
			settingGroup.DELETE("/hooks/:id", ginx.Handle(DeleteRepoHook))

			settingGroup.GET("/keys", ginx.Handle(ListRepoDeployKeys))
			settingGroup.POST("/keys", ginx.Handle(AddRepoDeployKey))
			settingGroup.DELETE("/keys/:id", ginx.Handle(DeleteRepoDeployKey))

			settingGroup.GET("/remotes", ginx.Handle(ListRepoRemotes))
			settingGroup.POST("/remotes", ginx.Handle(AddRepoRemote))
			settingGroup.DELETE("/remotes/:name", ginx.Handle(DeleteRepoRemote))

			settingGroup.GET("/pikoci", ginx.Handle(GetRepoPikoCI))
			settingGroup.POST("/pikoci", ginx.Handle(EnableRepoPikoCI))
			settingGroup.DELETE("/pikoci", ginx.Handle(DisableRepoPikoCI))

			gitHooks := settingGroup.Group("/hooks/git", context.GitHookService())
			{
				gitHooks.GET("", ginx.Handle(ListRepoGitHooks))
				gitHooks.GET("/:name", ginx.Handle(GetRepoGitHook))
				gitHooks.PATCH("/:name", ginx.Handle(UpdateRepoGitHook))
			}

			guardian := settingGroup.Group("/security/guardian")
			{
				guardian.GET("", ginx.Handle(GetRepoGuardianSettings))
				guardian.PATCH("", ginx.Handle(UpdateRepoGuardianSettings))
				guardian.POST("/patterns", ginx.Handle(AddRepoGuardianPattern))
				guardian.DELETE("/patterns/:id", ginx.Handle(DeleteRepoGuardianPattern))
			}

			janitor := settingGroup.Group("/janitor")
			{
				janitor.GET("/report", ginx.Handle(GetJanitorReport))
				janitor.POST("/scan", ginx.Handle(RunJanitorScan))
				janitor.POST("/cleanup", ginx.Handle(PerformJanitorCleanup))
			}
		}
	}

	orgs := r.Group("/orgs", context.Toggle(&context.ToggleOptions{SignInRequired: false}))
	{
		// Lecturas: públicas para el perfil y la lista de miembros, restringidas a
		// miembros de la organización para el listado de equipos.
		orgs.GET("/:org", ginx.Handle(GetOrg))
		orgs.GET("/:org/members", ginx.Handle(ListOrgMembers))
		orgs.GET("/:org/teams", context.OrgAssignment(true), ginx.Handle(ListOrgTeams))

		// Modificación del perfil y de los miembros: reservadas a la propiedad de la
		// organización. OrgAssignment concede el acceso también a los administradores
		// del sitio, y responde 404 cuando no se cumple, para no revelar que la
		// organización existe.
		orgs.PATCH("/:org", context.OrgAssignment(true, true), ginx.Handle(UpdateOrg))
		orgs.POST("/:org/avatar", context.OrgAssignment(true, true), ginx.Handle(UpdateOrgAvatar))
		orgs.DELETE("/:org/avatar", context.OrgAssignment(true, true), ginx.Handle(DeleteOrgAvatar))
		orgs.DELETE("/:org", context.OrgAssignment(true, true), ginx.Handle(DeleteOrg))
		orgs.DELETE("/:org/members/:username", context.OrgAssignment(true, true), ginx.Handle(RemoveOrgMember))
		orgs.POST("/:org/teams", context.OrgAssignment(true, true), ginx.Handle(CreateOrgTeam))

		// Equipos. Estas rutas no llevan ":org", así que OrgAssignment no puede
		// resolver la organización y la comprobación se hace dentro del handler
		// mediante requireTeamAccess, que deduce el permiso desde el propio equipo.
		orgs.GET("/teams/:teamID", ginx.Handle(GetTeam))
		orgs.PATCH("/teams/:teamID", ginx.Handle(UpdateOrgTeam))
		orgs.DELETE("/teams/:teamID", ginx.Handle(DeleteOrgTeam))

		orgs.POST("/teams/:teamID/members/:username", ginx.Handle(AddTeamMember))
		orgs.DELETE("/teams/:teamID/members/:username", ginx.Handle(RemoveTeamMember))

		orgs.POST("/teams/:teamID/repos/:reponame", ginx.Handle(AddTeamRepo))
		orgs.DELETE("/teams/:teamID/repos/:repoID", ginx.Handle(RemoveTeamRepo))
	}

	users := r.Group("/users", context.Toggle(&context.ToggleOptions{SignInRequired: false}))
	{
		users.GET("/:username", ginx.Handle(GetUserProfile))
		users.GET("/:username/repos", ginx.Handle(GetUserPublicRepositories))
		users.GET("/:username/stars", ginx.Handle(GetUserStars))
		users.GET("/:username/orgs", ginx.Handle(GetUserPublicOrganizations))
		users.GET("/:username/gists", ginx.Handle(ListUserGists))
	}

	gist := r.Group("/gists")
	{
		gist.GET("", ginx.Handle(ListPublicGists))
		gist.POST("", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(CreateGist))
		gist.GET("/:uuid", ginx.Handle(GetGist))
		gist.PUT("/:uuid", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(UpdateGist))
		gist.DELETE("/:uuid", context.Toggle(&context.ToggleOptions{SignInRequired: true}), ginx.Handle(DeleteGist))
	}
}
