package v2

import (
	"net/http"
	"net/url"
	"strings"

	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

type MigrateRepoForm struct {
	CloneAddr    string `json:"clone_addr" binding:"required"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	UID          int64  `json:"uid" binding:"required"`
	RepoName     string `json:"repo_name" binding:"required,max=100"`
	Private      bool   `json:"private"`
	Unlisted     bool   `json:"unlisted"`
	Description   string `json:"description" binding:"max=512"`
	IsArchived    bool   `json:"is_archived"`
	Website       string `json:"website" binding:"max=100"`
	DefaultBranch string `json:"default_branch" binding:"max=100"`
}

// MigrateRepo migrates a repository from a remote address.
func MigrateRepo(c *context.APIContext) {
	var f MigrateRepoForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	ctxUser, err := database.Handle.Users().GetByID(c.Req.Context(), f.UID)
	if err != nil {
		if database.IsErrUserNotExist(err) {
			c.StatusError(http.StatusNotFound, "Owner not found.", err)
		} else {
			c.Error(err, "get context user")
		}
		return
	}

	// Check permission
	if !c.User.IsAdmin && c.User.ID != ctxUser.ID {
		if !ctxUser.IsOrganization() || !ctxUser.IsOwnedBy(c.User.ID) {
			c.Status(http.StatusForbidden)
			return
		}
	}

	remoteAddr := strings.TrimSpace(f.CloneAddr)
	if strings.HasPrefix(remoteAddr, "git clone ") {
		remoteAddr = strings.TrimPrefix(remoteAddr, "git clone ")
	}

	if strings.HasPrefix(remoteAddr, "http://") || strings.HasPrefix(remoteAddr, "https://") {
		u, err := url.Parse(remoteAddr)
		if err == nil {
			if f.AuthUsername != "" || f.AuthPassword != "" {
				u.User = url.UserPassword(f.AuthUsername, f.AuthPassword)
			}
			remoteAddr = u.String()
		}
	}

	repo, err := database.MigrateRepository(c.User, ctxUser, database.MigrateRepoOptions{
		Name:        f.RepoName,
		Description: f.Description,
		IsPrivate:     f.Private || conf.Repository.ForcePrivate,
		IsUnlisted:    f.Unlisted,
		IsArchived:    f.IsArchived,
		Website:       f.Website,
		DefaultBranch: f.DefaultBranch,
		RemoteAddr:  remoteAddr,
	})

	if err != nil {
		if repo != nil {
			if errDelete := database.DeleteRepository(ctxUser.ID, repo.ID); errDelete != nil {
				log.Error("DeleteRepository: %v", errDelete)
			}
		}
		c.Error(err, "migrate repository")
		return
	}

	log.Trace("Repository migrated [%d]: %s/%s", repo.ID, ctxUser.Name, f.RepoName)
	c.JSON(http.StatusCreated, repo.APIFormat(ctxUser))
}
