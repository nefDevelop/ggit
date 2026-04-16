package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	issueservice "gogs.io/gogs/internal/services/issue"
)

// ListIssues returns a list of issues for the repository.
func ListIssues(c *context.APIContext) {
	isPullList := c.QueryBool("pulls")
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	opts := &database.IssuesOptions{
		RepoID: c.Repo.Repository.ID,
		IsPull: isPullList,
		Page:   page,
	}

	state := c.Query("state")
	if state == "closed" {
		opts.IsClosed = true
	}

	issues, err := database.Issues(opts)
	if err != nil {
		c.Error(err, "list issues")
		return
	}

	c.JSON(http.StatusOK, issues)
}

// CreateIssue creates a new issue for the repository.
func CreateIssue(c *context.APIContext) {
	var f issueservice.CreateIssueOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	newIssue := &database.Issue{
		RepoID:      c.Repo.Repository.ID,
		PosterID:    c.User.ID,
		Poster:      c.User,
		Title:       f.Title,
		Content:     f.Content,
		MilestoneID: f.MilestoneID,
		AssigneeID:  f.AssigneeID,
	}

	err := issueservice.Issues.CreateIssue(c.Req.Context(), c.Repo.Repository, newIssue, f)
	if err != nil {
		c.Error(err, "create issue")
		return
	}

	c.JSON(http.StatusCreated, newIssue)
}

// GetIssue returns a single issue by its index.
func GetIssue(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	c.JSON(http.StatusOK, issue)
}

// UpdateIssue updates an existing issue.
func UpdateIssue(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	u, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	// Check permissions: only poster or repo admin/writer
	if !c.Repo.IsWriter() && c.User.ID != u.PosterID {
		c.Status(http.StatusForbidden)
		return
	}

	var f issueservice.UpdateIssueOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// Update fields
	if f.Title != nil {
		u.Title = *f.Title
	}
	if f.Content != nil {
		u.Content = *f.Content
	}
	if f.MilestoneID != nil {
		u.MilestoneID = *f.MilestoneID
	}
	if f.AssigneeID != nil {
		u.AssigneeID = *f.AssigneeID
	}
	if f.IsClosed != nil {
		u.IsClosed = *f.IsClosed
	}
	if f.Deadline != nil {
		u.Deadline = *f.Deadline
	}

	if err := database.UpdateIssue(u); err != nil {
		c.Error(err, "update issue")
		return
	}

	c.JSON(http.StatusOK, u)
}
