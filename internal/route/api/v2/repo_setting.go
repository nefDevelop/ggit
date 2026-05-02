package v2

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/form"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// GetRepoSettings returns the repository settings.
func GetRepoSettings(c *context.APIContext) {
	repo := c.Repo.Repository
	res := map[string]any{
		"id":                      repo.ID,
		"name":                    repo.Name,
		"description":             repo.Description,
		"website":                 repo.Website,
		"private":                 repo.IsPrivate,
		"unlisted":                repo.IsUnlisted,
		"archived":                repo.IsArchived,
		"enable_wiki":             repo.EnableWiki,
		"allow_public_wiki":       repo.AllowPublicWiki,
		"enable_external_wiki":    repo.EnableExternalWiki,
		"external_wiki_url":       repo.ExternalWikiURL,
		"wiki_folder_name":        repo.WikiFolderName,
		"enable_issues":           repo.EnableIssues,
		"allow_public_issues":     repo.AllowPublicIssues,
		"enable_external_tracker": repo.EnableExternalTracker,
		"external_tracker_url":    repo.ExternalTrackerURL,
		"external_tracker_format": repo.ExternalTrackerFormat,
		"external_tracker_style":  repo.ExternalTrackerStyle,
		"enable_pulls":            repo.EnablePulls,
		"pulls_ignore_whitespace": repo.PullsIgnoreWhitespace,
		"pulls_allow_rebase":      repo.PullsAllowRebase,
		"pulls_allow_manual_merge": repo.PullsAllowManualMerge,
		"default_branch":          repo.DefaultBranch,
		"use_custom_avatar":       repo.UseCustomAvatar,
		"avatar_url":              repo.AvatarLink(),
	}

	c.JSON(http.StatusOK, res)
}

// UpdateRepoSettings updates the repository settings.
func UpdateRepoSettings(c *context.APIContext) {
	var f form.RepoSetting
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	repo := c.Repo.Repository
	visibilityChanged := repo.IsPrivate != f.Private || repo.IsUnlisted != f.Unlisted
	repo.Name = f.RepoName
	repo.LowerName = strings.ToLower(f.RepoName)
	repo.Description = f.Description
	repo.Website = f.Website
	repo.IsPrivate = f.Private
	repo.IsUnlisted = f.Unlisted
	repo.IsArchived = f.Archived
	repo.WikiFolderName = f.WikiFolderName

	// Advanced settings
	repo.EnableWiki = f.EnableWiki
	repo.AllowPublicWiki = f.AllowPublicWiki
	repo.EnableExternalWiki = f.EnableExternalWiki
	repo.ExternalWikiURL = f.ExternalWikiURL
	repo.EnableIssues = f.EnableIssues
	repo.AllowPublicIssues = f.AllowPublicIssues
	repo.EnableExternalTracker = f.EnableExternalTracker
	repo.ExternalTrackerURL = f.ExternalTrackerURL
	repo.ExternalTrackerFormat = f.TrackerURLFormat
	repo.ExternalTrackerStyle = f.TrackerIssueStyle
	repo.EnablePulls = f.EnablePulls
	repo.PullsIgnoreWhitespace = f.PullsIgnoreWhitespace
	repo.PullsAllowRebase = f.PullsAllowRebase
	repo.PullsAllowManualMerge = f.PullsAllowManualMerge
	if f.Branch != "" {
		repo.DefaultBranch = f.Branch
	} else if f.DefaultBranch != "" {
		repo.DefaultBranch = f.DefaultBranch
	}

	if err := database.UpdateRepository(repo, visibilityChanged); err != nil {
		c.Error(err, "update repository")
		return
	}

	c.JSON(http.StatusOK, repo)
}

// GetRepoCollaborators returns a list of collaborators for the repository.
func GetRepoCollaborators(c *context.APIContext) {
	users, err := c.Repo.Repository.GetCollaborators()
	if err != nil {
		c.Error(err, "get collaborators")
		return
	}
	c.JSON(http.StatusOK, users)
}

