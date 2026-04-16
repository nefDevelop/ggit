package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
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

	for _, issue := range issues {
		_ = issue.LoadAttributes()
	}

	c.JSON(http.StatusOK, issues)
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

	c.JSON(http.StatusOK, issue)
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
