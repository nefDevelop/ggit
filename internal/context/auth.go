package context

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/tool"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		// Cannot view any page before installation.
		if !conf.Security.InstallLock {
			gctx.RedirectSubpath("/install")
			return
		}

		// Check prohibit login users.
		if gctx.IsLogged && gctx.User.ProhibitLogin {
			gctx.Data["Title"] = gctx.Tr("auth.prohibit_login")
			gctx.Success("user/auth/prohibit_login")
			return
		}

		// Check non-logged users landing page.
		if !gctx.IsLogged && gctx.Request.RequestURI == "/" && conf.Server.LandingURL != "/" {
			gctx.RedirectSubpath(conf.Server.LandingURL)
			return
		}

		// Redirect to dashboard if user tries to visit any non-login page.
		if options.SignOutRequired && gctx.IsLogged && gctx.Request.RequestURI != "/" {
			gctx.RedirectSubpath("/")
			return
		}

		if !options.SignOutRequired && !options.DisableCSRF && gctx.Request.Method == http.MethodPost && !isAPIPath(gctx.Request.URL.Path) {
			// TODO: CSRF validation for Gin
		}

		if options.SignInRequired {
			if !gctx.IsLogged {
				// Restrict API calls with error message.
				if isAPIPath(gctx.Request.URL.Path) {
					gctx.JSON(http.StatusForbidden, map[string]string{
						"message": "Only authenticated user is allowed to call APIs.",
					})
					return
				}

				gctx.SetCookie("redirect_to", url.QueryEscape(conf.Server.Subpath+gctx.Request.RequestURI), 0, conf.Server.Subpath, "", false, false)
				gctx.RedirectSubpath("/user/login")
				return
			} else if !gctx.User.IsActive && conf.Auth.RequireEmailConfirmation {
				gctx.Title("auth.active_your_account")
				gctx.Success("user/auth/activate")
				return
			}
		}

		// Redirect to log in page if auto-signin info is provided and has not signed in.
		if !options.SignOutRequired && !gctx.IsLogged && !isAPIPath(gctx.Request.URL.Path) {
			u, _ := gctx.Cookie(conf.Security.CookieUsername)
			if len(u) > 0 {
				gctx.SetCookie("redirect_to", url.QueryEscape(conf.Server.Subpath+gctx.Request.RequestURI), 0, conf.Server.Subpath, "", false, false)
				gctx.RedirectSubpath("/user/login")
				return
			}
		}

		if options.AdminRequired {
			if !gctx.User.IsAdmin {
				gctx.Status(http.StatusForbidden)
				return
			}
			gctx.PageIs("Admin")
		}
		c.Next()
	}
}

func isAPIPath(url string) bool {
	return strings.HasPrefix(url, "/api/")
}

type AuthStore interface {
	// GetAccessTokenBySHA1 returns the access token with given SHA1. It returns
	// database.ErrAccessTokenNotExist when not found.
	GetAccessTokenBySHA1(ctx context.Context, sha1 string) (*database.AccessToken, error)
	// TouchAccessTokenByID updates the updated time of the given access token to
	// the current time.
	TouchAccessTokenByID(ctx context.Context, id int64) error

	// GetUserByID returns the user with given ID. It returns
	// database.ErrUserNotExist when not found.
	GetUserByID(ctx context.Context, id int64) (*database.User, error)
	// GetUserByUsername returns the user with given username. It returns
	// database.ErrUserNotExist when not found.
	GetUserByUsername(ctx context.Context, username string) (*database.User, error)
	// CreateUser creates a new user and persists to database. It returns
	// database.ErrNameNotAllowed if the given name or pattern of the name is not
	// allowed as a username, or database.ErrUserAlreadyExist when a user with same
	// name already exists, or database.ErrEmailAlreadyUsed if the email has been
	// verified by another user.
	CreateUser(ctx context.Context, username, email string, opts database.CreateUserOptions) (*database.User, error)
	// AuthenticateUser validates username and password via given login source ID.
	// It returns database.ErrUserNotExist when the user was not found.
	//
	// When the "loginSourceID" is negative, it aborts the process and returns
	// database.ErrUserNotExist if the user was not found in the database.
	//
	// When the "loginSourceID" is non-negative, it returns
	// database.ErrLoginSourceMismatch if the user has different login source ID
	// than the "loginSourceID".
	//
	// When the "loginSourceID" is positive, it tries to authenticate via given
	// login source and creates a new user when not yet exists in the database.
	AuthenticateUser(ctx context.Context, login, password string, loginSourceID int64) (*database.User, error)
}

