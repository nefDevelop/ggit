package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
)

// GetRepoTemplates returns the list of available .gitignore, License and Readme templates.
func GetRepoTemplates(c *context.APIContext) {
	c.JSON(http.StatusOK, map[string]any{
		"gitignores": database.Gitignores,
		"licenses":   database.Licenses,
		"readmes":    database.Readmes,
	})
}
