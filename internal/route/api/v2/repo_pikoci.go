package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	pikociservice "github.com/nef734/ggit/internal/services/pikoci"
)

// GetRepoPikoCI informa si el repositorio tiene CI y con qué pipeline.
func GetRepoPikoCI(c *context.APIContext) {
	res := map[string]any{
		"enabled":    conf.PikoCI.Enabled,
		"public_url": conf.PikoCI.PublicURL,
		"linked":     false,
	}

	link, err := database.Handle.PikoCIRepos().GetByRepoID(c.Req.Context(), c.Repo.Repository.ID)
	switch {
	case err == nil:
		res["linked"] = true
		res["pipeline"] = link.PipelineName
	case database.IsErrPikoCIRepoNotExist(err):
	default:
		c.Error(err, "get pikoci link")
		return
	}

	c.JSON(http.StatusOK, res)
}

// EnableRepoPikoCI publica el pipeline del repositorio en PikoCI.
func EnableRepoPikoCI(c *context.APIContext) {
	if !conf.PikoCI.Enabled {
		c.StatusError(http.StatusBadRequest, "PikoCI is disabled", nil)
		return
	}

	branch := c.Repo.Repository.DefaultBranch
	if branch == "" {
		branch = "main"
	}

	link, err := pikociservice.CI.SyncRepoPipeline(c.Req.Context(), c.Repo.Repository.ID,
		c.Repo.Owner.Name, c.Repo.Repository.Name, branch)
	if err != nil {
		c.Error(err, "sync pikoci pipeline")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":   "ok",
		"pipeline": link.PipelineName,
	})
}

// DisableRepoPikoCI borra el pipeline en PikoCI y el vínculo local.
func DisableRepoPikoCI(c *context.APIContext) {
	repoID := c.Repo.Repository.ID

	if conf.PikoCI.Enabled {
		if err := pikociservice.CI.DisableRepoPipeline(c.Req.Context(), repoID, ""); err != nil {
			c.Error(err, "disable pikoci pipeline")
			return
		}
	} else if err := database.Handle.PikoCIRepos().Delete(c.Req.Context(), repoID); err != nil {
		c.Error(err, "delete pikoci repository link")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