// authenticatedUserID returns the ID of the authenticated user, along with a bool value
// which indicates whether the user uses token authentication.
func authenticatedUserID(store AuthStore, c *gin.Context, sess any) (_ int64, isTokenAuth bool) {
	if !database.HasEngine {
		return 0, false
	}

	// Check access token.
	if isAPIPath(c.Request.URL.Path) {
		var tokenSHA string
		auHead := c.Request.Header.Get("Authorization")
		if auHead != "" {
			auths := strings.Fields(auHead)
			if len(auths) == 2 && auths[0] == "token" {
				tokenSHA = auths[1]
			}
		}

		// Let's see if token is valid.
		if len(tokenSHA) > 0 {
			t, err := store.GetAccessTokenBySHA1(c.Request.Context(), tokenSHA)
			if err != nil {
				if !database.IsErrAccessTokenNotExist(err) {
					log.Error("GetAccessTokenBySHA: %v", err)
				}
				return 0, false
			}
			if err = store.TouchAccessTokenByID(c.Request.Context(), t.ID); err != nil {
				log.Error("Failed to touch access token: %v", err)
			}
			return t.UserID, true
		}
	}

	s := sess.(sessions.Session)
	uid := s.Get("uid")
	sid := s.Get("sid")
	if uid == nil || sid == nil {
		return 0, false
	}

	id, ok := uid.(int64)
	if !ok {
		return 0, false
	}

	sessionID, ok := sid.(string)
	if !ok {
		return 0, false
	}

	// Verify session in database
	_, err := database.Handle.LoginSessions().GetBySID(c.Request.Context(), sessionID)
	if err != nil {
		// Session not found or revoked
		s.Clear()
		_ = s.Save()
		return 0, false
	}

	_, err = store.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if !database.IsErrUserNotExist(err) {
			log.Error("Failed to get user by ID: %v", err)
		}
		return 0, false
	}
	return id, false
}

// authenticatedUser returns the user object of the authenticated user, along with two bool values
// which indicate whether the user uses HTTP Basic Authentication or token authentication respectively.
func authenticatedUser(store AuthStore, ctx *gin.Context, sess any) (_ *database.User, isBasicAuth, isTokenAuth bool) {
	if !database.HasEngine {
		return nil, false, false
	}

	uid, isTokenAuth := authenticatedUserID(store, ctx, sess)

	if uid <= 0 {
		if conf.Auth.EnableReverseProxyAuthentication {
			webAuthUser := ctx.Request.Header.Get(conf.Auth.ReverseProxyAuthenticationHeader)
			if len(webAuthUser) > 0 {
				user, err := store.GetUserByUsername(ctx.Request.Context(), webAuthUser)
				if err != nil {
					if !database.IsErrUserNotExist(err) {
						log.Error("Failed to get user by name: %v", err)
						return nil, false, false
					}

					// Check if enabled auto-registration.
					if conf.Auth.EnableReverseProxyAutoRegistration {
						user, err = store.CreateUser(
							ctx.Request.Context(),
							webAuthUser,
							uuid.New().String()+"@localhost",
							database.CreateUserOptions{
								Activated: true,
							},
						)
						if err != nil {
							log.Error("Failed to create user %q: %v", webAuthUser, err)
							return nil, false, false
						}
					}
				}
				return user, false, false
			}
		}

		// Check with basic auth.
		baHead := ctx.Request.Header.Get("Authorization")
		if len(baHead) > 0 {
			auths := strings.Fields(baHead)
			if len(auths) == 2 && auths[0] == "Basic" {
				uname, passwd, _ := tool.BasicAuthDecode(auths[1])

				u, err := store.AuthenticateUser(ctx.Request.Context(), uname, passwd, -1)
				if err != nil {
					if !auth.IsErrBadCredentials(err) {
						log.Error("Failed to authenticate user: %v", err)
					}
					return nil, false, false
				}

				return u, true, false
			}
		}
		return nil, false, false
	}

	u, err := store.GetUserByID(ctx.Request.Context(), uid)
	if err != nil {
		log.Error("GetUserByID: %v", err)
		return nil, false, false
	}
	return u, false, isTokenAuth
}

// AuthenticateByToken attempts to authenticate a user by the given access
// token. It returns database.ErrAccessTokenNotExist when the access token does not
// exist.
func AuthenticateByToken(store AuthStore, ctx context.Context, token string) (*database.User, error) {
	t, err := store.GetAccessTokenBySHA1(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "get access token by SHA1")
	}
	if err = store.TouchAccessTokenByID(ctx, t.ID); err != nil {
		// NOTE: There is no need to fail the auth flow if we can't touch the token.
		log.Error("Failed to touch access token [id: %d]: %v", t.ID, err)
	}

	user, err := store.GetUserByID(ctx, t.UserID)
	if err != nil {
		return nil, errors.Wrapf(err, "get user by ID [user_id: %d]", t.UserID)
	}
	return user, nil
}
