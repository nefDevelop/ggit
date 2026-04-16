package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/form"
)

// ListReleases returns a list of releases for the repository.
func ListReleases(c *context.APIContext) {
	releases, err := database.GetReleasesByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get releases")
		return
	}

	for _, r := range releases {
		_ = r.LoadAttributes()
	}

	c.JSON(http.StatusOK, releases)
}

// GetRelease returns a single release by its ID.
func GetRelease(c *context.APIContext) {
	id := c.ParamsInt64("id")
	rel, err := database.GetReleaseByID(id)
	if err != nil {
		c.NotFoundOrError(err, "get release by ID")
		return
	}

	if rel.RepoID != c.Repo.Repository.ID {
		c.NotFound()
		return
	}

	_ = rel.LoadAttributes()
	c.JSON(http.StatusOK, rel)
}

// CreateRelease creates a new release.
func CreateRelease(c *context.APIContext) {
	var f form.NewRelease
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	// Basic validation
	if !c.Repo.GitRepo.HasBranch(f.Target) {
		c.StatusError(http.StatusUnprocessableEntity, "Target branch does not exist.", nil)
		return
	}

	rel := &database.Release{
		RepoID:       c.Repo.Repository.ID,
		PublisherID:  c.User.ID,
		Title:        f.Title,
		TagName:      f.TagName,
		Target:       f.Target,
		Note:         f.Content,
		IsDraft:      len(f.Draft) > 0,
		IsPrerelease: f.Prerelease,
	}

	if err := database.NewRelease(c.Repo.GitRepo, rel, f.Files); err != nil {
		c.Error(err, "new release")
		return
	}

	c.JSON(http.StatusCreated, rel)
}

// UpdateRelease updates an existing release.
func UpdateRelease(c *context.APIContext) {
	id := c.ParamsInt64("id")
	rel, err := database.GetReleaseByID(id)
	if err != nil {
		c.NotFoundOrError(err, "get release by ID")
		return
	}

	if rel.RepoID != c.Repo.Repository.ID {
		c.NotFound()
		return
	}

	var f form.EditRelease
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	isPublish := rel.IsDraft && f.Draft == ""
	rel.Title = f.Title
	rel.Note = f.Content
	rel.IsDraft = len(f.Draft) > 0
	rel.IsPrerelease = f.Prerelease

	if err := database.UpdateRelease(c.User, c.Repo.GitRepo, rel, isPublish, f.Files); err != nil {
		c.Error(err, "update release")
		return
	}

	c.JSON(http.StatusOK, rel)
}

// DeleteRelease deletes a release.
func DeleteRelease(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if err := database.DeleteReleaseOfRepoByID(c.Repo.Repository.ID, id); err != nil {
		c.Error(err, "delete release")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}
