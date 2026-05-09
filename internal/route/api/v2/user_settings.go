package v2

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/email"
	log "github.com/nef734/ggit/internal/log"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/nef734/ggit/internal/userx"
	"github.com/pquerna/otp/totp"
)

type TwoFactorStatus struct {
	Enabled bool `json:"enabled"`
}

type TwoFactorEnrollment struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type ConfirmTwoFactorPayload struct {
	Passcode string `json:"passcode" binding:"required"`
}

type EditProfile struct {
	Name            string `json:"name" binding:"max=35"`
	FullName        string `json:"full_name" binding:"max=100"`
	Email           string `json:"email" binding:"omitempty,email,max=254"`
	Website         string `json:"website" binding:"omitempty,url,max=100"`
	Location        string `json:"location" binding:"max=50"`
	AvatarEmail     string `json:"avatar_email" binding:"omitempty,email,max=254"`
	DisableNewsFeed bool   `json:"disable_news_feed"`
}

type EditPassword struct {
	OldPassword string `json:"old_password" binding:"required,max=255"`
	Password    string `json:"password" binding:"required,max=255"`
	Retype      string `json:"retype"`
}

type AddEmail struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

type NewAccessToken struct {
	Name string `json:"name" binding:"required"`
}

type AddSSHKey struct {
	Title   string `json:"title" binding:"required,max=50"`
	Content string `json:"content" binding:"required"`
}

// GetUserProfileSettings returns the current user's profile settings.
func GetUserProfileSettings(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"username":               c.User.Name,
		"name":                   c.User.Name,
		"full_name":              c.User.FullName,
		"email":                  c.User.Email,
		"website":                c.User.Website,
		"location":               c.User.Location,
		"avatar_email":           c.User.AvatarEmail,
		"disable_news_feed":      c.User.DisableNewsFeed,
		"enable_dicebear_avatar": conf.Picture.EnableDicebearAvatar,
		"avatar_url":             c.User.AvatarURL(),
		"use_custom_avatar":      c.User.UseCustomAvatar,
	})
}

