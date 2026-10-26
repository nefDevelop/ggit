package v2

import (
	stdcontext "context"
	"net/http"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
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
