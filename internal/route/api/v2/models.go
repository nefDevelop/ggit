package v2

import (
	"context"
	"path"
	"time"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/route/api/v2/types"
)

func ToUserResponse(u *database.User) *types.UserResponse {
	if u == nil {
		return nil
	}
	return &types.UserResponse{
		ID:               u.ID,
		Username:         u.Name,
		FullName:         u.FullName,
		Email:            u.Email,
		AvatarURL:        u.AvatarURL(),
		IsAdmin:          u.IsAdmin,
		DisableNewsFeed:  u.DisableNewsFeed,
		ProhibitLogin:    u.ProhibitLogin,
		AllowGitHook:     u.AllowGitHook,
		AllowImportLocal: u.AllowImportLocal,
		MaxRepoCreation:  u.MaxRepoCreation,
		CreatedAt:        time.Unix(u.CreatedUnix, 0),
	}
}

func ToUserEmailResponse(e *domain.EmailAddress) *types.UserEmailResponse {
	if e == nil {
		return nil
	}
	return &types.UserEmailResponse{
		Email:    e.Email,
		Verified: e.IsActivated,
		Primary:  e.IsPrimary,
	}
}

func ToRepositoryResponse(r *database.Repository) *types.RepositoryResponse {
	if r == nil {
		return nil
	}
	return &types.RepositoryResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Owner:       ToUserResponse(r.Owner),
		NumStars:    r.NumStars,
		NumForks:    r.NumForks,
		NumWatches:  r.NumWatches,
		IsPrivate:   r.IsPrivate,
		CreatedAt:   time.Unix(r.CreatedUnix, 0),
		Topics:      r.Topics,
	}
}

func ToOrganizationResponse(o *database.User, owner *database.User) *types.OrganizationResponse {
	if o == nil {
		return nil
	}
	return &types.OrganizationResponse{
		ID:          o.ID,
		Username:    o.Name,
		FullName:    o.FullName,
		Description: o.Description,
		Owner:       ToUserResponse(owner),
		NumTeams:    o.NumTeams,
		NumMembers:  o.NumMembers,
	}
}

func ToAuthSourceResponse(s *domain.LoginSource) *types.AuthSourceResponse {
	if s == nil {
		return nil
	}
	return &types.AuthSourceResponse{
		ID:          s.ID,
		Name:        s.Name,
		TypeName:    auth.Name(auth.Type(s.Type)),
		IsActivated: s.IsActived,
		// CreatedAt:   s.Created, // domain.LoginSource doesn't have Created yet?
		RawConfig: s.Config,
	}
}

func ToWebhookResponse(w *database.Webhook) *types.WebhookResponse {
	if w == nil {
		return nil
	}
	return &types.WebhookResponse{
		ID:          w.ID,
		URL:         w.URL,
		ContentType: w.ContentType.Name(),
		Secret:      w.Secret,
		Events:      w.EventsArray(),
		IsSSL:       w.IsSSL,
		IsActive:    w.IsActive,
		LastStatus:  int(w.LastStatus),
		Updated:     w.UpdatedUnix,
	}
}

func ToWebhookHistoryResponse(t *database.HookTask) *types.WebhookHistoryResponse {
	if t == nil {
		return nil
	}

	reqHeaders := make(map[string]string)
	if t.RequestInfo != nil {
		reqHeaders = t.RequestInfo.Headers
	}

	respHeaders := make(map[string]string)
	respStatus := 0
	respBody := ""
	if t.ResponseInfo != nil {
		respHeaders = t.ResponseInfo.Headers
		respStatus = t.ResponseInfo.Status
		respBody = t.ResponseInfo.Body
	}

	return &types.WebhookHistoryResponse{
		ID:              t.ID,
		UUID:            t.UUID,
		EventType:       string(t.EventType),
		IsSucceed:       t.IsSucceed,
		Delivered:       time.Unix(0, t.Delivered),
		RequestHeaders:  reqHeaders,
		RequestBody:     t.PayloadContent,
		ResponseHeaders: respHeaders,
		ResponseBody:    respBody,
		ResponseStatus:  respStatus,
	}
}

