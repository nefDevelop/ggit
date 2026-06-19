package lfs

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/authx"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/lfsx"
	log "github.com/nef734/ggit/internal/log"
)

// RegisterRoutes registers LFS routes using given router, and inherits all
// groups and middleware.
func RegisterRoutes(r gin.IRouter) {
	verifyAccept := verifyHeader("Accept", contentType, http.StatusNotAcceptable)
	verifyContentTypeJSON := verifyHeader("Content-Type", contentType, http.StatusBadRequest)
	verifyContentTypeStream := verifyHeader("Content-Type", "application/octet-stream", http.StatusBadRequest)

	store := NewStore()
	r.POST("/objects/batch", authenticate(store), authorize(store, domain.AccessModeRead), verifyAccept, verifyContentTypeJSON, serveBatch(store))
	basicGroup := r.Group("/objects/basic", authenticate(store))
	{
		basic := &basicHandler{
			store:          store,
			defaultStorage: lfsx.Storage(conf.LFS.Storage),
			storagers: map[lfsx.Storage]lfsx.Storager{
				lfsx.StorageLocal: &lfsx.LocalStorage{Root: conf.LFS.ObjectsPath, TempDir: conf.LFS.ObjectsTempPath},
			},
		}
		basicGroup.GET("/:oid", verifyOID(), authorize(store, domain.AccessModeRead), basic.serveDownload)
		basicGroup.PUT("/:oid", verifyOID(), authorize(store, domain.AccessModeWrite), verifyContentTypeStream, basic.serveUpload)
		basicGroup.POST("/verify", authorize(store, domain.AccessModeWrite), verifyAccept, verifyContentTypeJSON, basic.serveVerify)
	}
}

// authenticate tries to authenticate user via HTTP Basic Auth. It first tries to authenticate
// as plain username and password, then use username as access token if previous step failed.
func authenticate(store Store) gin.HandlerFunc {
	askCredentials := func(w http.ResponseWriter) {
		w.Header().Set("Lfs-Authenticate", `Basic realm="Git LFS"`)
		responseJSON(w, http.StatusUnauthorized, responseError{
			Message: "Credentials needed",
		})
	}

	return func(c *gin.Context) {
		username, password := authx.DecodeBasic(c.Request.Header)
		if username == "" {
			askCredentials(c.Writer)
			c.Abort()
			return
		}

		user, err := store.AuthenticateUser(c.Request.Context(), username, password, -1)
		if err != nil && !auth.IsErrBadCredentials(err) {
			internalServerError(c.Writer)
			log.Error("Failed to authenticate user [name: %s]: %v", username, err)
			c.Abort()
			return
		}

		if err == nil && store.IsTwoFactorEnabled(c.Request.Context(), user.ID) {
			c.String(http.StatusBadRequest, "Users with 2FA enabled are not allowed to authenticate via username and password.")
			c.Abort()
			return
		}

		// If username and password combination failed, try again using either username
		// or password as the token.
		if auth.IsErrBadCredentials(err) {
			user, err = context.AuthenticateByToken(store, c.Request.Context(), username)
			if err != nil && !database.IsErrAccessTokenNotExist(err) {
				internalServerError(c.Writer)
				log.Error("%v", err)
				c.Abort()
				return
			} else if database.IsErrAccessTokenNotExist(err) {
				// Try again using the password field as the token.
				user, err = context.AuthenticateByToken(store, c.Request.Context(), password)
				if err != nil {
					if database.IsErrAccessTokenNotExist(err) {
						askCredentials(c.Writer)
					} else {
						c.Status(http.StatusInternalServerError)
						log.Error("%v", err)
					}
					c.Abort()
					return
				}
			}
		}

		log.Trace("[LFS] Authenticated user: %s", user.Name)

		c.Set("actor", user)
	}
}

// authorize tries to authorize the user to the context repository with given access mode.
func authorize(store Store, mode domain.AccessMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := c.Get("actor")
		if !ok {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		u := actor.(*database.User)
		username := c.Param("username")
		reponame := strings.TrimSuffix(c.Param("reponame"), ".git")

		owner, err := store.GetUserByUsername(c.Request.Context(), username)
		if err != nil {
			if database.IsErrUserNotExist(err) {
				c.Status(http.StatusNotFound)
			} else {
				internalServerError(c.Writer)
				log.Error("Failed to get user [name: %s]: %v", username, err)
			}
			c.Abort()
			return
		}

		repo, err := store.GetRepositoryByName(c.Request.Context(), owner.ID, reponame)
		if err != nil {
			if database.IsErrRepoNotExist(err) {
				c.Status(http.StatusNotFound)
			} else {
				internalServerError(c.Writer)
				log.Error("Failed to get repository [owner_id: %d, name: %s]: %v", owner.ID, reponame, err)
			}
			c.Abort()
			return
		}

		if !store.AuthorizeRepositoryAccess(c.Request.Context(), u.ID, repo.ID, mode,
			domain.AccessModeOptions{
				OwnerID: repo.OwnerID,
				Private: repo.IsPrivate,
			},
		) {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		log.Trace("[LFS] Authorized user %q to %q", u.Name, username+"/"+reponame)

		c.Set("owner", owner)
		c.Set("repo", repo)
	}
}

// verifyHeader checks if the HTTP header contains given value.
// When not, response given "failCode" as status code.
func verifyHeader(key, value string, failCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		vals := c.Request.Header.Values(key)
		for _, val := range vals {
			if strings.Contains(val, value) {
				return
			}
		}

		log.Trace("[LFS] HTTP header %q does not contain value %q", key, value)
		c.Status(failCode)
		c.Abort()
	}
}

// verifyOID checks if the ":oid" URL parameter is valid.
func verifyOID() gin.HandlerFunc {
	return func(c *gin.Context) {
		oid := lfsx.OID(c.Param("oid"))
		if !lfsx.ValidOID(oid) {
			responseJSON(c.Writer, http.StatusBadRequest, responseError{
				Message: "Invalid oid",
			})
			c.Abort()
			return
		}

		c.Set("oid", oid)
	}
}

func internalServerError(w http.ResponseWriter) {
	responseJSON(w, http.StatusInternalServerError, responseError{
		Message: "Internal server error",
	})
}
