package v2

import (
	"net/http"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/form"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// ListReleases returns a list of releases for the repository.
func ListReleases(c *context.APIContext) {
	releases, err := database.GetReleasesByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get releases")
		return
	}

	resp := make([]*apiv2types.RepositoryReleaseResponse, len(releases))
	for i, r := range releases {
		_ = r.LoadAttributes()
		resp[i] = r.APIFormat()
	}

	c.JSON(http.StatusOK, resp)
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
	c.JSON(http.StatusOK, rel.APIFormat())
}

// CreateRelease creates a new release.
func CreateRelease(c *context.APIContext) {
	var f form.NewRelease
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if c.Repo.GitRepo == nil {
		c.StatusError(http.StatusInternalServerError, "Git repository is not initialized.", nil)
		return
	}

	if c.User == nil {
		c.StatusError(http.StatusUnauthorized, "Only authenticated user is allowed to call APIs.", nil)
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
		IsDraft:      f.Draft,
		IsPrerelease: f.Prerelease,
	}

	if f.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, f.CreatedAt); err == nil {
			rel.CreatedUnix = t.Unix()
		}
	}

	if err := database.NewRelease(c.Repo.GitRepo, rel, f.Files); err != nil {
		c.Error(err, "new release")
		return
	}

	// Load attributes (like Publisher) to avoid nil pointer in APIFormat
	if err := rel.LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	c.JSON(http.StatusCreated, rel.APIFormat())
}

// UploadReleaseAttachment uploads a release attachment.
func UploadReleaseAttachment(c *context.APIContext) {
	if !conf.Release.Attachment.Enabled {
		c.Status(http.StatusNotImplemented)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err, "get file")
		return
	}
	defer file.Close()

	// Size check
	if conf.Release.Attachment.MaxSize > 0 && header.Size > conf.Release.Attachment.MaxSize*1024*1024 {
		c.StatusError(http.StatusRequestEntityTooLarge, "File is too large.", nil)
		return
	}

	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	if n > 0 {
		buf = buf[:n]
	}

	attach, err := database.NewAttachment(header.Filename, buf, file)
	if err != nil {
		c.Error(err, "new attachment")
		return
	}

	c.JSON(http.StatusOK, attach)
}

// GetAttachment serves an attachment file.
func GetAttachment(c *context.APIContext) {
	uuid := c.Params(":uuid")
	attach, err := database.GetAttachmentByUUID(uuid)
	if err != nil {
		c.NotFoundOrError(err, "get attachment")
		return
	}

	c.ServeFile(attach.LocalPath(), attach.Name)
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

	if c.Repo.GitRepo == nil {
		c.StatusError(http.StatusInternalServerError, "Git repository is not initialized.", nil)
		return
	}

	if c.User == nil {
		c.StatusError(http.StatusUnauthorized, "Only authenticated user is allowed to call APIs.", nil)
		return
	}

	isPublish := rel.IsDraft && !f.Draft
	rel.Title = f.Title
	rel.Note = f.Content
	rel.IsDraft = f.Draft
	rel.IsPrerelease = f.Prerelease

	if err := database.UpdateRelease(c.User, c.Repo.GitRepo, rel, isPublish, f.Files); err != nil {
		c.Error(err, "update release")
		return
	}

	// Load attributes (like Publisher and Attachments) to ensure standard API response
	_ = rel.LoadAttributes()

	c.JSON(http.StatusOK, rel.APIFormat())
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
