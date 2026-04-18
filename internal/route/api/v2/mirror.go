package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
)

// SyncMirrorRepo triggers an immediate sync for a mirror repository.
func SyncMirrorRepo(c *context.APIContext) {
	if !c.Repo.Repository.IsMirror {
		c.StatusError(http.StatusBadRequest, "Repository is not a mirror", nil)
		return
	}

	database.MirrorQueue.Add(c.Repo.Repository.ID)
	c.Status(http.StatusOK)
}
