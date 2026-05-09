package v2

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	log "github.com/nef734/ggit/internal/log"
	reposervice "github.com/nef734/ggit/internal/services/repo"
)

type MigrateRepoForm struct {
	CloneAddr     string `json:"clone_addr" binding:"required"`
	AuthUsername  string `json:"auth_username"`
	AuthPassword  string `json:"auth_password"`
	UID           int64  `json:"uid" binding:"required"`
	RepoName      string `json:"repo_name" binding:"required,max=100"`
	Private       bool   `json:"private"`
	Unlisted      bool   `json:"unlisted"`
	Description   string `json:"description" binding:"max=512"`
	IsArchived    bool   `json:"is_archived"`
	Website       string `json:"website" binding:"max=100"`
	DefaultBranch string `json:"default_branch" binding:"max=100"`
}

// MigrateRepo migrates a repository from a remote address.
func MigrateRepo(c *context.APIContext) {
	var f MigrateRepoForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.Error(err, "bind JSON")
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
		if !database.FromDomainUser(ctxUser).IsOrganization() || !database.FromDomainUser(ctxUser).IsOwnedBy(c.User.ID) {
			c.Status(http.StatusForbidden)
			return
		}
	}

	remoteAddr := strings.TrimPrefix(strings.TrimSpace(f.CloneAddr), "git clone ")

	if strings.HasPrefix(remoteAddr, "http://") || strings.HasPrefix(remoteAddr, "https://") {
		u, err := url.Parse(remoteAddr)
		if err == nil {
			if f.AuthUsername != "" || f.AuthPassword != "" {
				u.User = url.UserPassword(f.AuthUsername, f.AuthPassword)
			}
			remoteAddr = u.String()
		}
	}

	repo, err := database.MigrateRepository(c.User, database.FromDomainUser(ctxUser), database.MigrateRepoOptions{
		Name:          f.RepoName,
		Description:   f.Description,
		IsPrivate:     f.Private || conf.Repository.ForcePrivate,
		IsUnlisted:    f.Unlisted,
		IsArchived:    f.IsArchived,
		Website:       f.Website,
		DefaultBranch: f.DefaultBranch,
		RemoteAddr:    remoteAddr,
	})

	if err != nil {
		if repo != nil {
			if errDelete := database.DeleteRepository(ctxUser.ID, repo.ID); errDelete != nil {
				log.Error("DeleteRepository: %v", errDelete)
			}
		}

		if database.IsErrRepoAlreadyExist(err) {
			c.StatusError(http.StatusConflict, err.Error(), err)
		} else {
			c.Error(err, "migrate repository")
		}
		return
	}

	log.Trace("Repository migrated [%d]: %s/%s", repo.ID, ctxUser.Name, f.RepoName)

	// Analyze languages after migration
	if err := reposervice.UpdateRepositoryLanguages(repo); err != nil {
		log.Error("UpdateRepositoryLanguages after migration: %v", err)
	}

	c.JSON(http.StatusCreated, repo.APIFormat(database.FromDomainUser(ctxUser)))
}
