package v2

import (
	"net/http"
	"time"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

func toIssueMilestone(m *database.Milestone) *apiv2types.IssueMilestoneResponse {
	state := apiv2types.IssueStateOpen
	if m.IsClosed {
		state = apiv2types.IssueStateClosed
	}

	return &apiv2types.IssueMilestoneResponse{
		ID:           m.ID,
		Title:        m.Name,
		Description:  m.Content,
		State:        state,
		OpenIssues:   m.NumOpenIssues,
		ClosedIssues: m.NumClosedIssues,
		Deadline:     &m.Deadline,
	}
}

// ListMilestones returns a list of milestones for the repository.
func ListRepoMilestones(c *context.APIContext) {
	milestones, err := database.GetMilestonesByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get milestones by repository ID")
		return
	}

	apiMilestones := make([]*apiv2types.IssueMilestoneResponse, len(milestones))
	for i := range milestones {
		apiMilestones[i] = toIssueMilestone(milestones[i])
	}
	c.JSON(http.StatusOK, apiMilestones)
}

// GetMilestone returns a single milestone for the repository.
func GetRepoMilestone(c *context.APIContext) {
	milestone, err := database.GetMilestoneByRepoID(c.Repo.Repository.ID, c.ParamsInt64("id"))
	if err != nil {
		c.NotFoundOrError(err, "get milestone by repository ID")
		return
	}

	c.JSON(http.StatusOK, toIssueMilestone(milestone))
}

type MilestoneOptions struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"due_on"`
	State       *string    `json:"state"`
}

// CreateMilestone creates a new milestone for the repository.
func CreateRepoMilestone(c *context.APIContext) {
	var f MilestoneOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if f.Deadline == nil {
		defaultDeadline, _ := time.ParseInLocation("2006-01-02", "9999-12-31", time.Local)
		f.Deadline = &defaultDeadline
	}

	milestone := &database.Milestone{
		RepoID:   c.Repo.Repository.ID,
		Name:     f.Title,
		Content:  f.Description,
		Deadline: *f.Deadline,
	}

	if err := database.NewMilestone(milestone); err != nil {
		c.Error(err, "new milestone")
		return
	}

	c.JSON(http.StatusCreated, toIssueMilestone(milestone))
}

// UpdateMilestone updates an existing milestone.
func UpdateRepoMilestone(c *context.APIContext) {
	milestone, err := database.GetMilestoneByRepoID(c.Repo.Repository.ID, c.ParamsInt64("id"))
	if err != nil {
		c.NotFoundOrError(err, "get milestone by repository ID")
		return
	}

	var f MilestoneOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if len(f.Title) > 0 {
		milestone.Name = f.Title
	}
	if len(f.Description) > 0 {
		milestone.Content = f.Description
	}
	if f.Deadline != nil && !f.Deadline.IsZero() {
		milestone.Deadline = *f.Deadline
	}

	if f.State != nil {
		if err = milestone.ChangeStatus(apiv2types.IssueStateClosed == apiv2types.IssueStateType(*f.State)); err != nil {
			c.Error(err, "change status")
			return
		}
	} else if err = database.UpdateMilestone(milestone); err != nil {
		c.Error(err, "update milestone")
		return
	}

	c.JSON(http.StatusOK, toIssueMilestone(milestone))
}

// DeleteMilestone deletes a milestone from the repository.
func DeleteRepoMilestone(c *context.APIContext) {
	if err := database.DeleteMilestoneOfRepoByID(c.Repo.Repository.ID, c.ParamsInt64("id")); err != nil {
		c.Error(err, "delete milestone of repository by ID")
		return
	}

	c.NoContent()
}
