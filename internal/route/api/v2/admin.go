package v2

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"gogs.io/gogs/internal/auth"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/cron"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/email"
	"gogs.io/gogs/internal/process"
	apiv2types "gogs.io/gogs/internal/route/api/v2/types"
	"gogs.io/gogs/internal/tool"
)

// GetAdminAuditLog returns all actions from the system (for admins).
func GetAdminAuditLog(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	actions, err := database.Handle.Actions().ListAuditLog(c.Req.Context(), page)
	if err != nil {
		c.Error(err, "list audit log")
		return
	}

	resp := make([]*apiv2types.AuditLogResponse, len(actions))
	for i := range actions {
		resp[i] = ToAuditLogResponse(actions[i])
	}
	c.JSONSuccess(resp)
}

var initTime = time.Now()

// GetAdminStats returns global statistics and system status.
func GetAdminStats(c *context.APIContext) {
	stats := database.GetStatistic(c.Req.Context())

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	sysStatus := apiv2types.SystemStatusResponse{
		Uptime:       tool.TimeSincePro(initTime),
		NumGoroutine: runtime.NumGoroutine(),

		MemAllocated: tool.FileSize(int64(m.Alloc)),
		MemTotal:     tool.FileSize(int64(m.TotalAlloc)),
		MemSys:       tool.FileSize(int64(m.Sys)),
		Lookups:      m.Lookups,
		MemMallocs:   m.Mallocs,
		MemFrees:     m.Frees,

		HeapAlloc:    tool.FileSize(int64(m.HeapAlloc)),
		HeapSys:      tool.FileSize(int64(m.HeapSys)),
		HeapIdle:     tool.FileSize(int64(m.HeapIdle)),
		HeapInuse:    tool.FileSize(int64(m.HeapInuse)),
		HeapReleased: tool.FileSize(int64(m.HeapReleased)),
		HeapObjects:  m.HeapObjects,

		StackInuse:  tool.FileSize(int64(m.StackInuse)),
		StackSys:    tool.FileSize(int64(m.StackSys)),
		MSpanInuse:  tool.FileSize(int64(m.MSpanInuse)),
		MSpanSys:    tool.FileSize(int64(m.MSpanSys)),
		MCacheInuse: tool.FileSize(int64(m.MCacheInuse)),
		MCacheSys:   tool.FileSize(int64(m.MCacheSys)),
		BuckHashSys: tool.FileSize(int64(m.BuckHashSys)),
		GCSys:       tool.FileSize(int64(m.GCSys)),
		OtherSys:    tool.FileSize(int64(m.OtherSys)),

		NextGC:       tool.FileSize(int64(m.NextGC)),
		LastGC:       fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000),
		PauseTotalNs: fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000),
		PauseNs:      fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000),
		NumGC:        m.NumGC,
	}

	c.JSON(http.StatusOK, map[string]any{
		"data":   stats,
		"system": sysStatus,
	})
}

type AdminOperation int

const (
	CleanInactivateUser AdminOperation = iota + 1
	CleanRepoArchives
	CleanMissingRepos
	GitGCRepos
	SyncSSHAuthorizedKey
	SyncRepositoryHooks
	ReinitMissingRepository
	SyncRepositoryUpdate
	SyncGitFsck
)

