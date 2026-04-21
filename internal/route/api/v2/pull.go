package v2

import (
	"net/http"
	"os"
	"time"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/gitx"
	apiv2types "gogs.io/gogs/internal/route/api/v2/types"
	issueservice "gogs.io/gogs/internal/services/issue"
)

// ListPullRequests returns a list of pull requests for the repository.
func ListPullRequests(c *context.APIContext) {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	opts := &database.IssuesOptions{
		RepoID: c.Repo.Repository.ID,
		IsPull: true,
		Page:   page,
	}

	state := c.Query("state")
	if state == "closed" {
		opts.IsClosed = true
	}

	issues, err := database.Issues(opts)
	if err != nil {
		c.Error(err, "list pull requests")
		return
	}

	apiPullRequests := make([]*apiv2types.PullRequestResponse, len(issues))
	for i, issue := range issues {
		_ = issue.LoadAttributes()
		apiPullRequests[i] = issue.PullRequest.APIFormat()
	}

	c.JSON(http.StatusOK, apiPullRequests)
}

// GetPullRequest returns a single pull request by its index.
func GetPullRequest(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	if !issue.IsPull {
		c.NotFound()
		return
	}

	if err = issue.LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	c.JSON(http.StatusOK, issue.PullRequest.APIFormat())
}

// MergePullRequest merges a pull request.
func MergePullRequest(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	if !issue.IsPull {
		c.NotFound()
		return
	}

	if err = issue.LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	if issue.PullRequest.HasMerged {
		c.StatusError(http.StatusUnprocessableEntity, "Pull request has already been merged.", nil)
		return
	}

	if !issue.PullRequest.CanAutoMerge() {
		c.StatusError(http.StatusUnprocessableEntity, "Pull request is not mergeable.", nil)
		return
	}

	if err = issue.PullRequest.Merge(c.User, c.Repo.GitRepo, database.MergeStyleRegular, ""); err != nil {
		c.Error(err, "merge pull request")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "merged",
	})
}

// GetPullRequestDiff returns the diff of a pull request.
func GetPullRequestDiff(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	// I should get the issue first.
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}
	
	pull, err := database.GetPullRequestByIssueID(issue.ID)
	if err != nil {
		c.NotFoundOrError(err, "get pull request by issue ID")
		return
	}

	patchPath, err := c.Repo.Repository.PatchPath(pull.Index)
	if err != nil {
		c.Error(err, "get patch path")
		return
	}

	f, err := os.Open(patchPath)
	if err != nil {
		c.Error(err, "open patch file")
		return
	}
	defer f.Close()

	diff, err := gitx.ParseDiff(f, conf.Git.MaxDiffFiles, conf.Git.MaxDiffLines, conf.Git.MaxDiffLineChars)
	if err != nil {
		c.Error(err, "parse diff")
		return
	}

	c.JSON(http.StatusOK, diff)
}

// CreatePullRequest creates a new pull request.
func CreatePullRequest(c *context.APIContext) {
	var f struct {
		Title      string    `json:"title" binding:"required"`
		Content    string    `json:"content"`
		BaseBranch string    `json:"base_branch" binding:"required"`
		HeadBranch string    `json:"head_branch" binding:"required"`
		CreatedAt  time.Time `json:"created_at"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	repo := c.Repo.Repository
	baseBranch := f.BaseBranch
	headBranch := f.HeadBranch

	// Check if branches exist
	if !c.Repo.GitRepo.HasBranch(baseBranch) {
		c.StatusError(http.StatusUnprocessableEntity, "Base branch does not exist.", nil)
		return
	}
	if !c.Repo.GitRepo.HasBranch(headBranch) {
		c.StatusError(http.StatusUnprocessableEntity, "Head branch does not exist.", nil)
		return
	}

	// For simplicity, same repo PRs for now.
	// In a real scenario, we'd handle fork repos too.
	meta, err := gitx.Module.PullRequestMeta(repo.RepoPath(), repo.RepoPath(), headBranch, baseBranch)
	if err != nil {
		c.Error(err, "get pull request meta")
		return
	}

	patch, err := c.Repo.GitRepo.DiffBinary(meta.MergeBase, headBranch)
	if err != nil {
		c.Error(err, "get patch")
		return
	}

	pullIssue := &database.Issue{
		RepoID:   repo.ID,
		Index:    repo.NextIssueIndex(),
		Title:    f.Title,
		PosterID: c.User.ID,
		Poster:   c.User,
		IsPull:   true,
		Content:  f.Content,
	}
	if !f.CreatedAt.IsZero() {
		pullIssue.CreatedUnix = f.CreatedAt.Unix()
		pullIssue.UpdatedUnix = pullIssue.CreatedUnix
	}
	pullRequest := &database.PullRequest{
		HeadRepoID:   repo.ID,
		BaseRepoID:   repo.ID,
		HeadUserName: c.Repo.Owner.Name,
		HeadBranch:   headBranch,
		BaseBranch:   baseBranch,
		HeadRepo:     repo,
		BaseRepo:     repo,
		MergeBase:    meta.MergeBase,
		Type:         database.PullRequestTypeGogs,
	}

	if err := issueservice.Issues.CreatePullRequest(c.Req.Context(), repo, pullIssue, pullRequest, issueservice.CreatePullRequestOptions{
		CreateIssueOptions: issueservice.CreateIssueOptions{
			RepoID:    repo.ID,
			PosterID:  c.User.ID,
			Title:     f.Title,
			Content:   f.Content,
			CreatedAt: f.CreatedAt,
		},
		HeadRepoID: repo.ID,
		BaseRepoID: repo.ID,
		HeadBranch: headBranch,
		BaseBranch: baseBranch,
		Patch:      patch,
	}); err != nil {
		c.Error(err, "create pull request")
		return
	}

	c.JSON(http.StatusCreated, pullIssue)
}
