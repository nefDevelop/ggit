package v2

import (
	gocontext "context"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/email"
	authservice "github.com/nef734/ggit/internal/services/auth"
	"github.com/nef734/ggit/internal/strx"
	"github.com/nef734/ggit/internal/tool"
	"github.com/nef734/ggit/internal/userx"
	log "unknwon.dev/clog/v2"
)

type loginPayload struct {
	UserName    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Remember    bool   `json:"remember"`
	LoginSource int64  `json:"login_source"`
}

// @Summary User Login
// @Description Authenticates a user and establishes a session via cookies.
// @Tags user
// @Accept json
// @Produce json
// @Param payload body loginPayload true "Login credentials"
// @Success 200 {object} map[string]any "ok and user object"
// @Success 202 {object} map[string]any "two_factor_required"
// @Failure 400 {object} error "Bad request"
// @Failure 401 {object} string "Invalid username or password"
// @Router /user/login [post]
// Login authenticates a user and establishes a session via cookies.
func Login(c *context.APIContext) {
	if !database.HasEngine || authservice.Auth == nil {
		c.StatusError(http.StatusServiceUnavailable, "Ggit has not been installed yet. Please visit /install.", nil)
		return
	}

	var payload loginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	u, err := authservice.Auth.Authenticate(c.Req.Context(), authservice.LoginOptions{
		Username:    payload.UserName,
		Password:    payload.Password,
		LoginSource: payload.LoginSource,
	})
	if err != nil {
		if auth.IsErrBadCredentials(err) || database.IsErrLoginSourceMismatch(err) {
			c.StatusError(http.StatusUnauthorized, "Invalid username or password.", err)
		} else {
			c.Error(err, "authenticate user")
		}
		return
	}

	if database.Handle.TwoFactors().IsEnabled(c.Req.Context(), u.ID) {
		c.Session.Set("twoFactorRemember", payload.Remember)
		c.Session.Set("twoFactorUserID", u.ID)
		
		c.JSON(http.StatusAccepted, map[string]any{
			"status": "two_factor_required",
		})
		return
	}

	// We proceed to set session cookies.

	if payload.Remember {
		days := 86400 * conf.Security.LoginRememberDays
		c.SetCookie(conf.Security.CookieUsername, u.Name, days, conf.Server.Subpath, "", conf.Security.CookieSecure, true)
		c.SetSuperSecureCookie(u.Rands+u.Password, conf.Security.CookieRememberName, u.Name, days, conf.Server.Subpath, "", conf.Security.CookieSecure, true)
	}

	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.Name)

	sid, _ := strx.RandomChars(32)
	c.Session.Set("sid", sid)
	_, err = database.Handle.LoginSessions().Create(c.Req.Context(), u.ID, sid, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		log.Error("Failed to create login session: %v", err)
	}

	if err := c.Session.Save(); err != nil {
		c.Error(err, "save session")
		return
	}

	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.Subpath)
	if conf.Security.EnableLoginStatusCookie {
		c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Server.Subpath)
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"user":   ToUserResponse(u),
	})
}

// CheckAuth verifies the current session and returns the user object.
func CheckAuth(c *context.APIContext) {
	if !c.IsLogged {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"user":   ToUserResponse(c.User),
	})
}

// Logout terminates the current session.
func Logout(c *context.APIContext) {
	sid := c.Session.Get("sid")
	if sessionID, ok := sid.(string); ok {
		_ = database.Handle.LoginSessions().DeleteBySID(c.Req.Context(), sessionID)
	}

	c.Session.Clear()
	_ = c.Session.Save()
	c.SetCookie(conf.Security.CookieUsername, "", -1, conf.Server.Subpath)
	c.SetCookie(conf.Security.CookieRememberName, "", -1, conf.Server.Subpath)
	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.Subpath)
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