// RunAdminOperation executes an administrative maintenance operation.
func RunAdminOperation(c *context.APIContext) {
	var f struct {
		Op int `json:"op" binding:"Required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	var err error
	var successKey string
	switch AdminOperation(f.Op) {
	case CleanInactivateUser:
		successKey = "admin.dashboard.delete_inactivate_accounts_success"
		err = database.Handle.Users().DeleteInactivated()
	case CleanRepoArchives:
		successKey = "admin.dashboard.delete_repo_archives_success"
		err = database.DeleteRepositoryArchives()
	case CleanMissingRepos:
		successKey = "admin.dashboard.delete_missing_repos_success"
		err = database.DeleteMissingRepositories()
	case GitGCRepos:
		successKey = "admin.dashboard.git_gc_repos_success"
		err = database.GitGcRepos()
	case SyncSSHAuthorizedKey:
		successKey = "admin.dashboard.resync_all_sshkeys_success"
		err = database.RewriteAuthorizedKeys()
	case SyncRepositoryHooks:
		successKey = "admin.dashboard.resync_all_hooks_success"
		err = database.SyncRepositoryHooks()
	case ReinitMissingRepository:
		successKey = "admin.dashboard.reinit_missing_repos_success"
		err = database.ReinitMissingRepositories()
	case SyncRepositoryUpdate:
		successKey = "admin.dashboard.resync_all_stats_success"
		database.CheckRepoStats()
	case SyncGitFsck:
		successKey = "admin.dashboard.git_fsck_success"
		database.GitFsck()
	default:
		c.ErrorStatus(http.StatusBadRequest, fmt.Errorf("Unknown operation ID: %d", f.Op))
		return
	}

	if err != nil {
		c.Error(err, "run admin operation")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"msg":    successKey,
	})
}

// GetAdminConfig returns the current application configuration.
func GetAdminConfig(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"app":        conf.App,
		"server":     conf.Server,
		"ssh":        conf.SSH,
		"repository": conf.Repository,
		"database":   conf.Database,
		"security":   conf.Security,
		"email":      conf.Email,
		"auth":       conf.Auth,
		"user":       conf.User,
		"session":    conf.Session,
		"cache":      conf.Cache,
		"attachment": conf.Attachment,
		"release":    conf.Release,
		"picture":    conf.Picture,
		"http":       conf.HTTP,
		"webhook":    conf.Webhook,
		"git":        conf.Git,
	})
}

// SendAdminTestEmail sends a test email to the specified address.
func SendAdminTestEmail(c *context.APIContext) {
	var f struct {
		Email string `json:"email" binding:"Required;Email;MaxSize(254)"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := email.SendTestMail(f.Email); err != nil {
		c.Error(err, "send test email")
		return
	}

	c.NoContent()
}

// GetAdminMonitor returns system monitor data.
func GetAdminMonitor(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"cron":      cron.ListTasks(),
		"processes": process.Processes,
	})
}

// GetAdminUsers returns a list of all users.
func GetAdminUsers(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	users, err := database.Handle.Users().List(c.Req.Context(), page, 50)
	if err != nil {
		c.Error(err, "list users")
		return
	}

	resp := make([]*apiv2types.UserResponse, len(users))
	for i, u := range users {
		resp[i] = ToUserResponse(u)
	}
	c.JSON(http.StatusOK, resp)
}

