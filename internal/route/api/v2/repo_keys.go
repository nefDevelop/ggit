package v2

import (
	"net/http"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
)

// ListRepoDeployKeys returns a list of deploy keys for the repository.
func ListRepoDeployKeys(c *context.APIContext) {
	keys, err := database.ListDeployKeys(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "list deploy keys")
		return
	}

	res := make([]any, len(keys))
	for i, k := range keys {
		if err := k.GetContent(); err != nil {
			c.Error(err, "get key content")
			return
		}
		res[i] = ToDeployKeyResponse(k)
	}

	c.JSON(http.StatusOK, res)
}

// AddRepoDeployKey adds a new deploy key to the repository.
func AddRepoDeployKey(c *context.APIContext) {
	var f struct {
		Title    string `json:"title" binding:"required"`
		Key      string `json:"key" binding:"required"`
		ReadOnly bool   `json:"read_only"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	content, err := database.CheckPublicKeyString(f.Key)
	if err != nil {
		c.ErrorStatus(http.StatusUnprocessableEntity, err)
		return
	}

	mode := database.AccessModeWrite
	if f.ReadOnly {
		mode = database.AccessModeRead
	}

	key, err := database.AddDeployKey(c.Repo.Repository.ID, f.Title, content, mode)
	if err != nil {
		if database.IsErrKeyAlreadyExist(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Key content already exists.", err)
		} else if database.IsErrDeployKeyAlreadyExist(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Deploy key already exists for this repository.", err)
		} else {
			c.Error(err, "add deploy key")
		}
		return
	}

	if err := key.GetContent(); err != nil {
		c.Error(err, "get key content")
		return
	}

	c.JSON(http.StatusCreated, ToDeployKeyResponse(key))
}

// DeleteRepoDeployKey deletes a deploy key from the repository.
func DeleteRepoDeployKey(c *context.APIContext) {
	id := c.ParamsInt64(":id")
	if err := database.DeleteDeployKey(c.User, id); err != nil {
		c.Error(err, "delete deploy key")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
