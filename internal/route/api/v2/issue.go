package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	issueservice "github.com/nef734/ggit/internal/services/issue"
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

	apiIssues := make([]*apiv2types.IssueResponse, len(issues))
	for i, issue := range issues {
		apiIssues[i] = issue.APIFormat()
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
	if !f.CreatedAt.IsZero() {
		issue.CreatedUnix = f.CreatedAt.Unix()
		issue.UpdatedUnix = issue.CreatedUnix
	}

	if err := issueservice.Issues.CreateIssue(c.Req.Context(), c.Repo.Repository, issue, f); err != nil {
		c.Error(err, "create issue")
		return
	}

	// The index is assigned inside NewIssue, on a model reached through three
	// conversions: the service takes issue.ToDomain() and the store rebuilds it
	// with FromDomainIssue. That copy is the one that gets the number, so the
	// struct held here still reads zero and the response used to say
	// number: 0. Re-reading returns what was actually persisted.
	created, err := database.GetIssueByID(issue.ID)
	if err != nil {
		c.Error(err, "get created issue")
		return
	}

	c.JSON(http.StatusCreated, created.APIFormat())
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

	apiComments := make([]*apiv2types.IssueCommentResponse, len(comments))
	for i, comment := range comments {
		apiComments[i] = comment.APIFormat()
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

	// The raw model was serialised here, and it does not match what the rest of
	// the API returns: Comment.Content has no json tag so it came out as
	// "Content", and Poster is json:"-". The client reads "body" and "user", so
	// the comment was posted successfully and then rendered empty.
	c.JSON(http.StatusCreated, comment.APIFormat())
}
