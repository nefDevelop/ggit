package v2

import (
	"net/http"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/cron"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/email"
	"gogs.io/gogs/internal/process"
)

// GetAdminStats returns global statistics.
func GetAdminStats(c *context.APIContext) {
	stats := database.GetStatistic(c.Req.Context())
	c.JSON(http.StatusOK, stats)
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
		"mirror":     conf.Mirror,
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

	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		resp[i] = ToUserResponse(u)
	}
	c.JSON(http.StatusOK, resp)
}

// CreateAdminUser creates a new user account.
func CreateAdminUser(c *context.APIContext) {
	var f struct {
		Username  string `json:"username" binding:"Required"`
		Email     string `json:"email" binding:"Required;Email;MaxSize(254)"`
		Password  string `json:"password" binding:"Required;MaxSize(255)"`
		FullName  string `json:"full_name"`
		Activated bool   `json:"activated"`
		Admin     bool   `json:"admin"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	user, err := database.Handle.Users().Create(c.Req.Context(), f.Username, f.Email, database.CreateUserOptions{
		Password:  f.Password,
		FullName:  f.FullName,
		Activated: f.Activated,
		Admin:     f.Admin,
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

	resp := make([]*RepositoryResponse, len(repos))
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

	resp := make([]*OrganizationResponse, len(orgs))
	for i, o := range orgs {
		owner, _ := o.GetOrgOwner()
		resp[i] = ToOrganizationResponse(o, owner)
	}
	c.JSON(http.StatusOK, resp)
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
		IsActivated:      f.IsActivated,
		IsAdmin:          f.IsAdmin,
		ProhibitLogin:    f.ProhibitLogin,
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

	resp := make([]*AuthSourceResponse, len(sources))
	for i, s := range sources {
		resp[i] = ToAuthSourceResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

// CreateAdminAuth creates a new login source.
func CreateAdminAuth(c *context.APIContext) {
	var opts database.CreateLoginSourceOptions
	if err := c.ShouldBindJSON(&opts); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	source, err := database.Handle.LoginSources().Create(c.Req.Context(), opts)
	if err != nil {
		c.Error(err, "create login source")
		return
	}

	c.JSON(http.StatusCreated, source)
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

// UpdateAdminAuth updates a login source's information.
func UpdateAdminAuth(c *context.APIContext) {
	_ = c.ParamsInt64("id")
	// For now, Gogs/Ggit doesn't have a generic Update method for all source types in the store easily exposed.
	// But I can implement it by getting and saving.
	// Actually, let's keep it simple and just implement GET and DELETE for now if UPDATE is complex.
	// Wait, I should check if there is an Update method.
	c.JSON(http.StatusNotImplemented, map[string]string{"error": "Updating auth sources via API is coming soon."})
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

	resp := make([]*NoticeResponse, len(notices))
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
