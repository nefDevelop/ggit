package v2

import (
	"net/http"
	"strings"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// EditRepoFile handles web-based file creation or editing.
func EditRepoFile(c *context.APIContext) {
	var f struct {
		TreePath      string `json:"tree_path" binding:"required"`
		Content       string `json:"content"`
		CommitSummary string `json:"commit_summary"`
		CommitMessage string `json:"commit_message"`
		IsNewFile     bool   `json:"is_new_file"`
		LastCommitID  string `json:"last_commit_id"`
		Branch        string `json:"branch"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	message := strings.TrimSpace(f.CommitSummary)
	if message == "" {
		if f.IsNewFile {
			message = "Create " + f.TreePath
		} else {
			message = "Update " + f.TreePath
		}
	}
	if f.CommitMessage != "" {
		message += "\n\n" + f.CommitMessage
	}

	targetBranch := f.Branch
	if targetBranch == "" {
		targetBranch = c.Repo.BranchName
	}

	err := c.Repo.Repository.UpdateRepoFile(c.User, database.UpdateRepoFileOptions{
		OldBranch:   c.Repo.BranchName,
		NewBranch:   targetBranch,
		OldTreeName: f.TreePath,
		NewTreeName: f.TreePath,
		Message:     message,
		Content:     strings.ReplaceAll(f.Content, "\r", ""),
		IsNewFile:   f.IsNewFile,
	})
	if err != nil {
		c.Error(err, "update repo file")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteRepoFile handles web-based file deletion.
func DeleteRepoFile(c *context.APIContext) {
	var f struct {
		TreePath      string `json:"tree_path" binding:"required"`
		CommitSummary string `json:"commit_summary"`
		CommitMessage string `json:"commit_message"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	message := strings.TrimSpace(f.CommitSummary)
	if message == "" {
		message = "Delete " + f.TreePath
	}
	if f.CommitMessage != "" {
		message += "\n\n" + f.CommitMessage
	}

	err := c.Repo.Repository.DeleteRepoFile(c.User, database.DeleteRepoFileOptions{
		LastCommitID: c.Repo.CommitID,
		OldBranch:    c.Repo.BranchName,
		NewBranch:    c.Repo.BranchName,
		TreePath:     f.TreePath,
		Message:      message,
	})
	if err != nil {
		c.Error(err, "delete repo file")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
