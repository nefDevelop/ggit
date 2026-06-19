package v2

import (
	"net/http"
	"os"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/gitx"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	issueservice "github.com/nef734/ggit/internal/services/issue"
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

	if err = issue.PullRequest.CheckBranchProtection(c.Req.Context(), c.User); err != nil {
		c.StatusError(http.StatusForbidden, err.Error(), err)
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
		Type:         database.PullRequestTypeGgit,
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

// ListPullReviews returns all reviews for a pull request.
func ListPullReviews(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue")
		return
	}

	reviews, err := database.Handle.PullReviews().ListByIssueID(c.Req.Context(), issue.ID)
	if err != nil {
		c.Error(err, "list reviews")
		return
	}

	// We should probably load reviewers too.
	apiReviews := make([]*apiv2types.PullReviewResponse, len(reviews))
	for i, r := range reviews {
		du, _ := database.Handle.Users().GetByID(c.Req.Context(), r.ReviewerID)
		r.Reviewer = database.FromDomainUser(du)
		r.Comments, _ = database.Handle.PullReviewComments().ListByReviewID(c.Req.Context(), r.ID)
		apiReviews[i] = r.APIFormat()
	}

	c.JSON(http.StatusOK, apiReviews)
}

// ResolvePullReviewComment resolves or unresolves a review comment.
func ResolvePullReviewComment(c *context.APIContext) {
	commentID := c.ParamsInt64("id")
	var f struct {
		Resolved bool `json:"resolved"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := database.Handle.PullReviewComments().SetResolved(c.Req.Context(), commentID, f.Resolved); err != nil {
		c.Error(err, "set resolved")
		return
	}

	c.NoContent()
}

// CreatePullReview creates a new review for a pull request.
func CreatePullReview(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue")
		return
	}

	var f struct {
		Type    int    `json:"type" binding:"required"` // 1: Approve, 2: Comment, 3: Reject
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	review := &database.PullReview{
		Type:       database.PullReviewType(f.Type),
		ReviewerID: c.User.ID,
		IssueID:    issue.ID,
		Content:    f.Content,
	}

	if err := database.Handle.PullReviews().Create(c.Req.Context(), review); err != nil {
		c.Error(err, "create review")
		return
	}

	c.JSON(http.StatusCreated, review)
}

// AddPullReviewComment adds an inline comment to a review.
func AddPullReviewComment(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	reviewID := c.ParamsInt64("id")

	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue")
		return
	}

	var f struct {
		TreePath string `json:"tree_path" binding:"required"`
		Line     int64  `json:"line"`
		CommitID string `json:"commit_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	comment := &database.PullReviewComment{
		ReviewID: reviewID,
		IssueID:  issue.ID,
		TreePath: f.TreePath,
		Line:     f.Line,
		CommitID: f.CommitID,
		Content:  f.Content,
	}

	if err := database.Handle.PullReviewComments().Create(c.Req.Context(), comment); err != nil {
		c.Error(err, "create review comment")
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetPullRequestProtection returns the branch protection status for a pull request.
func GetPullRequestProtection(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue")
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

	pb, err := database.GetProtectBranchOfRepoByName(issue.Repo.ID, issue.PullRequest.BaseBranch)
	if err != nil && !database.IsErrBranchNotExist(err) {
		c.Error(err, "get protect branch")
		return
	}

	stats, err := database.Handle.PullReviews().GetReviewStatsByIssueID(c.Req.Context(), issue.ID)
	if err != nil {
		c.Error(err, "get review stats")
		return
	}

	// Count unresolved conversations
	unresolvedCount := 0
	comments, _ := database.Handle.PullReviewComments().ListByIssueID(c.Req.Context(), issue.ID)
	for _, comment := range comments {
		if !comment.IsResolved {
			unresolvedCount++
		}
	}

	combinedStatus := database.Handle.CommitStatuses().GetCombinedStatusBySHA(c.Req.Context(), issue.Repo.ID, issue.PullRequest.HeadCommitID)

	requiredApprovals := 0
	if pb != nil && pb.Protected {
		requiredApprovals = pb.RequiredApprovals
	}
	if requiredApprovals == 0 {
		requiredApprovals = issue.Repo.RequiredApprovals
	}

	err = issue.PullRequest.CheckMergeChecks(c.Req.Context(), c.User)
	mergeable := err == nil

	c.JSON(http.StatusOK, map[string]any{
		"required_approvals": requiredApprovals,
		"approvals":          stats.Approvals,
		"rejections":         stats.Rejections,
		"unresolved_count":   unresolvedCount,
		"combined_status":    string(combinedStatus),
		"is_protected":       pb != nil && pb.Protected,
		"mergeable":          mergeable,
		"merge_error": func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	})
}
