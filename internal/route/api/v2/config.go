package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// GetRepoTemplates returns the list of available .gitignore, License and Readme templates.
func GetRepoTemplates(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"gitignores":      database.Gitignores,
		"global_ignores": database.GlobalIgnores,
		"licenses":        database.Licenses,
		"readmes":         database.Readmes,
	})
}
