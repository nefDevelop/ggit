package v2

import (
	"encoding/json"
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/services/janitor"
)

// GetJanitorReport returns the latest janitor report for the repository.
func GetJanitorReport(c *context.APIContext) {
	repo := c.Repo.Repository
	if repo.JanitorReport == "" {
		c.JSONSuccess(map[string]any{
			"issues":      []any{},
			"total_waste": 0,
			"last_check":  repo.JanitorLastCheckUnix,
		})
		return
	}

	var report janitor.Report
	if err := json.Unmarshal([]byte(repo.JanitorReport), &report); err != nil {
		c.Error(err, "unmarshal janitor report")
		return
	}

	c.JSONSuccess(report)
}

// RunJanitorScan triggers a new janitor scan.
func RunJanitorScan(c *context.APIContext) {
	report, err := janitor.Scan(c.Req.Context(), c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "run janitor scan")
		return
	}
	c.JSONSuccess(report)
}

// PerformJanitorCleanup executes the suggested cleanup.
func PerformJanitorCleanup(c *context.APIContext) {
	var f struct {
		Paths []string `json:"paths" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	err := janitor.Cleanup(c.Req.Context(), c.Repo.Repository.ID, c.User, f.Paths)
	if err != nil {
		c.Error(err, "perform janitor cleanup")
		return
	}

	c.NoContent()
}