// CreateAdminUser creates a new user account.
func CreateAdminUser(c *context.APIContext) {
	var f struct {
		Username      string `json:"username" binding:"Required"`
		Email         string `json:"email" binding:"Required;Email;MaxSize(254)"`
		Password      string `json:"password" binding:"Required;MaxSize(255)"`
		FullName      string `json:"full_name"`
		Website       string `json:"website"`
		Location      string `json:"location"`
		Activated     bool   `json:"activated"`
		Admin         bool   `json:"admin"`
		ProhibitLogin bool   `json:"prohibit_login"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	user, err := database.Handle.Users().Create(c.Req.Context(), f.Username, f.Email, database.CreateUserOptions{
		Password:      f.Password,
		FullName:      f.FullName,
		Website:       f.Website,
		Location:      f.Location,
		Activated:     f.Activated,
		Admin:         f.Admin,
		ProhibitLogin: f.ProhibitLogin,
	})
	if err != nil {
		c.Error(err, "create user")
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(user))
}

// GetAdminUser returns a user's information.
func GetAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	user, err := database.Handle.Users().GetByID(c.Req.Context(), userID)
	if err != nil {
		c.Error(err, "get user")
		return
	}
	c.JSON(http.StatusOK, ToUserResponse(user))
}

// GetAdminRepos returns a list of all repositories.
func GetAdminRepos(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	repos, err := database.RepositoriesWithUsers(page, 50)
	if err != nil {
		c.Error(err, "list repositories")
		return
	}

	resp := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, r := range repos {
		resp[i] = ToRepositoryResponse(r)
	}
	c.JSON(http.StatusOK, resp)
}

// GetAdminOrgs returns a list of all organizations.
func GetAdminOrgs(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	orgs, err := database.Organizations(page, 50)
	if err != nil {
		c.Error(err, "list organizations")
		return
	}

	resp := make([]*apiv2types.OrganizationResponse, len(orgs))
	for i, o := range orgs {
		owner, _ := o.GetOrgOwner()
		resp[i] = ToOrganizationResponse(o, owner)
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteAdminRepo deletes a repository administratively.
func DeleteAdminRepo(c *context.APIContext) {
	repoID := c.ParamsInt64("id")
	repo, err := database.GetRepositoryByID(repoID)
	if err != nil {
		c.Error(err, "get repository")
		return
	}

	if err := database.DeleteRepository(repo.OwnerID, repo.ID); err != nil {
		c.Error(err, "delete repository")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// UpdateAdminUser updates a user's information administratively.
func UpdateAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	var f struct {
		FullName        *string `json:"full_name"`
		Website         *string `json:"website"`
		Location        *string `json:"location"`
		MaxRepoCreation *int    `json:"max_repo_creation"`
		IsActivated     *bool   `json:"activated"`
		IsAdmin         *bool   `json:"admin"`
		ProhibitLogin   *bool   `json:"prohibit_login"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	opts := database.UpdateUserOptions{
		FullName:        f.FullName,
		Website:         f.Website,
		Location:        f.Location,
		MaxRepoCreation: f.MaxRepoCreation,
		IsActivated:     f.IsActivated,
		IsAdmin:         f.IsAdmin,
		ProhibitLogin:   f.ProhibitLogin,
	}

	if err := database.Handle.Users().Update(c.Req.Context(), userID, opts); err != nil {
		c.Error(err, "update user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteAdminUser deletes a user account.
func DeleteAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	if err := database.Handle.Users().DeleteByID(c.Req.Context(), userID, false); err != nil {
		c.Error(err, "delete user")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetAdminAuths returns a list of all login sources.
func GetAdminAuths(c *context.APIContext) {
	sources, err := database.Handle.LoginSources().List(c.Req.Context(), database.ListLoginSourceOptions{})
	if err != nil {
		c.Error(err, "list login sources")
		return
	}

	resp := make([]*apiv2types.AuthSourceResponse, len(sources))
	for i, s := range sources {
		resp[i] = ToAuthSourceResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

// CreateAdminAuth creates a new login source.
// CreateAdminAuth creates a new authentication source.
func CreateAdminAuth(c *context.APIContext) {
	var f struct {
		Type      int    `json:"type" binding:"Required"`
		Name      string `json:"name" binding:"Required"`
		Activated bool   `json:"activated"`
		Default   bool   `json:"default"`
		Config    any    `json:"config" binding:"Required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	source, err := database.Handle.LoginSources().Create(c.Req.Context(), database.CreateLoginSourceOptions{
		Type:      auth.Type(f.Type),
		Name:      f.Name,
		Activated: f.Activated,
		Default:   f.Default,
		Config:    f.Config,
	})
	if err != nil {
		c.Error(err, "create login source")
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
		"id":     source.ID,
	})
}

// GetAdminAuth returns a login source's information.
func GetAdminAuth(c *context.APIContext) {
	authID := c.ParamsInt64("id")
	source, err := database.Handle.LoginSources().GetByID(c.Req.Context(), authID)
	if err != nil {
		c.Error(err, "get login source")
		return
	}
	c.JSON(http.StatusOK, ToAuthSourceResponse(source))
}

// UpdateAdminAuth updates an existing authentication source.
func UpdateAdminAuth(c *context.APIContext) {
	id := c.ParamsInt64("id")
	var f struct {
		Name      *string `json:"name"`
		Activated *bool   `json:"activated"`
		Default   *bool   `json:"default"`
		Config    any     `json:"config"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	opts := database.UpdateLoginSourceOptions{
		Name:      f.Name,
		Activated: f.Activated,
		Default:   f.Default,
		Config:    f.Config,
	}

	if err := database.Handle.LoginSources().Update(c.Req.Context(), id, opts); err != nil {
		c.Error(err, "update login source")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteAdminAuth deletes a login source.
func DeleteAdminAuth(c *context.APIContext) {
	authID := c.ParamsInt64("id")
	if err := database.Handle.LoginSources().DeleteByID(c.Req.Context(), authID); err != nil {
		c.Error(err, "delete login source")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetAdminNotices returns a list of system notices.
func GetAdminNotices(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	notices, err := database.Handle.Notices().List(c.Req.Context(), page, 50)
	if err != nil {
		c.Error(err, "list notices")
		return
	}

	resp := make([]*apiv2types.NoticeResponse, len(notices))
	for i, n := range notices {
		resp[i] = ToNoticeResponse(n)
	}

	c.JSON(http.StatusOK, map[string]any{
		"notices": resp,
		"total":   database.Handle.Notices().Count(c.Req.Context()),
	})
}

// DeleteAdminNotice deletes a system notice.
func DeleteAdminNotice(c *context.APIContext) {
	noticeID := c.ParamsInt64("id")
	if err := database.Handle.Notices().DeleteByIDs(c.Req.Context(), noticeID); err != nil {
		c.Error(err, "delete notice")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// ClearAdminNotices deletes all system notices.
func ClearAdminNotices(c *context.APIContext) {
	if err := database.Handle.Notices().DeleteAll(c.Req.Context()); err != nil {
		c.Error(err, "clear notices")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
