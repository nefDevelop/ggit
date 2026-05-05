package v2

import (
	"net/http"
	"strings"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/gitx"
)

// GetComparison returns the comparison between two branches.
func GetComparison(c *context.APIContext) {
	info := c.Params("*")
	parts := strings.Split(info, "...")
	if len(parts) != 2 {
		c.ErrorStatus(http.StatusBadRequest, nil)
		return
	}

	baseBranch := parts[0]
	headBranch := parts[1]

	baseRepo := c.Repo.Repository
	baseRepoPath := baseRepo.RepoPath()

	// If branches are the same and it's the same repository, return empty result
	if baseBranch == headBranch {
		c.JSON(http.StatusOK, map[string]any{
			"base_branch": baseBranch,
			"head_branch": headBranch,
			"merge_base":  "",
			"commits":     []any{},
			"num_commits": 0,
			"diff":        nil,
		})
		return
	}

	// For simplicity in this first version, we assume same repository.
	// In the future, we can expand to fork comparison.
	meta, err := gitx.Module.PullRequestMeta(baseRepoPath, baseRepoPath, headBranch, baseBranch)
	if err != nil {
		c.NotFoundOrError(gitx.NewError(err), "get pull request meta")
		return
	}

	if meta.MergeBase == "" {
		c.StatusError(http.StatusNotFound, "No common ancestor found", nil)
		return
	}

	// Get Diff
	diff, err := gitx.RepoDiff(c.Repo.GitRepo,
		meta.MergeBase, 500, 5000, 500,
		git.DiffOptions{Base: baseBranch}, // simplified
	)
	if err != nil {
		c.Error(err, "get diff")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"base_branch": baseBranch,
		"head_branch": headBranch,
		"merge_base":  meta.MergeBase,
		"commits":     matchUsersWithCommitEmails(c.Req.Context(), meta.Commits),
		"num_commits": len(meta.Commits),
		"diff":        diff,
	})
}