// UpdateUserProfileSettings updates the current user's profile settings.
func UpdateUserProfileSettings(c *context.APIContext) {
	var f EditProfile
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// Non-local users are not allowed to change their username
	if c.User.IsLocal() && f.Name != "" {
		// Check if the username (including cases) had been changed
		if c.User.Name != f.Name {
			err := database.Handle.Users().ChangeUsername(c.Req.Context(), c.User.ID, f.Name)
			if err != nil {
				switch {
				case database.IsErrUserAlreadyExist(err):
					c.StatusError(http.StatusUnprocessableEntity, "Username has already been taken.", err)
				case database.IsErrNameNotAllowed(err):
					c.StatusError(http.StatusBadRequest, "Username is not allowed.", err)
				default:
					c.Error(err, "change username")
				}
				return
			}

			log.Trace("Username changed: %s -> %s", c.User.Name, f.Name)
		}
	}

	err := database.Handle.Users().Update(
		c.Req.Context(),
		c.User.ID,
		database.UpdateUserOptions{
			FullName:        &f.FullName,
			Website:         &f.Website,
			Location:        &f.Location,
			AvatarEmail:     &f.AvatarEmail,
			DisableNewsFeed: &f.DisableNewsFeed,
		},
	)
	if err != nil {
		c.Error(err, "update user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// UpdateUserPassword updates the current user's password.
func UpdateUserPassword(c *context.APIContext) {
	var f EditPassword
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if !userx.ValidatePassword(c.User.Password, c.User.Salt, f.OldPassword) {
		c.StatusError(http.StatusUnauthorized, "Old password is incorrect.", nil)
		return
	}

	if f.Password != f.Retype {
		c.StatusError(http.StatusBadRequest, "Passwords do not match.", nil)
		return
	}

	err := database.Handle.Users().Update(
		c.Req.Context(),
		c.User.ID,
		database.UpdateUserOptions{
			Password: &f.Password,
		},
	)
	if err != nil {
		c.Error(err, "update user password")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUserEmails returns all email addresses for the current user.
func GetUserEmails(c *context.APIContext) {
	emails, err := database.Handle.Users().ListEmails(c.Req.Context(), c.User.ID)
	if err != nil {
		c.Error(err, "get email addresses")
		return
	}

	res := make([]*apiv2types.UserEmailResponse, len(emails))
	for i, e := range emails {
		res[i] = ToUserEmailResponse(e)
	}
	c.JSON(http.StatusOK, res)
}

// AddUserEmail adds a new email address for the current user.
func AddUserEmail(c *context.APIContext) {
	var f AddEmail
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	err := database.Handle.Users().AddEmail(c.Req.Context(), c.User.ID, f.Email, !conf.Auth.RequireEmailConfirmation)
	if err != nil {
		if database.IsErrEmailAlreadyUsed(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Email address has already been used.", err)
		} else {
			c.Error(err, "add email address")
		}
		return
	}

	// Send confirmation email
	if conf.Auth.RequireEmailConfirmation {
		if err := email.SendActivateEmailMail(database.NewMailerUser(c.User), f.Email, c.Tr("mail.activate_email")); err != nil {
			log.Error("Failed to send activate email mail: %v", err)
		}

		if err := c.Cache.Put("MailResendLimit_"+c.User.LowerName, c.User.LowerName, 180); err != nil {
			log.Error("Set cache 'MailResendLimit' failed: %v", err)
		}
		c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
			"info":   c.Tr("settings.add_email_confirmation_sent", f.Email, conf.Auth.ActivateCodeLives/60),
		})
	} else {
		c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	}
}

// DeleteUserEmail deletes an email address for the current user.
func DeleteUserEmail(c *context.APIContext) {
	emailStr := c.Query("email")
	if c.User.Email == emailStr {
		c.StatusError(http.StatusBadRequest, "Cannot delete primary email address.", nil)
		return
	}

	err := database.Handle.Users().DeleteEmail(c.Req.Context(), c.User.ID, emailStr)
	if err != nil {
		c.Error(err, "delete email address")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// MarkEmailPrimary marks an email address as primary for the current user.
func MarkEmailPrimary(c *context.APIContext) {
	emailStr := c.Query("email")
	err := database.Handle.Users().MarkEmailPrimary(c.Req.Context(), c.User.ID, emailStr)
	if err != nil {
		c.Error(err, "mark email as primary")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUserSSHKeys returns all SSH keys for the current user.
func GetUserSSHKeys(c *context.APIContext) {
	keys, err := database.ListPublicKeys(c.User.ID)
	if err != nil {
		c.Error(err, "list public keys")
		return
	}
	c.JSON(http.StatusOK, keys)
}

// AddUserSSHKey adds a new SSH key for the current user.
func AddUserSSHKey(c *context.APIContext) {
	var f AddSSHKey
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	content, err := database.CheckPublicKeyString(f.Content)
	if err != nil {
		if !database.IsErrKeyUnableVerify(err) {
			c.StatusError(http.StatusBadRequest, "Invalid SSH key: "+err.Error(), err)
			return
		}
	}

	key, err := database.AddPublicKey(c.User.ID, f.Title, content)
	if err != nil {
		switch {
		case database.IsErrKeyAlreadyExist(err):
			c.StatusError(http.StatusUnprocessableEntity, "SSH key has already been used.", err)
		case database.IsErrKeyNameAlreadyUsed(err):
			c.StatusError(http.StatusUnprocessableEntity, "SSH key name has already been used.", err)
		default:
			c.Error(err, "add public key")
		}
		return
	}

	c.JSON(http.StatusOK, key)
}

// DeleteUserSSHKey deletes an SSH key for the current user.
func DeleteUserSSHKey(c *context.APIContext) {
	err := database.DeletePublicKey(c.User, c.ParamsInt64("id"))
	if err != nil {
		c.Error(err, "delete public key")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUserAccessTokens returns all access tokens for the current user.
func GetUserAccessTokens(c *context.APIContext) {
	tokens, err := database.Handle.AccessTokens().List(c.Req.Context(), c.User.ID)
	if err != nil {
		c.Error(err, "list access tokens")
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// AddUserAccessToken adds a new access token for the current user.
func AddUserAccessToken(c *context.APIContext) {
	var f NewAccessToken
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	t, token, err := database.Handle.AccessTokens().Create(c.Req.Context(), c.User.ID, f.Name)
	if err != nil {
		if database.IsErrAccessTokenAlreadyExist(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Access token name has already been used.", err)
		} else {
			c.Error(err, "add access token")
		}
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id":   t.ID,
		"name": t.Name,
		"sha1": token,
	})
}

// DeleteUserAccessToken deletes an access token for the current user.
func DeleteUserAccessToken(c *context.APIContext) {
	err := database.Handle.AccessTokens().DeleteByID(c.Req.Context(), c.User.ID, c.ParamsInt64("id"))
	if err != nil {
		c.Error(err, "delete access token")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// UpdateUserAvatar updates the current user's avatar.
func UpdateUserAvatar(c *context.APIContext) {
	var f struct {
		Avatar *multipart.FileHeader `form:"avatar"`
	}

	header, err := c.FormFile("avatar")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			log.Error("Failed to get avatar from form: %v", err)
		}
	} else {
		f.Avatar = header
		log.Trace("Avatar file received: %s (size: %d)", header.Filename, header.Size)
	}

	if f.Avatar != nil {
		log.Trace("Processing avatar for user %d: %s", c.User.ID, f.Avatar.Filename)
		content, err := f.Avatar.Open()
		if err != nil {
			c.StatusError(http.StatusInternalServerError, "Open avatar file failed", err)
			return
		}
		defer func() { _ = content.Close() }()

		data, err := io.ReadAll(content)
		if err != nil {
			c.StatusError(http.StatusInternalServerError, "Read avatar file failed", err)
			return
		}

		if err := database.Handle.Users().UseCustomAvatar(c.Req.Context(), c.User.ID, data); err != nil {
			c.StatusError(http.StatusInternalServerError, "Use custom avatar failed", err)
			return
		}
		log.Trace("Successfully updated avatar for user %d", c.User.ID)
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":     "ok",
		"avatar_url": c.User.AvatarURL(),
	})
}

// DeleteUserAvatar deletes the current user's custom avatar.
func DeleteUserAvatar(c *context.APIContext) {
	err := database.Handle.Users().DeleteCustomAvatar(c.Req.Context(), c.User.ID)
	if err != nil {
		c.Error(err, "delete custom avatar")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUserRepositories returns all repositories for the current user.
func GetUserRepositories(c *context.APIContext) {
	repos, err := database.GetUserAndCollaborativeRepositories(c.User.ID)
	if err != nil {
		c.Error(err, "get user and collaborative repositories")
		return
	}
	if err = database.RepositoryList(repos).LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	res := make([]*apiv2types.RepositoryResponse, len(repos))
	for i, repo := range repos {
		if err = repo.GetOwner(); err != nil {
			c.Error(err, "get owner")
			return
		}
		res[i] = repo.APIFormat(repo.Owner)
		res[i].ContextUser = c.User.Name
	}
	c.JSON(http.StatusOK, res)
}

// LeaveUserRepository leaves a collaborative repository.
func LeaveUserRepository(c *context.APIContext) {
	repo, err := database.GetRepositoryByID(c.ParamsInt64("id"))
	if err != nil {
		c.NotFoundOrError(err, "get repository by ID")
		return
	}

	if err = repo.DeleteCollaboration(c.User.ID); err != nil {
		c.Error(err, "delete collaboration")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetUserOrganizations returns all organizations for the current user.
func GetUserOrganizations(c *context.APIContext) {
	orgs, err := database.GetOrgsByUserID(c.User.ID, true)
	if err != nil {
		c.Error(err, "get organizations by user ID")
		return
	}

	res := make([]*apiv2types.UserResponse, len(orgs))
	for i, org := range orgs {
		res[i] = org.APIFormat()
	}
	c.JSON(http.StatusOK, res)
}

// LeaveUserOrganization leaves an organization.
func LeaveUserOrganization(c *context.APIContext) {
	if err := database.RemoveOrgUser(c.ParamsInt64("id"), c.User.ID); err != nil {
		if database.IsErrLastOrgOwner(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Cannot leave organization because you are the last owner.", err)
		} else {
			c.Error(err, "remove organization user")
		}
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// CreateUserOrganization creates a new organization for the current user.
func CreateUserOrganization(c *context.APIContext) {
	var f struct {
		Name     string `json:"username" binding:"required,max=35"`
		FullName string `json:"full_name" binding:"max=100"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	org := &database.User{
		Name:     f.Name,
		FullName: f.FullName,
		Type:     database.UserTypeOrganization,
	}

	if err := database.CreateOrganization(org, c.User); err != nil {
		switch {
		case database.IsErrUserAlreadyExist(err):
			c.StatusError(http.StatusUnprocessableEntity, "Organization name has already been taken.", err)
		case database.IsErrNameNotAllowed(err):
			c.StatusError(http.StatusBadRequest, "Organization name is not allowed.", err)
		default:
			c.Error(err, "create organization")
		}
		return
	}

	c.JSON(http.StatusCreated, org.APIFormat())
}

// DeleteUserAccount deletes the current user's account.
func DeleteUserAccount(c *context.APIContext) {
	var payload struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if _, err := database.Handle.Users().Authenticate(c.Req.Context(), c.User.Name, payload.Password, c.User.LoginSource); err != nil {
		c.StatusError(http.StatusUnauthorized, "Invalid password.", err)
		return
	}

	if err := database.Handle.Users().DeleteByID(c.Req.Context(), c.User.ID, false); err != nil {
		switch {
		case database.IsErrUserOwnRepos(err):
			c.StatusError(http.StatusUnprocessableEntity, "Still owns repositories.", err)
		case database.IsErrUserHasOrgs(err):
			c.StatusError(http.StatusUnprocessableEntity, "Still has organizations.", err)
		default:
			c.Error(err, "delete user")
		}
		return
	}

	log.Trace("Account deleted: %s", c.User.Name)
	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetTwoFactorStatus returns whether 2FA is enabled for the current user.
func GetTwoFactorStatus(c *context.APIContext) {
	enabled := database.Handle.TwoFactors().IsEnabled(c.Req.Context(), c.User.ID)
	c.JSON(http.StatusOK, TwoFactorStatus{Enabled: enabled})
}

// EnableTwoFactor initiates the 2FA enrollment process.
func EnableTwoFactor(c *context.APIContext) {
	if database.Handle.TwoFactors().IsEnabled(c.Req.Context(), c.User.ID) {
		c.StatusError(http.StatusConflict, "Two-factor authentication is already enabled.", nil)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      conf.App.BrandName,
		AccountName: c.User.Email,
	})
	if err != nil {
		c.Error(err, "generate TOTP")
		return
	}

	img, err := key.Image(240, 240)
	if err != nil {
		c.Error(err, "generate image")
		return
	}

	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		c.Error(err, "encode image")
		return
	}

	qrCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	// Store secret in cache for confirmation (similar to legacy session)
	if err := c.Cache.Put(userx.TwoFactorCacheKey(c.User.ID, "enrollment"), key.Secret(), 600); err != nil {
		c.Error(err, "store 2FA secret in cache")
		return
	}

	c.JSON(http.StatusOK, TwoFactorEnrollment{
		Secret: key.Secret(),
		QRCode: qrCode,
	})
}

// ConfirmTwoFactor completes the 2FA enrollment process.
func ConfirmTwoFactor(c *context.APIContext) {
	var f ConfirmTwoFactorPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	secret, ok := c.Cache.Get(userx.TwoFactorCacheKey(c.User.ID, "enrollment")).(string)
	if !ok || secret == "" {
		c.StatusError(http.StatusBadRequest, "Two-factor enrollment session expired or not found.", nil)
		return
	}

	if !totp.Validate(f.Passcode, secret) {
		c.StatusError(http.StatusBadRequest, "Invalid passcode.", nil)
		return
	}

	if err := database.Handle.TwoFactors().Create(c.Req.Context(), c.User.ID, conf.Security.SecretKey, secret); err != nil {
		c.Error(err, "create two factor")
		return
	}

	c.Cache.Delete(userx.TwoFactorCacheKey(c.User.ID, "enrollment"))

	recoveryCodes, err := database.GetRecoveryCodesByUserID(c.User.ID)
	if err != nil {
		c.Error(err, "get recovery codes")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":         "ok",
		"recovery_codes": recoveryCodes,
	})
}

// DisableTwoFactor disables 2FA for the current user.
func DisableTwoFactor(c *context.APIContext) {
	if !database.Handle.TwoFactors().IsEnabled(c.Req.Context(), c.User.ID) {
		c.StatusError(http.StatusConflict, "Two-factor authentication is not enabled.", nil)
		return
	}

	if err := database.DeleteTwoFactor(c.User.ID); err != nil {
		c.Error(err, "delete two factor")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetTwoFactorRecoveryCodes returns recovery codes for the current user.
func GetTwoFactorRecoveryCodes(c *context.APIContext) {
	if !database.Handle.TwoFactors().IsEnabled(c.Req.Context(), c.User.ID) {
		c.StatusError(http.StatusForbidden, "Two-factor authentication is not enabled.", nil)
		return
	}

	recoveryCodes, err := database.GetRecoveryCodesByUserID(c.User.ID)
	if err != nil {
		c.Error(err, "get recovery codes")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"recovery_codes": recoveryCodes,
	})
}

// RegenerateTwoFactorRecoveryCodes generates new recovery codes.
func RegenerateTwoFactorRecoveryCodes(c *context.APIContext) {
	if !database.Handle.TwoFactors().IsEnabled(c.Req.Context(), c.User.ID) {
		c.StatusError(http.StatusForbidden, "Two-factor authentication is not enabled.", nil)
		return
	}

	if err := database.RegenerateRecoveryCodes(c.User.ID); err != nil {
		c.Error(err, "regenerate recovery codes")
		return
	}

	recoveryCodes, err := database.GetRecoveryCodesByUserID(c.User.ID)
	if err != nil {
		c.Error(err, "get recovery codes")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":         "ok",
		"recovery_codes": recoveryCodes,
	})
}

// GetUserLoginSessions returns all active login sessions for the current user.
func GetUserLoginSessions(c *context.APIContext) {
	sessions, err := database.Handle.LoginSessions().ListByUserID(c.Req.Context(), c.User.ID)
	if err != nil {
		c.Error(err, "list login sessions")
		return
	}

	currentSID := ""
	if sid, ok := c.Session.Get("sid").(string); ok {
		currentSID = sid
	}

	res := make([]*apiv2types.LoginSessionResponse, len(sessions))
	for i, s := range sessions {
		maskedSID := ""
		if len(s.SID) > 8 {
			maskedSID = s.SID[:4] + "...." + s.SID[len(s.SID)-4:]
		}

		res[i] = &apiv2types.LoginSessionResponse{
			ID:          s.ID,
			SID:         maskedSID,
			UserAgent:   s.UserAgent,
			RemoteAddr:  s.RemoteAddr,
			Current:     s.SID == currentSID,
			CreatedUnix: s.CreatedUnix,
			UpdatedUnix: s.UpdatedUnix,
			CreatedAt:   time.Unix(s.CreatedUnix, 0),
			UpdatedAt:   time.Unix(s.UpdatedUnix, 0),
		}
	}
	c.JSON(http.StatusOK, res)
}

// DeleteOtherUserLoginSessions deletes all login sessions for the current user except the current one.
func DeleteOtherUserLoginSessions(c *context.APIContext) {
	sid, ok := c.Session.Get("sid").(string)
	if !ok {
		c.StatusError(http.StatusUnauthorized, "No active session found.", nil)
		return
	}

	err := database.Handle.LoginSessions().DeleteOthersByUserID(c.Req.Context(), c.User.ID, sid)
	if err != nil {
		c.Error(err, "delete other login sessions")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteUserLoginSession deletes a specific login session for the current user.
func DeleteUserLoginSession(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if id <= 0 {
		c.StatusError(http.StatusBadRequest, "Invalid session ID.", nil)
		return
	}

	err := database.Handle.LoginSessions().DeleteByIDAndUserID(c.Req.Context(), id, c.User.ID)
	if err != nil {
		c.Error(err, "delete login session")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// GetBlockedUsers returns a list of users blocked by the current user.
func GetBlockedUsers(c *context.APIContext) {
	users, err := database.Handle.UserBlocks().ListByUserID(c.Req.Context(), c.User.ID)
	if err != nil {
		c.Error(err, "list blocked users")
		return
	}

	res := make([]*apiv2types.UserResponse, len(users))
	for i, u := range users {
		res[i] = database.FromDomainUser(u).APIFormat()
	}
	c.JSON(http.StatusOK, res)
}

// BlockUser blocks a user.
func BlockUser(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}

	if err := database.Handle.UserBlocks().Block(c.Req.Context(), c.User.ID, user.ID); err != nil {
		c.Error(err, "block user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// UnblockUser unblocks a user.
func UnblockUser(c *context.APIContext) {
	username := c.Params("username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}

	if err := database.Handle.UserBlocks().Unblock(c.Req.Context(), c.User.ID, user.ID); err != nil {
		c.Error(err, "unblock user")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
