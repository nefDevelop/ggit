package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
)

type CreateGistForm struct {
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Files       []struct {
		Filename string `json:"filename" binding:"Required"`
		Content  string `json:"content" binding:"Required"`
		Language string `json:"language"`
	} `json:"files" binding:"Required"`
}

func CreateGist(c *context.APIContext) {
	var form CreateGistForm
	if err := c.BindJSON(&form); err != nil {
		c.StatusErrorf(http.StatusBadRequest, "BindJSON: %v", err)
		return
	}

	gist := &domain.Gist{
		UserID:      c.User.ID,
		Description: form.Description,
		IsPublic:    form.IsPublic,
	}
	for _, f := range form.Files {
		gist.Files = append(gist.Files, &domain.GistFile{
			Filename: f.Filename,
			Content:  f.Content,
			Language: f.Language,
		})
	}

	if err := database.Handle.Gists().Create(gist); err != nil {
		c.Error(err, "Create gist")
		return
	}

	c.JSON(http.StatusCreated, gist)
}

func GetGist(c *context.APIContext) {
	uuid := c.Params(":uuid")
	gist, err := database.Handle.Gists().GetByUUID(uuid)
	if err != nil {
		if domain.IsErrGistNotFound(err) {
			c.Status(http.StatusNotFound)
		} else {
			c.Error(err, "Get gist by UUID")
		}
		return
	}

	if !gist.IsPublic && (c.User == nil || c.User.ID != gist.UserID) {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gist)
}

func ListUserGists(c *context.APIContext) {
	username := c.Params(":username")
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		if database.IsErrUserNotExist(err) {
			c.Status(http.StatusNotFound)
		} else {
			c.Error(err, "Get user by username")
		}
		return
	}

	isOwner := c.User != nil && c.User.ID == user.ID
	gists, err := database.Handle.Gists().ListByUserID(user.ID, isOwner)
	if err != nil {
		c.Error(err, "List user gists")
		return
	}

	c.JSON(http.StatusOK, gists)
}

func ListPublicGists(c *context.APIContext) {
	gists, err := database.Handle.Gists().ListPublic()
	if err != nil {
		c.Error(err, "List public gists")
		return
	}

	c.JSON(http.StatusOK, gists)
}

func DeleteGist(c *context.APIContext) {
	uuid := c.Params(":uuid")
	gist, err := database.Handle.Gists().GetByUUID(uuid)
	if err != nil {
		if domain.IsErrGistNotFound(err) {
			c.Status(http.StatusNotFound)
		} else {
			c.Error(err, "Get gist for deletion")
		}
		return
	}

	if c.User.ID != gist.UserID && !c.User.IsAdmin {
		c.Status(http.StatusForbidden)
		return
	}

	if err := database.Handle.Gists().Delete(uuid); err != nil {
		c.Error(err, "Delete gist")
		return
	}

	c.Status(http.StatusNoContent)
}