// AddRepoCollaborator adds a new collaborator to the repository.
func AddRepoCollaborator(c *context.APIContext) {
	userName := c.Params("collaborator")
	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	if err = c.Repo.Repository.AddCollaborator(u); err != nil {
		c.Error(err, "add collaborator")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// UpdateRepoCollaborator updates the collaborator's access mode.
func UpdateRepoCollaborator(c *context.APIContext) {
	userName := c.Params("collaborator")
	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	var f struct {
		Mode int `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err = c.Repo.Repository.ChangeCollaborationAccessMode(u.ID, database.AccessMode(f.Mode)); err != nil {
		c.Error(err, "change collaboration access mode")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// DeleteRepoCollaborator deletes a collaborator from the repository.
func DeleteRepoCollaborator(c *context.APIContext) {
	userName := c.Params("collaborator")
	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	if err = c.Repo.Repository.DeleteCollaboration(u.ID); err != nil {
		c.Error(err, "delete collaboration")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetRepoHooks returns a list of webhooks for the repository.
func GetRepoHooks(c *context.APIContext) {
	ws, err := database.GetWebhooksByRepoID(c.Repo.Repository.ID)
	if err != nil {
		c.Error(err, "get webhooks")
		return
	}
	
	resp := make([]*apiv2types.WebhookResponse, len(ws))
	for i, w := range ws {
		resp[i] = ToWebhookResponse(w)
	}
	c.JSON(http.StatusOK, resp)
}

func parseHookEvents(events []string, w *database.Webhook) {
	w.HookEvent = &database.HookEvent{
		ChooseEvents: true,
	}
	for _, e := range events {
		switch e {
		case "push":
			w.Push = true
		case "create":
			w.Create = true
		case "delete":
			w.Delete = true
		case "fork":
			w.Fork = true
		case "issues":
			w.Issues = true
		case "issue_comment":
			w.IssueComment = true
		case "pull_request":
			w.PullRequest = true
		case "release":
			w.Release = true
		}
	}
}

// CreateRepoHook creates a new webhook for the repository.
func CreateRepoHook(c *context.APIContext) {
	var f apiv2types.CreateWebhookRequest
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	w := &database.Webhook{
		RepoID:      c.Repo.Repository.ID,
		URL:         f.URL,
		ContentType: database.ToHookContentType(f.ContentType),
		Secret:      f.Secret,
		IsActive:    f.Active,
	}
	parseHookEvents(f.Events, w)

	if err := database.CreateWebhook(w); err != nil {
		c.Error(err, "create webhook")
		return
	}

	c.JSON(http.StatusCreated, ToWebhookResponse(w))
}

// UpdateRepoHook updates a webhook for the repository.
func UpdateRepoHook(c *context.APIContext) {
	id := c.ParamsInt64("id")
	var f apiv2types.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	w, err := database.GetWebhookOfRepoByID(c.Repo.Repository.ID, id)
	if err != nil {
		c.Error(err, "get webhook")
		return
	}

	if f.URL != "" {
		w.URL = f.URL
	}
	if f.ContentType != "" {
		w.ContentType = database.ToHookContentType(f.ContentType)
	}
	if f.Secret != "" {
		w.Secret = f.Secret
	}
	if f.Active != nil {
		w.IsActive = *f.Active
	}
	if len(f.Events) > 0 {
		parseHookEvents(f.Events, w)
	}

	if err = database.UpdateWebhook(w); err != nil {
		c.Error(err, "update webhook")
		return
	}

	c.JSON(http.StatusOK, ToWebhookResponse(w))
}

// TestRepoHook triggers a test delivery for a webhook.
func TestRepoHook(c *context.APIContext) {
	id := c.ParamsInt64("id")
	
	// Create a dummy push payload for testing
	p := &apiv2types.WebhookPushPayload{
		Ref: "refs/heads/master",
		Before: "0000000000000000000000000000000000000000",
		After: "0000000000000000000000000000000000000000",
		Commits: []*apiv2types.WebhookPayloadCommit{},
		Pusher: ToUserResponse(c.User),
		Sender: ToUserResponse(c.User),
		Repo: ToRepositoryResponse(c.Repo.Repository),
	}

	if err := database.TestWebhook(c.Repo.Repository, database.HookEventTypePush, p, id); err != nil {
		c.Error(err, "test webhook")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// GetRepoHookHistory returns the delivery history for a specific webhook.
func GetRepoHookHistory(c *context.APIContext) {
	id := c.ParamsInt64("id")
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	tasks, err := database.HookTasks(id, page)
	if err != nil {
		c.Error(err, "get hook tasks")
		return
	}

	resp := make([]*apiv2types.WebhookHistoryResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = ToWebhookHistoryResponse(t)
	}
	c.JSON(http.StatusOK, resp)
}

// RedeliverRepoHookTask redelivers a specific webhook task.
func RedeliverRepoHookTask(c *context.APIContext) {
	id := c.ParamsInt64("id")
	uuid := c.Params("uuid")

	task, err := database.GetHookTaskOfWebhookByUUID(id, uuid)
	if err != nil {
		c.Error(err, "get hook task")
		return
	}

	if err := database.RedeliverHookTask(task); err != nil {
		c.Error(err, "redeliver hook task")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// DeleteRepoHook deletes a webhook from the repository.
func DeleteRepoHook(c *context.APIContext) {
	id := c.ParamsInt64("id")
	if err := database.DeleteWebhookOfRepoByID(c.Repo.Repository.ID, id); err != nil {
		c.Error(err, "delete webhook")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

// DeleteRepository deletes the repository.
func DeleteRepository(c *context.APIContext) {
	var f struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if f.Name != c.Repo.Repository.Name {
		c.StatusError(http.StatusUnprocessableEntity, "Confirmation name does not match.", nil)
		return
	}

	if err := database.DeleteRepository(c.User.ID, c.Repo.Repository.ID); err != nil {
		c.Error(err, "delete repository")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// ListRepoRemotes returns a list of remotes for the repository.
func ListRepoRemotes(c *context.APIContext) {
	remotes, err := c.Repo.GitRepo.Remotes()
	if err != nil {
		c.Error(err, "get remotes")
		return
	}
	c.JSON(http.StatusOK, remotes)
}

// AddRepoRemote adds a new remote to the repository.
func AddRepoRemote(c *context.APIContext) {
	var f struct {
		Name string `json:"name" binding:"required"`
		URL  string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := c.Repo.GitRepo.RemoteAdd(f.Name, f.URL); err != nil {
		c.Error(err, "add remote")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteRepoRemote deletes a remote from the repository.
func DeleteRepoRemote(c *context.APIContext) {
	name := c.Params("name")
	if err := c.Repo.GitRepo.RemoteRemove(name); err != nil {
		c.Error(err, "delete remote")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// TransferRepoOwnership transfers the repository ownership to another user/org.
func TransferRepoOwnership(c *context.APIContext) {
	var f struct {
		NewOwner string `json:"new_owner" binding:"required"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if err := database.TransferOwnership(c.User, f.NewOwner, c.Repo.Repository); err != nil {
		c.Error(err, "transfer ownership")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// UpdateRepoAvatar updates the repository avatar.
func UpdateRepoAvatar(c *context.APIContext) {
	var f struct {
		Avatar *multipart.FileHeader `form:"avatar"`
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
		file, header, err := c.Req.FormFile("avatar")
		if err == nil {
			f.Avatar = header
			_ = file.Close()
		}
	}

	if f.Avatar != nil {
		content, err := f.Avatar.Open()
		if err != nil {
			c.Error(err, "open avatar file")
			return
		}
		defer func() { _ = content.Close() }()

		data, err := io.ReadAll(content)
		if err != nil {
			c.Error(err, "read avatar file")
			return
		}

		if err := c.Repo.Repository.UploadAvatar(data); err != nil {
			c.Error(err, "upload avatar")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":     "ok",
		"avatar_url": c.Repo.Repository.AvatarLink(),
	})
}

// DeleteRepoAvatar deletes the repository avatar.
func DeleteRepoAvatar(c *context.APIContext) {
	if err := c.Repo.Repository.DeleteAvatar(); err != nil {
		c.Error(err, "delete avatar")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}
