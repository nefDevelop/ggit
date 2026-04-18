package v2

import (


	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/route/api/v2/types"
)

func ToUserResponse(u *database.User) *types.UserResponse {
	if u == nil {
		return nil
	}
	return &types.UserResponse{
		ID:              u.ID,
		Username:        u.Name,
		FullName:        u.FullName,
		Email:           u.Email,
		AvatarURL:       u.AvatarURL(),
		IsAdmin:         u.IsAdmin,
		DisableNewsFeed: u.DisableNewsFeed,
		CreatedAt:       u.CreatedUnix,
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
		CreatedAt:   r.CreatedUnix,
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

func ToAuthSourceResponse(s *database.LoginSource) *types.AuthSourceResponse {
	if s == nil {
		return nil
	}
	return &types.AuthSourceResponse{
		ID:          s.ID,
		Name:        s.Name,
		TypeName:    s.TypeName(),
		IsActivated: s.IsActived,
		CreatedAt:   s.Created,
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
		CreatedAt:   n.CreatedUnix,
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
		ReadOnly:  key.Mode <= database.AccessModeRead,
	}
}
