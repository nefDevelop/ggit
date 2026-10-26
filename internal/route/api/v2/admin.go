package v2

import (
	json "github.com/nef734/ggit/internal/jsonx"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/cron"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/email"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/process"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/nef734/ggit/internal/services/backup"
	"github.com/nef734/ggit/internal/services/health"
	"github.com/nef734/ggit/internal/services/janitor"
	"github.com/nef734/ggit/internal/services/repo"
	"github.com/nef734/ggit/internal/tool"
	"gopkg.in/ini.v1"
)

// @Summary Get Admin Audit Log
// @Description Returns all actions from the system (for admins).
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {array} apiv2types.AuditLogResponse
// @Router /admin/audit-log [get]
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
		resp[i] = ToAuditLogResponse(database.FromDomainAction(actions[i]))
	}
	c.JSONSuccess(resp)
}

var initTime = time.Now()

// @Summary Get Admin Statistics
// @Description Returns global statistics and system status.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]any
// @Router /admin/stats [get]
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

type adminOperationPayload struct {
	Op int `json:"op" binding:"required"`
}

// @Summary Run Admin Operation
// @Description Executes an administrative maintenance operation.
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body adminOperationPayload true "Operation details"
// @Success 200 {object} map[string]any
// @Router /admin/operations [post]
// RunAdminOperation executes an administrative maintenance operation.
func RunAdminOperation(c *context.APIContext) {
	var f adminOperationPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	var err error
	var successKey string
	switch AdminOperation(f.Op) {
	case CleanInactivateUser:
		successKey = "admin.dashboard.delete_inactivate_accounts_success"
		err = database.Handle.Users().DeleteInactivated(c.Req.Context())
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

// @Summary Get Admin Configuration
// @Description Returns the current application configuration.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]map[string]string
// @Router /admin/config [get]
// GetAdminConfig returns the current application configuration.
func GetAdminConfig(c *context.APIContext) {
	sections := conf.File.Sections()
	resp := make(map[string]map[string]string)
	for _, s := range sections {
		if s.Name() == ini.DefaultSection {
			resp["app"] = s.KeysHash()
			continue
		}
		resp[s.Name()] = s.KeysHash()
	}
	c.JSON(http.StatusOK, resp)
}

// SendAdminTestEmail sends a test email to the specified address.
func SendAdminTestEmail(c *context.APIContext) {
	var f struct {
		Email string `json:"email" binding:"required,email,max=254"`
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

// @Summary Get Admin Monitor
// @Description Returns system monitor data.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]any
// @Router /admin/health-snapshot [get]
// GetAdminMonitor returns system monitor data.
func GetAdminMonitor(c *context.APIContext) {
	log.Trace("API: GetAdminMonitor called")
	tasks := cron.ListTasks()
	type cronTask struct {
		Description string `json:"Description"`
		Spec        string `json:"Spec"`
		Next        string `json:"Next"`
		Prev        string `json:"Prev"`
		ExecTimes   int64  `json:"ExecTimes"`
	}
	cronTasks := make([]cronTask, len(tasks))
	for i, t := range tasks {
		cronTasks[i] = cronTask{
			Description: t.Description,
			Spec:        t.Spec,
			Next:        t.Next.Format(time.RFC3339),
			Prev:        t.Prev.Format(time.RFC3339),
			ExecTimes:   int64(t.ExecTimes),
		}
	}

	type procInfo struct {
		PID         int64  `json:"PID"`
		Description string `json:"Description"`
		Start       string `json:"Start"`
	}
	procInfos := make([]procInfo, len(process.Processes))
	for i, p := range process.Processes {
		procInfos[i] = procInfo{
			PID:         p.PID,
			Description: p.Description,
			Start:       p.Start.Format(time.RFC3339),
		}
	}

	log.Trace("API: Sending monitor response (Tasks: %d, Procs: %d)", len(cronTasks), len(procInfos))
	
	respData := map[string]any{
		"ping":        "pong",
		"cron_v2":     cronTasks,
		"processes":   procInfos,
		"debug_info": map[string]any{
			"install_lock": conf.Security.InstallLock,
			"custom_conf":  conf.CustomConf,
			"cron_enabled": conf.Cron.RepoHealthCheck.Enabled,
			"tasks_count":  len(tasks),
			"v":           5,
		},
	}
	
	c.JSON(http.StatusOK, respData)
	log.Trace("API: Monitor response sent (v5)")
}

// RunAdminCronTask runs a cron task manually.
func RunAdminCronTask(c *context.APIContext) {
	var f struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	cron.RunTask(f.Description)
	c.NoContent()
}

// KillAdminProcess kills a running process.
func KillAdminProcess(c *context.APIContext) {
	pid := c.ParamsInt64(":pid")
	if err := process.Kill(pid); err != nil {
		c.Error(err, "kill process")
		return
	}
	c.NoContent()
}

// @Summary Get Admin Users
// @Description Returns a list of all users.
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {array} apiv2types.UserResponse
// @Router /admin/users [get]
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
		resp[i] = ToUserResponse(database.FromDomainUser(u))
	}
	c.JSON(http.StatusOK, resp)
}

type adminCreateUserPayload struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email,max=254"`
	Password         string `json:"password" binding:"required,max=255"`
	FullName         string `json:"full_name"`
	Website          string `json:"website"`
	Location         string `json:"location"`
	Activated        bool   `json:"activated"`
	Admin            bool   `json:"admin"`
	ProhibitLogin    bool   `json:"prohibit_login"`
	AllowGitHook     bool   `json:"allow_git_hook"`
	AllowImportLocal bool   `json:"allow_import_local"`
}

// @Summary Create User (Admin)
// @Description Creates a new user account administratively.
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body adminCreateUserPayload true "User details"
// @Success 201 {object} apiv2types.UserResponse
// @Router /admin/users [post]
// CreateAdminUser creates a new user account.
func CreateAdminUser(c *context.APIContext) {
	var f adminCreateUserPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	user, err := database.Handle.Users().Create(c.Req.Context(), f.Username, f.Email, database.CreateUserOptions{
		Password:         f.Password,
		FullName:         f.FullName,
		Website:          f.Website,
		Location:         f.Location,
		Activated:        f.Activated,
		Admin:            f.Admin,
		ProhibitLogin:    f.ProhibitLogin,
		AllowGitHook:     f.AllowGitHook,
		AllowImportLocal: f.AllowImportLocal,
	})
	if err != nil {
		c.Error(err, "create user")
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(database.FromDomainUser(user)))
}

// @Summary Get User (Admin)
// @Description Returns a user's information.
// @Tags admin
// @Produce json
// @Param id path int64 true "User ID"
// @Success 200 {object} apiv2types.UserResponse
// @Router /admin/users/{id} [get]
// GetAdminUser returns a user's information.
func GetAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	u, err := database.Handle.Users().GetByID(c.Req.Context(), userID)
	if err != nil {
		c.Error(err, "get user")
		return
	}
	c.JSON(http.StatusOK, ToUserResponse(database.FromDomainUser(u)))
}