type registerPayload struct {
	UserName      string `json:"username" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	MiddleName    string `json:"middle_name"`    // Honeypot
	CaptchaAnswer string `json:"captcha_answer"` // Question answer
}

// @Summary User Registration
// @Description Creates a new user account.
// @Tags user
// @Accept json
// @Produce json
// @Param payload body registerPayload true "Registration details"
// @Success 201 {object} map[string]any "ok"
// @Failure 400 {object} error "Bad request"
// @Failure 403 {object} string "Registration disabled"
// @Failure 409 {object} string "Conflict"
// @Router /user/register [post]
// Register creates a new user account.
func Register(c *context.APIContext) {
	if conf.Auth.DisableRegistration {
		c.StatusError(http.StatusForbidden, "Registration is disabled.", nil)
		return
	}

	var f registerPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// 1. Honeypot check (must be empty)
	if len(f.MiddleName) > 0 {
		log.Warn("Registration blocked by honeypot: %s", f.UserName)
		c.Status(http.StatusNoContent) // Stealth drop
		return
	}

	// 2. Question Captcha check (if enabled)
	if conf.Auth.EnableRegistrationQuestion {
		if strings.TrimSpace(f.CaptchaAnswer) != conf.Auth.RegistrationAnswer {
			c.StatusError(http.StatusBadRequest, "Incorrect answer to security question.", nil)
			return
		}
	}

	if !database.HasEngine || authservice.Auth == nil {
		c.StatusError(http.StatusServiceUnavailable, "Ggit has not been installed yet. Please visit /install.", nil)
		return
	}

	u, err := authservice.Auth.Register(c.Req.Context(), authservice.SignupOptions{
		Username: f.UserName,
		Email:    f.Email,
		Password: f.Password,
	})
	if err != nil {
		switch {
		case database.IsErrUserAlreadyExist(err):
			c.StatusError(http.StatusConflict, "Username has been already taken.", err)
		case database.IsErrEmailAlreadyUsed(err):
			c.StatusError(http.StatusConflict, "Email address has been already used.", err)
		case database.IsErrNameNotAllowed(err):
			c.StatusError(http.StatusBadRequest, "Username is not allowed.", err)
		default:
			c.Error(err, "register user")
		}
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
		"requires_email_activation": !u.IsActive && conf.Auth.RequireEmailConfirmation && u.ID > 1,
	})
}

// -----------------------------------------------------------------------------
// Two-Factor & Recovery
// -----------------------------------------------------------------------------

type loginTwoFactorPayload struct {
	Passcode string `json:"passcode" binding:"required"`
}

// LoginTwoFactor handles the second step of 2FA authentication.
func LoginTwoFactor(c *context.APIContext) {
	userID, ok := c.Session.Get("twoFactorUserID").(int64)
	if !ok {
		c.StatusError(http.StatusUnauthorized, "Two-factor session expired or invalid.", nil)
		return
	}

	var payload loginTwoFactorPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	t, err := database.Handle.TwoFactors().GetByUserID(c.Req.Context(), userID)
	if err != nil {
		c.Error(err, "get two factor by user ID")
		return
	}

	valid, err := t.ValidateTOTP(payload.Passcode)
	if err != nil {
		c.Error(err, "validate TOTP")
		return
	} else if !valid {
		c.StatusError(http.StatusUnauthorized, "Invalid passcode.", nil)
		return
	}

	du, err := database.Handle.Users().GetByID(c.Req.Context(), userID)
	if err != nil {
		c.Error(err, "get user by ID")
		return
	}
	u := database.FromDomainUser(du)

	// Prevent same passcode from being reused
	if c.Cache.IsExist(userx.TwoFactorCacheKey(u.ID, payload.Passcode)) {
		c.StatusError(http.StatusTooManyRequests, "Passcode has already been used.", nil)
		return
	}
	if err := c.Cache.Put(userx.TwoFactorCacheKey(u.ID, payload.Passcode), 1, 60); err != nil {
		log.Error("Failed to put cache 'two factor passcode': %v", err)
	}

	remember, _ := c.Session.Get("twoFactorRemember").(bool)

	// Proceed to setup session cookies
	if remember {
		days := 86400 * conf.Security.LoginRememberDays
		c.SetCookie(conf.Security.CookieUsername, u.Name, days, conf.Server.Subpath, "", conf.Security.CookieSecure, true)
		c.SetSuperSecureCookie(u.Rands+u.Password, conf.Security.CookieRememberName, u.Name, days, conf.Server.Subpath, "", conf.Security.CookieSecure, true)
	}

	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.Name)

	sid, _ := strx.RandomChars(32)
	c.Session.Set("sid", sid)
	_, err = database.Handle.LoginSessions().Create(c.Req.Context(), u.ID, sid, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		log.Error("Failed to create login session: %v", err)
	}

	c.Session.Delete("twoFactorRemember")
	c.Session.Delete("twoFactorUserID")

	c.SetCookie(conf.Session.CSRFCookieName, "", -1, conf.Server.Subpath)
	if conf.Security.EnableLoginStatusCookie {
		c.SetCookie(conf.Security.LoginStatusCookieName, "true", 0, conf.Server.Subpath)
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"user": map[string]any{
			"id":       u.ID,
			"username": u.Name,
			"email":    u.Email,
			"avatar":   u.AvatarURL(),
			"is_admin": u.IsAdmin,
			"disable_news_feed": u.DisableNewsFeed,
		},
	})
}

// -----------------------------------------------------------------------------
// Forgot & Reset Password
// -----------------------------------------------------------------------------

type forgotPasswordPayload struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword initiates the password recovery flow.
func ForgotPassword(c *context.APIContext) {
	if !conf.Email.Enabled {
		c.StatusError(http.StatusForbidden, "Email service is disabled.", nil)
		return
	}

	var payload forgotPasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	u, err := database.Handle.Users().GetByEmail(c.Req.Context(), payload.Email)
	if err != nil {
		if database.IsErrUserNotExist(err) {
			// For security, do not reveal if the email exists. Always return OK.
			c.JSON(http.StatusOK, map[string]any{"status": "ok", "message": "If the email exists, a reset link has been sent."})
			return
		}
		c.Error(err, "get user by email")
		return
	}

	if !database.FromDomainUser(u).IsLocal() {
		c.StatusError(http.StatusForbidden, "Non-local accounts cannot reset passwords here.", nil)
		return
	}

	if c.Cache.IsExist(userx.MailResendCacheKey(u.ID)) {
		c.StatusError(http.StatusTooManyRequests, "Please wait before requesting another email.", nil)
		return
	}

	if err = email.SendResetPasswordMail(database.NewMailerUser(database.FromDomainUser(u)), c.Tr("mail.reset_password")); err != nil {
		log.Error("Failed to send reset password mail: %v", err)
	}
	if err = c.Cache.Put(userx.MailResendCacheKey(u.ID), 1, 180); err != nil {
		log.Error("Failed to put cache key 'mail resend': %v", err)
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok", "message": "Password reset email sent."})
}

type resetPasswordPayload struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func parseUserFromCode(code string) *database.User {
	if len(code) <= tool.TimeLimitCodeLength {
		return nil
	}
	hexStr := code[tool.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		if user, err := database.Handle.Users().GetByUsername(gocontext.TODO(), string(b)); user != nil {
			return database.FromDomainUser(user)
		} else if !database.IsErrUserNotExist(err) {
			log.Error("Failed to get user by name %q: %v", string(b), err)
		}
	}
	return nil
}

func verifyUserActiveCode(code string) *database.User {
	minutes := conf.Auth.ActivateCodeLives
	if user := parseUserFromCode(code); user != nil {
		prefix := code[:tool.TimeLimitCodeLength]
		data := strconv.FormatInt(user.ID, 10) + user.Email + user.LowerName + user.Password + user.Rands

		if tool.VerifyTimeLimitCode(data, minutes, prefix) {
			return user
		}
	}
	return nil
}

// ResetPassword finalizes the password recovery flow.
func ResetPassword(c *context.APIContext) {
	var payload resetPasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	u := verifyUserActiveCode(payload.Code)
	if u == nil {
		c.StatusError(http.StatusUnauthorized, "Invalid or expired reset code.", nil)
		return
	}

	if len(payload.Password) < 6 {
		c.StatusError(http.StatusBadRequest, "Password is too short.", nil)
		return
	}

	err := database.Handle.Users().Update(c.Req.Context(), u.ID, database.UpdateUserOptions{Password: &payload.Password})
	if err != nil {
		c.Error(err, "update user")
		return
	}

	log.Trace("User password reset: %s", u.Name)
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// -----------------------------------------------------------------------------
// Account Validation
// -----------------------------------------------------------------------------

type activatePayload struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email"` // Optional, used for secondary email activation
}

