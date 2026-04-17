package v2

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"mime/multipart"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/pquerna/otp/totp"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/email"
	"gogs.io/gogs/internal/form"
	userroute "gogs.io/gogs/internal/route/user"
	"gogs.io/gogs/internal/userx"
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
	Website         string `json:"website" binding:"omitempty,url,max=100"`
	Location        string `json:"location" binding:"max=50"`
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
		"name":      c.User.Name,
		"full_name": c.User.FullName,
		"email":     c.User.Email,
		"website":   c.User.Website,
		"location":  c.User.Location,
		"disable_news_feed": c.User.DisableNewsFeed,
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
				case database.IsErrUserAlreadyExist(errors.Cause(err)):
					c.StatusError(http.StatusUnprocessableEntity, "Username has already been taken.", err)
				case database.IsErrNameNotAllowed(errors.Cause(err)):
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
	c.JSON(http.StatusOK, emails)
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

	t, err := database.Handle.AccessTokens().Create(c.Req.Context(), c.User.ID, f.Name)
	if err != nil {
		if database.IsErrAccessTokenAlreadyExist(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Access token name has already been used.", err)
		} else {
			c.Error(err, "add access token")
		}
		return
	}

	c.JSON(http.StatusOK, t)
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
		Source   string                `json:"source"`
		Avatar   *multipart.FileHeader `form:"avatar"`
		Gravatar string                `json:"gravatar" binding:"omitempty,email,max=254"`
	}
	// Try to bind as multipart form first
	if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
		f.Source = c.Req.PostFormValue("source")
		f.Gravatar = c.Req.PostFormValue("gravatar")
		file, header, err := c.Req.FormFile("avatar")
		if err == nil {
			f.Avatar = header
			_ = file.Close()
		}
	} else {
		// Try JSON for Gravatar source
		if err := c.ShouldBindJSON(&f); err != nil {
			c.ErrorStatus(http.StatusBadRequest, err)
			return
		}
	}

	// Map to legacy form for the existing helper
	legacyForm := form.Avatar{
		Source:   f.Source,
		Avatar:   f.Avatar,
		Gravatar: f.Gravatar,
	}
	if err := userroute.UpdateAvatarSetting(c.Context, legacyForm, c.User); err != nil {
		c.StatusError(http.StatusInternalServerError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"avatar": c.User.AvatarURL(),
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
	c.JSON(http.StatusOK, repos)
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
	c.JSON(http.StatusOK, orgs)
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
