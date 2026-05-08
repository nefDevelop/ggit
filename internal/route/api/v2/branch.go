package v2

import (
	"net/http"
	"time"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

type Branch struct {
	Name        string `json:"name"`
	Commit      any    `json:"commit"`
	IsProtected bool   `json:"is_protected"`
}

// ListRepoBranches returns a list of branches for the repository.
func ListRepoBranches(c *context.APIContext) {
	if c.Repo.Repository.IsBare {
		c.JSON(http.StatusOK, map[string]any{
			"default": &Branch{
				Name:        c.Repo.Repository.DefaultBranch,
				Commit:      nil,
				IsProtected: false,
			},
			"active": []*Branch{},
			"stale":  []*Branch{},
		})
		return
	}

	rawBranches, err := c.Repo.Repository.GetBranches()
	if err != nil {
		c.Error(err, "get branches")
		return
	}

	protectBranches, err := database.GetProtectBranchesByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get protect branches")
		return
	}

	branches := make([]*Branch, len(rawBranches))
	for i := range rawBranches {
		commit, err := rawBranches[i].GetCommit()
		if err != nil {
			c.Error(err, "get commit")
			return
		}

		branches[i] = &Branch{
			Name: rawBranches[i].Name,
			Commit: map[string]any{
				"id":      commit.ID.String(),
				"summary": commit.Summary(),
				"committer": map[string]any{
					"name":  commit.Committer.Name,
					"email": commit.Committer.Email,
					"when":  commit.Committer.When,
				},
			},
		}

		for _, pb := range protectBranches {
			if branches[i].Name == pb.Name {
				branches[i].IsProtected = true
				break
			}
		}
	}

	now := time.Now()
	active := make([]*Branch, 0)
	stale := make([]*Branch, 0)
	var defaultBranch *Branch

	for _, b := range branches {
		if b.Name == c.Repo.Repository.DefaultBranch {
			defaultBranch = b
			continue
		}
		if b.Commit.(map[string]any)["committer"].(map[string]any)["when"].(time.Time).Add(30 * 24 * time.Hour).After(now) {
			active = append(active, b)
		} else {
			stale = append(stale, b)
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"default": defaultBranch,
		"active":  active,
		"stale":   stale,
	})
}

// DeleteRepoBranch deletes a branch from the repository.
func DeleteRepoBranch(c *context.APIContext) {
	branchName := c.Params("branch")

	if branchName == c.Repo.Repository.DefaultBranch {
		c.StatusError(http.StatusBadRequest, "Cannot delete default branch.", nil)
		return
	}

	protectBranch, err := database.GetProtectBranchOfRepoByName(c.Repo.Repository.ID, branchName)
	if err != nil && !database.IsErrBranchNotExist(err) {
		c.Error(err, "get protect branch")
		return
	}
	if protectBranch != nil && protectBranch.Protected {
		c.StatusError(http.StatusBadRequest, "Cannot delete protected branch.", nil)
		return
	}

	if c.Repo.GitRepo == nil {
		c.StatusError(http.StatusInternalServerError, "Git repository is not initialized.", nil)
		return
	}

	if err := c.Repo.GitRepo.DeleteBranch(branchName, git.DeleteBranchOptions{Force: true}); err != nil {
		c.Error(err, "delete branch")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// CreateRepoBranch creates a new branch in the repository.
func CreateRepoBranch(c *context.APIContext) {
	var f struct {
		OldRefName string `json:"old_ref_name" binding:"required"`
		NewRefName string `json:"new_ref_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if c.Repo.GitRepo == nil {
		c.StatusError(http.StatusInternalServerError, "Git repository is not initialized.", nil)
		return
	}

	if !c.Repo.GitRepo.HasBranch(f.OldRefName) && !c.Repo.GitRepo.HasTag(f.OldRefName) && len(f.OldRefName) != 40 {
		c.StatusError(http.StatusUnprocessableEntity, "Old ref name does not exist.", nil)
		return
	}

	if c.Repo.GitRepo.HasBranch(f.NewRefName) {
		c.StatusError(http.StatusUnprocessableEntity, "Branch already exists.", nil)
		return
	}

	if err := c.Repo.Repository.CreateNewBranch(f.OldRefName, f.NewRefName); err != nil {
		c.Error(err, "create branch")
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
		"name":   f.NewRefName,
	})
}

// GetBranchProtection returns protection rules for a branch.
func GetBranchProtection(c *context.APIContext) {
	branchName := c.Params("branch")
	pb, err := database.GetProtectBranchOfRepoByName(c.Repo.Repository.ID, branchName)
	if err != nil {
		if database.IsErrBranchNotExist(err) {
			c.JSON(http.StatusOK, map[string]any{"protected": false})
		} else {
			c.Error(err, "get protect branch")
		}
		return
	}
	c.JSON(http.StatusOK, pb)
}

// UpdateBranchProtection updates protection rules for a branch.
func UpdateBranchProtection(c *context.APIContext) {
	branchName := c.Params("branch")
	var f struct {
		Protected              bool   `json:"protected"`
		RequirePullRequest     bool   `json:"require_pull_request"`
		EnableWhitelist        bool   `json:"enable_whitelist"`
		WhitelistUserIDs       string `json:"whitelist_user_ids"`
		WhitelistTeamIDs       string `json:"whitelist_team_ids"`
		RequiredApprovals      int    `json:"required_approvals"`
		EnableMergeWhitelist   bool   `json:"enable_merge_whitelist"`
		MergeWhitelistUserIDs  string `json:"merge_whitelist_user_ids"`
		MergeWhitelistTeamIDs  string `json:"merge_whitelist_team_ids"`
		BlockOnOutdatedBranch bool   `json:"block_on_outdated_branch"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	pb, err := database.GetProtectBranchOfRepoByName(c.Repo.Repository.ID, branchName)
	if err != nil && !database.IsErrBranchNotExist(err) {
		c.Error(err, "get protect branch")
		return
	}

	if pb == nil {
		pb = &database.ProtectBranch{
			RepoID: c.Repo.Repository.ID,
			Name:   branchName,
		}
	}

	pb.Protected = f.Protected
	pb.RequirePullRequest = f.RequirePullRequest
	pb.EnableWhitelist = f.EnableWhitelist
	pb.WhitelistUserIDs = f.WhitelistUserIDs
	pb.WhitelistTeamIDs = f.WhitelistTeamIDs
	pb.RequiredApprovals = f.RequiredApprovals
	pb.EnableMergeWhitelist = f.EnableMergeWhitelist
	pb.MergeWhitelistUserIDs = f.MergeWhitelistUserIDs
	pb.MergeWhitelistTeamIDs = f.MergeWhitelistTeamIDs
	pb.BlockOnOutdatedBranch = f.BlockOnOutdatedBranch

	if err := database.UpdateProtectBranch(pb); err != nil {
		c.Error(err, "update protect branch")
		return
	}

	c.JSON(http.StatusOK, pb)
}

// DeleteBranchProtection deletes protection rules for a branch.
func DeleteBranchProtection(c *context.APIContext) {
	branchName := c.Params("branch")
	if err := database.DeleteProtectBranch(c.Repo.Repository.ID, branchName); err != nil {
		c.Error(err, "delete protect branch")
		return
	}
	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
