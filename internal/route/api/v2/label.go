package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/route/api/v1/types"
)

func toIssueLabel(l *database.Label) *types.IssueLabel {
	return &types.IssueLabel{
		ID:    l.ID,
		Name:  l.Name,
		Color: l.Color,
	}
}

// ListLabels returns a list of labels for the repository.
func ListRepoLabels(c *context.APIContext) {
	labels, err := database.GetLabelsByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get labels by repository ID")
		return
	}

	apiLabels := make([]*types.IssueLabel, len(labels))
	for i := range labels {
		apiLabels[i] = toIssueLabel(labels[i])
	}
	c.JSON(http.StatusOK, apiLabels)
}

// GetLabel returns a single label for the repository.
func GetRepoLabel(c *context.APIContext) {
	label, err := database.GetLabelOfRepoByID(c.Repo.Repository.ID, c.ParamsInt64("id"))
	if err != nil {
		c.NotFoundOrError(err, "get label of repository by ID")
		return
	}

	c.JSON(http.StatusOK, toIssueLabel(label))
}

type LabelOptions struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

// CreateLabel creates a new label for the repository.
func CreateRepoLabel(c *context.APIContext) {
	var f LabelOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	label := &database.Label{
		RepoID: c.Repo.Repository.ID,
		Name:   f.Name,
		Color:  f.Color,
	}
	if err := database.NewLabels(label); err != nil {
		c.Error(err, "new labels")
		return
	}

	c.JSON(http.StatusCreated, toIssueLabel(label))
}

// UpdateLabel updates an existing label.
func UpdateRepoLabel(c *context.APIContext) {
	label, err := database.GetLabelOfRepoByID(c.Repo.Repository.ID, c.ParamsInt64("id"))
	if err != nil {
		c.NotFoundOrError(err, "get label of repository by ID")
		return
	}

	var f LabelOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	label.Name = f.Name
	label.Color = f.Color
	if err := database.UpdateLabel(label); err != nil {
		c.Error(err, "update label")
		return
	}

	c.JSON(http.StatusOK, toIssueLabel(label))
}

// DeleteLabel deletes a label from the repository.
func DeleteRepoLabel(c *context.APIContext) {
	if err := database.DeleteLabel(c.Repo.Repository.ID, c.ParamsInt64("id")); err != nil {
		c.Error(err, "delete label")
		return
	}

	c.NoContent()
}

// Issue labels

// ListIssueLabels lists all labels of an issue.
func ListIssueLabels(c *context.APIContext) {
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64("index"))
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	if err = issue.GetLabels(); err != nil {
		c.Error(err, "get labels")
		return
	}

	apiLabels := make([]*types.IssueLabel, len(issue.Labels))
	for i := range issue.Labels {
		apiLabels[i] = toIssueLabel(issue.Labels[i])
	}
	c.JSON(http.StatusOK, apiLabels)
}

// AddIssueLabels adds labels to an issue.
func AddIssueLabels(c *context.APIContext) {
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64("index"))
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	var f struct {
		Labels []int64 `json:"labels"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err = issue.AddLabels(c.User, f.Labels); err != nil {
		c.Error(err, "add labels")
		return
	}

	c.NoContent()
}

// ReplaceIssueLabels replaces all labels of an issue.
func ReplaceIssueLabels(c *context.APIContext) {
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64("index"))
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	var f struct {
		Labels []int64 `json:"labels"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err = issue.ReplaceLabels(f.Labels); err != nil {
		c.Error(err, "replace labels")
		return
	}

	c.NoContent()
}

// ClearIssueLabels removes all labels from an issue.
func ClearIssueLabels(c *context.APIContext) {
	issue, err := database.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64("index"))
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}

	if err = issue.ClearLabels(); err != nil {
		c.Error(err, "clear labels")
		return
	}

	c.NoContent()
}
