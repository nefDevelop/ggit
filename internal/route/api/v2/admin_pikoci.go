package v2

import (
	stdcontext "context"
	"net/http"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	pikociservice "github.com/nef734/ggit/internal/services/pikoci"
)

const pikoCIHealthTimeout = 5 * time.Second

// GetAdminPikoCIStatus devuelve la configuración visible y la salud del
// sidecar. Si PikoCI no responde se informa reachable=false en lugar de
// devolver un error: el panel necesita justamente mostrar ese estado.
func GetAdminPikoCIStatus(c *context.APIContext) {
	res := map[string]any{
		"enabled":    conf.PikoCI.Enabled,
		"mode":       conf.PikoCI.Mode,
		"url":        conf.PikoCI.URL,
		"public_url": conf.PikoCI.PublicURL,
	}
	if !conf.PikoCI.Enabled {
		c.JSON(http.StatusOK, res)
		return
	}

	ctx, cancel := stdcontext.WithTimeout(c.Req.Context(), pikoCIHealthTimeout)
	defer cancel()

	health, err := pikociservice.CI.Health(ctx)
	if err != nil {
		res["reachable"] = false
		res["error"] = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res["reachable"] = health.Status == "ok"
	res["health"] = health
	c.JSON(http.StatusOK, res)
}

// GetAdminPikoCIPipelineTemplate devuelve la plantilla HCL que ggit publica en
// PikoCI, con el bloque que reporta estados a ggit.
func GetAdminPikoCIPipelineTemplate(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]string{
		"template": pikociservice.PipelineTemplate,
	})
}

type pikoCIRepoPayload struct {
	PipelineName string `json:"pipeline_name" binding:"required"`
	WebhookToken string `json:"webhook_token" binding:"required"`
}

// UpdateAdminPikoCIRepo vincula un repositorio con su pipeline de PikoCI.
// Hasta que exista el aprovisionamiento automático, el administrador copia
// aquí el token del webhook del resource.
func UpdateAdminPikoCIRepo(c *context.APIContext) {
	repoID := c.ParamsInt64("repoID")
	if _, err := database.Handle.Repositories().GetByID(c.Req.Context(), repoID); err != nil {
		c.NotFoundOrError(err, "get repository by ID")
		return
	}

	var f pikoCIRepoPayload
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	link := &database.PikoCIRepo{
		RepoID:       repoID,
		PipelineName: f.PipelineName,
		WebhookToken: f.WebhookToken,
	}
	if err := database.Handle.PikoCIRepos().Upsert(c.Req.Context(), link); err != nil {
		c.Error(err, "save pikoci repository link")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":        "ok",
		"repo_id":       repoID,
		"pipeline_name": f.PipelineName,
	})
}

// SyncAdminPikoCIRepo publica o republica el pipeline del repositorio con las
// variables de ggit y guarda el webhook token resultante.
func SyncAdminPikoCIRepo(c *context.APIContext) {
	if !conf.PikoCI.Enabled {
		c.StatusError(http.StatusBadRequest, "PikoCI is disabled", nil)
		return
	}

	repoID := c.ParamsInt64("repoID")
	repo, err := database.Handle.Repositories().GetByID(c.Req.Context(), repoID)
	if err != nil {
		c.NotFoundOrError(err, "get repository by ID")
		return
	}
	owner, err := database.Handle.Users().GetByID(c.Req.Context(), repo.OwnerID)
	if err != nil {
		c.NotFoundOrError(err, "get repository owner")
		return
	}

	branch := repo.DefaultBranch
	if branch == "" {
		branch = "main"
	}

	link, err := pikociservice.CI.SyncRepoPipeline(c.Req.Context(), repo.ID, owner.Name, repo.Name, branch)
	if err != nil {
		c.Error(err, "sync pikoci pipeline")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":   "ok",
		"repo_id":  repoID,
		"pipeline": link.PipelineName,
	})
}

// DeleteAdminPikoCIRepo borra el pipeline en PikoCI (si está habilitado) y el
// vínculo local.
func DeleteAdminPikoCIRepo(c *context.APIContext) {
	repoID := c.ParamsInt64("repoID")

	if conf.PikoCI.Enabled {
		if err := pikociservice.CI.DisableRepoPipeline(c.Req.Context(), repoID, ""); err != nil {
			c.Error(err, "disable pikoci pipeline")
			return
		}
	} else if err := database.Handle.PikoCIRepos().Delete(c.Req.Context(), repoID); err != nil {
		c.Error(err, "delete pikoci repository link")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"repo_id": repoID,
	})
}
