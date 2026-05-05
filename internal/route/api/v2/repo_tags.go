package v2

import (
	"net/http"
	"sort"
	"time"

	"github.com/nef734/ggit/internal/context"
)

type TagResponse struct {
	Name      string    `json:"name"`
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Tagger    string    `json:"tagger"`
	CreatedAt time.Time `json:"created_at"`
	ZipURL    string    `json:"zip_url"`
	TarURL    string    `json:"tar_url"`
}

// ListRepoTags returns all tags for a repository.
func ListRepoTags(c *context.APIContext) {
	tags, err := c.Repo.GitRepo.Tags()
	if err != nil {
		c.Error(err, "list tags")
		return
	}

	username := c.Params(":username")
	reponame := c.Params(":reponame")

	resp := make([]*TagResponse, 0, len(tags))
	for _, tag := range tags {
		resp = append(resp, &TagResponse{
			Name:   tag,
			ZipURL: "/api/v2/" + username + "/" + reponame + "/archive/" + tag + ".zip",
			TarURL: "/api/v2/" + username + "/" + reponame + "/archive/" + tag + ".tar.gz",
		})
	}

	// Sort alphabetically descending (newest semver-ish tags first)
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Name > resp[j].Name
	})

	c.JSON(http.StatusOK, resp)
}

type CreateTagRequest struct {
	Name    string `json:"name" binding:"required"`
	Target  string `json:"target" binding:"required"` // branch name or commit SHA
	Message string `json:"message"`
}

// CreateRepoTag creates a new tag.
func CreateRepoTag(c *context.APIContext) {
	var f CreateTagRequest
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := c.Repo.GitRepo.CreateTag(f.Name, f.Target); err != nil {
		c.Error(err, "create tag")
		return
	}

	username := c.Params(":username")
	reponame := c.Params(":reponame")

	c.JSON(http.StatusCreated, &TagResponse{
		Name:   f.Name,
		ZipURL: "/api/v2/" + username + "/" + reponame + "/archive/" + f.Name + ".zip",
		TarURL: "/api/v2/" + username + "/" + reponame + "/archive/" + f.Name + ".tar.gz",
	})
}

// DeleteRepoTag deletes a tag.
func DeleteRepoTag(c *context.APIContext) {
	tagName := c.Params("tag")
	if tagName == "" {
		c.StatusError(http.StatusBadRequest, "Tag name is required.", nil)
		return
	}

	if err := c.Repo.GitRepo.DeleteTag(tagName); err != nil {
		c.Error(err, "delete tag")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
