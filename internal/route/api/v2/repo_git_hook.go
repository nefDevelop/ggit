package v2

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/nef734/ggit/internal/context"
)

type GitHook struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Content  string `json:"content"`
}

// ListRepoGitHooks returns a list of git hooks for the repository.
func ListRepoGitHooks(c *context.APIContext) {
	hooks := []string{"pre-receive", "update", "post-receive"}
	res := make([]*GitHook, len(hooks))
	for i, name := range hooks {
		hookPath := filepath.Join(c.Repo.Repository.RepoPath(), "custom_hooks", name)
		isActive := false
		content := ""
		if info, err := os.Stat(hookPath); err == nil && !info.IsDir() {
			isActive = true
			data, _ := os.ReadFile(hookPath)
			content = string(data)
		}
		res[i] = &GitHook{
			Name:     name,
			IsActive: isActive,
			Content:  content,
		}
	}
	c.JSON(http.StatusOK, res)
}

// GetRepoGitHook returns a specific git hook for the repository.
func GetRepoGitHook(c *context.APIContext) {
	name := c.Params("name")
	hookPath := filepath.Join(c.Repo.Repository.RepoPath(), "custom_hooks", name)
	
	content := ""
	isActive := false
	if info, err := os.Stat(hookPath); err == nil && !info.IsDir() {
		isActive = true
		data, _ := os.ReadFile(hookPath)
		content = string(data)
	}
	
	c.JSON(http.StatusOK, &GitHook{
		Name:     name,
		IsActive: isActive,
		Content:  content,
	})
}

// UpdateRepoGitHook updates a specific git hook for the repository.
func UpdateRepoGitHook(c *context.APIContext) {
	name := c.Params("name")
	var f struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}
	
	customHooksDir := filepath.Join(c.Repo.Repository.RepoPath(), "custom_hooks")
	if err := os.MkdirAll(customHooksDir, 0755); err != nil {
		c.Error(err, "create custom hooks directory")
		return
	}
	
	hookPath := filepath.Join(customHooksDir, name)
	if f.Content == "" {
		_ = os.Remove(hookPath)
	} else {
		if err := os.WriteFile(hookPath, []byte(f.Content), 0755); err != nil {
			c.Error(err, "write hook file")
			return
		}
	}
	
	c.NoContent()
}