func ToNoticeResponse(n *database.Notice) *types.NoticeResponse {
	if n == nil {
		return nil
	}
	return &types.NoticeResponse{
		ID:          n.ID,
		Type:        n.Type.TrStr(),
		Description: n.Description,
		CreatedAt:   time.Unix(n.CreatedUnix, 0),
	}
}

func ToDeployKeyResponse(key *database.DeployKey) *types.RepositoryDeployKeyResponse {
	if key == nil {
		return nil
	}
	return &types.RepositoryDeployKeyResponse{
		ID:        key.ID,
		Key:       key.Content,
		URL:       "", // Not used in frontend yet
		Title:     key.Name,
		CreatedAt: key.Created,
		ReadOnly:  key.Mode <= domain.AccessModeRead,
	}
}

func ToAuditLogResponse(a *database.Action) *types.AuditLogResponse {
	if a == nil {
		return nil
	}

	opName := "Unknown"
	switch database.ActionType(a.OpType) {
	case database.ActionCreateRepo:
		opName = "Create Repository"
	case database.ActionRenameRepo:
		opName = "Rename Repository"
	case database.ActionStarRepo:
		opName = "Star Repository"
	case database.ActionWatchRepo:
		opName = "Watch Repository"
	case database.ActionCommitRepo:
		opName = "Push Commits"
	case database.ActionCreateIssue:
		opName = "Create Issue"
	case database.ActionCreatePullRequest:
		opName = "Create Pull Request"
	case database.ActionTransferRepo:
		opName = "Transfer Repository"
	case database.ActionPushTag:
		opName = "Push Tag"
	case database.ActionCommentIssue:
		opName = "Comment Issue"
	case database.ActionMergePullRequest:
		opName = "Merge Pull Request"
	case database.ActionCloseIssue:
		opName = "Close Issue"
	case database.ActionReopenIssue:
		opName = "Reopen Issue"
	case database.ActionClosePullRequest:
		opName = "Close Pull Request"
	case database.ActionReopenPullRequest:
		opName = "Reopen Pull Request"
	case database.ActionCreateBranch:
		opName = "Create Branch"
	case database.ActionDeleteBranch:
		opName = "Delete Branch"
	case database.ActionDeleteTag:
		opName = "Delete Tag"
	case database.ActionForkRepo:
		opName = "Fork Repository"
	case database.ActionUpdateConfig:
		opName = "Update Configuration"
	}

	return &types.AuditLogResponse{
		ID:        a.ID,
		OpType:    int(a.OpType),
		OpName:    opName,
		ActorID:   a.ActUserID,
		ActorName: a.ActUserName,
		RepoID:    a.RepoID,
		RepoPath:  path.Join(a.RepoUserName, a.RepoName),
		Content:   a.Content,
		CreatedAt: time.Unix(a.CreatedUnix, 0),
	}
}

func ToCommitStatusResponse(s *database.CommitStatus) *types.CommitStatusResponse {
	if s == nil {
		return nil
	}

	u, _ := database.Handle.Users().GetByID(context.Background(), s.CreatorID)

	return &types.CommitStatusResponse{
		ID:          s.ID,
		State:       string(s.State),
		TargetURL:   s.TargetURL,
		Description: s.Description,
		Context:     s.Context,
		Creator:     ToUserResponse(database.FromDomainUser(u)),
		CreatedAt:   time.Unix(s.CreatedAt, 0),
		UpdatedAt:   time.Unix(s.UpdatedAt, 0),
	}
}

func ToCombinedCommitStatusResponse(repo *database.Repository, sha string, state database.CommitStatusState, statuses []*database.CommitStatus) *types.CombinedCommitStatusResponse {
	respStatuses := make([]*types.CommitStatusResponse, len(statuses))
	for i, s := range statuses {
		respStatuses[i] = ToCommitStatusResponse(s)
	}

	return &types.CombinedCommitStatusResponse{
		State:    string(state),
		Statuses: respStatuses,
		Repo:     ToRepositoryResponse(repo),
		SHA:      sha,
	}
}
