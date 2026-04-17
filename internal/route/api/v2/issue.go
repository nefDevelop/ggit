package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/markup"
	apiv1types "gogs.io/gogs/internal/route/api/v1/types"
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
		RepoID:      c.Repo.Repository.ID,
		IsPull:      isPullList,
		Page:        page,
		Labels:      c.Query("labels"),
		MilestoneID: c.QueryInt64("milestone_id"),
		AssigneeID:  c.QueryInt64("assignee_id"),
		PosterID:    c.QueryInt64("poster_id"),
		SortType:    c.Query("sort"),
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

	apiIssues := make([]*apiv1types.Issue, len(issues))
	for i, issue := range issues {
		apiIssues[i] = issue.APIFormat()
		apiIssues[i].BodyHTML = string(markup.Render(markup.TypeMarkdown, issue.Content, c.Repo.RepoLink, c.Repo.Repository.ComposeMetas()))
	}

	c.JSON(http.StatusOK, apiIssues)
}

// CreateIssue creates a new issue for the repository.
func CreateIssue(c *context.APIContext) {
	var f issueservice.CreateIssueOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	issue := &database.Issue{
		RepoID:      c.Repo.Repository.ID,
		PosterID:    c.User.ID,
		Poster:      c.User,
		Title:       f.Title,
		Content:     f.Content,
		MilestoneID: f.MilestoneID,
		AssigneeID:  f.AssigneeID,
	}

	if err := issueservice.Issues.CreateIssue(c.Req.Context(), c.Repo.Repository, issue, f); err != nil {
		c.Error(err, "create issue")
		return
	}

	apiIssue := issue.APIFormat()
	apiIssue.BodyHTML = string(markup.Render(markup.TypeMarkdown, issue.Content, c.Repo.RepoLink, c.Repo.Repository.ComposeMetas()))

	c.JSON(http.StatusCreated, apiIssue)
}
// GetIssue returns a single issue by its index.
func GetIssue(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	res := issue.APIFormat()
	res.BodyHTML = string(markup.Render(markup.TypeMarkdown, issue.Content, c.Repo.RepoLink, c.Repo.Repository.ComposeMetas()))

	c.JSON(http.StatusOK, res)
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

// GetIssueComments returns all comments for a specific issue.
func GetIssueComments(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	comments, err := database.GetCommentsByIssueID(issue.ID)
	if err != nil {
		c.Error(err, "get comments by issue ID")
		return
	}

	apiComments := make([]*apiv1types.IssueComment, len(comments))
	for i, comment := range comments {
		apiComments[i] = comment.APIFormat()
		apiComments[i].BodyHTML = string(markup.Render(markup.TypeMarkdown, comment.Content, c.Repo.RepoLink, c.Repo.Repository.ComposeMetas()))
	}

	c.JSON(http.StatusOK, apiComments)
}

// CreateIssueComment creates a new comment on an issue.
func CreateIssueComment(c *context.APIContext) {
	idx := c.ParamsInt64("index")
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, idx)
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	var f struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	comment, err := database.CreateIssueComment(c.User, c.Repo.Repository, issue, f.Body, nil)
	if err != nil {
		c.Error(err, "create comment")
		return
	}

	c.JSON(http.StatusCreated, comment)
}
