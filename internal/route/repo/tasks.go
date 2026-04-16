package repo

import (
	"net/http"

	"github.com/cockroachdb/errors"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/cryptox"
	"gogs.io/gogs/internal/database"
)

func TriggerTask(c *context.Context) {
	branch := c.Query("branch")
	pusherID := c.QueryInt64("pusher")
	secret := c.Query("secret")
	if branch == "" || pusherID <= 0 || secret == "" {
		c.Error(errors.New("incomplete branch, pusher or secret"), "trigger task")
		return
	}

	username := c.Param("username")
	reponame := c.Param("reponame")

	owner, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		if database.IsErrUserNotExist(err) {
			c.Error(err, "owner does not exist")
		} else {
			c.Status(http.StatusInternalServerError)
			log.Error("Failed to get user [name: %s]: %v", username, err)
		}
		return
	}

	// 🚨 SECURITY: No need to check existence of the repository if the client
	// can't even get the valid secret. Mostly likely not a legitimate request.
	if secret != cryptox.MD5(owner.Salt) {
		c.Error(errors.New("invalid secret"), "valid secret")
		return
	}

	repo, err := database.Handle.Repositories().GetByName(c.Req.Context(), owner.ID, reponame)
	if err != nil {
		if database.IsErrRepoNotExist(err) {
			c.Error(err, "repository does not exist")
		} else {
			c.Status(http.StatusInternalServerError)
			log.Error("Failed to get repository [owner_id: %d, name: %s]: %v", owner.ID, reponame, err)
		}
		return
	}

	pusher, err := database.Handle.Users().GetByID(c.Req.Context(), pusherID)
	if err != nil {
		if database.IsErrUserNotExist(err) {
			c.Error(err, "pusher does not exist")
		} else {
			c.Status(http.StatusInternalServerError)
			log.Error("Failed to get user [id: %d]: %v", pusherID, err)
		}
		return
	}

	log.Trace("TriggerTask: %s/%s@%s by %q", owner.Name, repo.Name, branch, pusher.Name)

	go database.HookQueue.Add(repo.ID)
	go database.AddTestPullRequestTask(pusher, repo.ID, branch, true)
	c.Status(http.StatusAccepted)
}