// @Summary Get Admin Repositories
// @Description Returns a list of all repositories.
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {array} apiv2types.RepositoryResponse
// @Router /admin/repos [get]
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

type adminCreateOrgPayload struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name"`
	OwnerID  int64  `json:"owner_id" binding:"required"`
}

// @Summary Create Organization (Admin)
// @Description Creates a new organization administratively.
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body adminCreateOrgPayload true "Organization details"
// @Success 201 {object} apiv2types.OrganizationResponse
// @Router /admin/orgs [post]
// CreateAdminOrg creates a new organization administratively.
func CreateAdminOrg(c *context.APIContext) {
	var f adminCreateOrgPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	u, err := database.Handle.Users().GetByID(c.Req.Context(), f.OwnerID)
	if err != nil {
		c.StatusError(http.StatusNotFound, "Owner not found", err)
		return
	}
	owner := database.FromDomainUser(u)

	org := &database.User{
		Name:     f.Username,
		FullName: f.FullName,
		Type:     database.UserTypeOrganization,
	}

	if err := database.CreateOrganization(org, owner); err != nil {
		c.Error(err, "create organization")
		return
	}

	c.JSON(http.StatusCreated, ToOrganizationResponse(org, owner))
}

// @Summary Delete Repository (Admin)
// @Description Deletes a repository administratively.
// @Tags admin
// @Param id path int64 true "Repository ID"
// @Success 200 {object} map[string]string
// @Router /admin/repos/{id} [delete]
// DeleteAdminRepo deletes a repository administratively.
func DeleteAdminRepo(c *context.APIContext) {
	repoID := c.ParamsInt64("id")
	log.Trace("Admin: Deleting repository [repo_id: %d]", repoID)
	repo, err := database.GetRepositoryByID(repoID)
	if err != nil {
		c.Error(err, "get repository")
		return
	}

	log.Trace("Admin: Repository found: %s/%s. Starting database deletion...", repo.Owner.Name, repo.Name)
	removePikoCIPipeline(c.Req.Context(), repo.ID)
	if err := database.DeleteRepository(repo.OwnerID, repo.ID); err != nil {
		c.Error(err, "delete repository")
		return
	}

	log.Trace("Admin: Repository deleted successfully: [repo_id: %d]", repoID)
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

type adminUpdateUserPayload struct {
	FullName         *string `json:"full_name"`
	Website          *string `json:"website"`
	Location         *string `json:"location"`
	MaxRepoCreation  *int    `json:"max_repo_creation"`
	IsActivated      *bool   `json:"activated"`
	IsAdmin          *bool   `json:"admin"`
	ProhibitLogin    *bool   `json:"prohibit_login"`
	AllowGitHook     *bool   `json:"allow_git_hook"`
	AllowImportLocal *bool   `json:"allow_import_local"`
}

// @Summary Update User (Admin)
// @Description Updates a user's information administratively.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int64 true "User ID"
// @Param payload body adminUpdateUserPayload true "User details"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id} [patch]
// UpdateAdminUser updates a user's information administratively.
func UpdateAdminUser(c *context.APIContext) {
	userID := c.ParamsInt64("id")
	var f adminUpdateUserPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	opts := database.UpdateUserOptions{
		FullName:         f.FullName,
		Website:          f.Website,
		Location:         f.Location,
		MaxRepoCreation:  f.MaxRepoCreation,
		IsActivated:      f.IsActivated,
		IsAdmin:          f.IsAdmin,
		ProhibitLogin:    f.ProhibitLogin,
		AllowGitHook:     f.AllowGitHook,
		AllowImportLocal: f.AllowImportLocal,
	}

	if err := database.Handle.Users().Update(c.Req.Context(), userID, opts); err != nil {
		c.Error(err, "update user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// @Summary Delete User (Admin)
// @Description Deletes a user account administratively.
// @Tags admin
// @Param id path int64 true "User ID"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id} [delete]
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
	sources, err := database.Handle.LoginSources().List(c.Req.Context(), domain.ListLoginSourceOptions{})
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
		Type      int    `json:"type" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Activated bool   `json:"activated"`
		Default   bool   `json:"default"`
		Config    any    `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	source, err := database.Handle.LoginSources().Create(c.Req.Context(), domain.CreateLoginSourceOptions{
		Type:      f.Type,
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

type adminUpdateConfigPayload struct {
	Section string            `json:"section" binding:"required"`
	Data    map[string]string `json:"data" binding:"required"`
}

// @Summary Update Admin Configuration
// @Description Updates server configurations for a specific section.
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body adminUpdateConfigPayload true "Configuration details"
// @Success 200 {object} map[string]any
// @Router /admin/config [put]
// UpdateAdminConfig updates server configurations for a specific section.
func UpdateAdminConfig(c *context.APIContext) {
	var f adminUpdateConfigPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	section := conf.File.Section(f.Section)
	for k, v := range f.Data {
		section.Key(k).SetValue(v)
	}

	if err := conf.File.SaveTo(conf.CustomConf); err != nil {
		c.Error(err, "save configuration")
		return
	}

	// Reload config to memory (some parts might need specific mapping,
	// but for now we write to disk which is the source of truth for next restart)
	// For immediate effect on some settings, we'd need more complex mapping.
	// But according to the goal, we are prioritizing manageability.

	// Record action in Audit Log
	content, _ := json.Marshal(f.Data)
	if err := database.Handle.Actions().UpdateConfig(c.Req.Context(), c.User.ToDomain(), f.Section, string(content)); err != nil {
		log.Error("failed to update: %v", err)
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
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

	opts := domain.UpdateLoginSourceOptions{
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

// GetAdminUnmanagedRepos returns a list of repositories on disk that are not in the database.
func GetAdminUnmanagedRepos(c *context.APIContext) {
	unmanaged, err := repo.Repos.ListUnmanagedRepositories(c.Req.Context())
	if err != nil {
		c.Error(err, "list unmanaged repositories")
		return
	}
	c.JSONSuccess(unmanaged)
}

type adminAdoptRepoPayload struct {
	OwnerName string `json:"owner_name" binding:"required"`
	RepoName  string `json:"repo_name" binding:"required"`
}

// @Summary Adopt Repository (Admin)
// @Description Adopts an unmanaged repository folder into the database.
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body adminAdoptRepoPayload true "Adoption details"
// @Success 200 {object} apiv2types.RepositoryResponse
// @Router /admin/repos/adopt [post]
// AdoptAdminRepo adopts an unmanaged repository folder into the database.
func AdoptAdminRepo(c *context.APIContext) {
	var f adminAdoptRepoPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	r, err := repo.Repos.AdoptRepository(c.Req.Context(), f.OwnerName, f.RepoName)
	if err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	c.JSONSuccess(ToRepositoryResponse(r))
}

// CheckAdminRepoIntegrity forces an integrity check on a repository.
func CheckAdminRepoIntegrity(c *context.APIContext) {
	repoID := c.ParamsInt64("id")
	r, err := database.GetRepositoryByID(repoID)
	if err != nil {
		c.Error(err, "get repository")
		return
	}

	if err := health.Sentinel.RunIntegrityCheck(c.Req.Context(), r); err != nil {
		c.Error(err, "run integrity check")
		return
	}

	c.NoContent()
}

// ExportAdminBackup generates and downloads a metadata backup.
func ExportAdminBackup(c *context.APIContext) {
	path, err := backup.CreateMetadataBackup()
	if err != nil {
		c.Error(err, "create backup")
		return
	}

	c.File(path)

	// Cleanup old backups in the background
	go func() {
		_ = backup.CleanupOldBackups()
	}()
}

// GetAdminJanitorStats returns a summary of janitor reports for all repositories.
func GetAdminJanitorStats(c *context.APIContext) {
	repos, err := database.RepositoriesWithUsers(1, 1000)
	if err != nil {
		c.Error(err, "list repositories")
		return
	}

	type JanitorStat struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Owner       string `json:"owner"`
		Waste       int64  `json:"waste"`
		LastCheck   int64  `json:"last_check"`
		IssuesCount int    `json:"issues_count"`
	}

	stats := make([]JanitorStat, 0)
	for _, r := range repos {
		var report struct {
			TotalWaste int64 `json:"total_waste"`
			Issues     []any `json:"issues"`
		}
		if r.JanitorReport != "" {
			_ = json.Unmarshal([]byte(r.JanitorReport), &report)
		}

		stats = append(stats, JanitorStat{
			ID:          r.ID,
			Name:        r.Name,
			Owner:       r.MustOwner().Name,
			Waste:       report.TotalWaste,
			LastCheck:   r.JanitorLastCheckUnix,
			IssuesCount: len(report.Issues),
		})
	}
	c.JSONSuccess(stats)
}

// ScanAllJanitor triggers a janitor scan for all repositories.
func ScanAllJanitor(c *context.APIContext) {
	repos, err := database.RepositoriesWithUsers(1, 1000)
	if err != nil {
		c.Error(err, "list repositories")
		return
	}

	scanned := 0
	for _, r := range repos {
		_, err := janitor.Scan(c.Req.Context(), r.ID)
		if err != nil {
			log.Error("janitor scan failed for repo %d: %v", r.ID, err)
			continue
		}
		scanned++
	}

	c.JSONSuccess(map[string]any{
		"scanned": scanned,
		"total":   len(repos),
	})
}