func verifyActiveEmailCode(code, email string) *database.EmailAddress {
	minutes := conf.Auth.ActivateCodeLives
	if user := parseUserFromCode(code); user != nil {
		prefix := code[:tool.TimeLimitCodeLength]
		data := strconv.FormatInt(user.ID, 10) + email + user.LowerName + user.Password + user.Rands

		if tool.VerifyTimeLimitCode(data, minutes, prefix) {
			emailAddress, err := database.Handle.Users().GetEmail(gocontext.TODO(), user.ID, email, false)
			if err == nil {
				return emailAddress
			}
		}
	}
	return nil
}

// Activate verifies a user's account or secondary email address.
func Activate(c *context.APIContext) {
	var payload activatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// Case 1: Account Activation
	if payload.Email == "" {
		u := verifyUserActiveCode(payload.Code)
		if u == nil {
			c.StatusError(http.StatusUnauthorized, "Invalid or expired activation code.", nil)
			return
		}

		v := true
		err := database.Handle.Users().Update(c.Req.Context(), u.ID, database.UpdateUserOptions{
			GenerateNewRands: true,
			IsActivated:      &v,
		})
		if err != nil {
			c.Error(err, "update user")
			return
		}

		log.Trace("User activated: %s", u.Name)
		
		// Setup session immediately
		c.Session.Set("uid", u.ID)
		c.Session.Set("uname", u.Name)
		
		c.JSON(http.StatusOK, map[string]any{"status": "ok"})
		return
	}

	// Case 2: Email Address Activation
	emailAddr := verifyActiveEmailCode(payload.Code, payload.Email)
	if emailAddr == nil {
		c.StatusError(http.StatusUnauthorized, "Invalid or expired email activation code.", nil)
		return
	}

	err := database.Handle.Users().MarkEmailActivated(c.Req.Context(), emailAddr.UserID, emailAddr.Email)
	if err != nil {
		c.Error(err, "activate email")
		return
	}

	log.Trace("Email activated: %s", emailAddr.Email)
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetRegistrationSettings returns the current registration security settings.
func GetRegistrationSettings(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"enable_question": conf.Auth.EnableRegistrationQuestion,
		"question":        conf.Auth.RegistrationQuestion,
	})
}
