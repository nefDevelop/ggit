package v2

import (
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/gitx"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// GetComparison returns the comparison between two branches.
func GetComparison(c *context.APIContext) {
	info := c.Param("info")
	parts := strings.Split(info, "...")
	if len(parts) != 2 {
		c.ErrorStatus(http.StatusBadRequest, errors.New("invalid comparison format"))
		return
	}

	baseBranch := strings.Trim(parts[0], " /")
	headBranch := strings.Trim(parts[1], " /")

	baseRepo := c.Repo.Repository
	baseRepoPath := baseRepo.RepoPath()

	// Gather metadata for the UI (branches, forks, etc.)
	branches, _ := c.Repo.GitRepo.Branches()
	
	var allHeadRepos []*apiv2types.RepositoryResponse
	allHeadRepos = append(allHeadRepos, baseRepo.APIFormat(c.Repo.Owner))
	
	forks, _ := database.Handle.Repositories().ListForks(c.Req.Context(), baseRepo.ID, 1, 100)
	for _, fork := range forks {
		owner, _ := database.Handle.Users().GetByID(c.Req.Context(), fork.OwnerID)
		allHeadRepos = append(allHeadRepos, fork.APIFormat(owner))
	}

	resp := map[string]any{
		"base_repo":      baseRepo.APIFormat(c.Repo.Owner),
		"head_repo":      baseRepo.APIFormat(c.Repo.Owner),
		"all_head_repos": allHeadRepos,
		"base_branches":  branches,
		"head_branches":  branches,
		"base_branch":    baseBranch,
		"head_branch":    headBranch,
		"is_same_repo":   true,
	}

	// If branches are the same, return empty result with metadata
	if baseBranch == headBranch {
		resp["merge_base"] = ""
		resp["commits"] = []any{}
		resp["num_commits"] = 0
		resp["diff"] = nil
		c.JSON(http.StatusOK, resp)
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

	resp["merge_base"] = meta.MergeBase
	resp["commits"] = matchUsersWithCommitEmails(c.Req.Context(), meta.Commits)
	resp["num_commits"] = len(meta.Commits)
	resp["diff"] = diff

	c.JSON(http.StatusOK, resp)
}
