package v2

import (
	"time"

	"gogs.io/gogs/internal/database"
)

// UserResponse represents a user in the API v2.
type UserResponse struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	AvatarURL       string `json:"avatar_url"`
	IsAdmin         bool   `json:"is_admin"`
	DisableNewsFeed bool   `json:"disable_news_feed"`
	CreatedAt       int64  `json:"created_at"`
}

func ToUserResponse(u *database.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
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

// RepositoryResponse represents a repository in the API v2.
type RepositoryResponse struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       *UserResponse `json:"owner"`
	NumStars    int           `json:"num_stars"`
	NumForks    int           `json:"num_forks"`
	NumWatches  int           `json:"num_watches"`
	IsPrivate   bool          `json:"is_private"`
	CreatedAt   int64         `json:"created_at"`
}

func ToRepositoryResponse(r *database.Repository) *RepositoryResponse {
	if r == nil {
		return nil
	}
	return &RepositoryResponse{
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

// OrganizationResponse represents an organization in the API v2.
type OrganizationResponse struct {
	ID          int64         `json:"id"`
	Username    string        `json:"username"`
	FullName    string        `json:"full_name"`
	Description string        `json:"description"`
	Owner       *UserResponse `json:"owner"`
	NumTeams    int           `json:"num_teams"`
	NumMembers  int           `json:"num_members"`
}

func ToOrganizationResponse(o *database.User, owner *database.User) *OrganizationResponse {
	if o == nil {
		return nil
	}
	return &OrganizationResponse{
		ID:          o.ID,
		Username:    o.Name,
		FullName:    o.FullName,
		Description: o.Description,
		Owner:       ToUserResponse(owner),
		NumTeams:    o.NumTeams,
		NumMembers:  o.NumMembers,
	}
}

// AuthSourceResponse represents a login source in the API v2.
type AuthSourceResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	TypeName    string    `json:"type_name"`
	IsActivated bool      `json:"is_activated"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToAuthSourceResponse(s *database.LoginSource) *AuthSourceResponse {
	if s == nil {
		return nil
	}
	return &AuthSourceResponse{
		ID:          s.ID,
		Name:        s.Name,
		TypeName:    s.TypeName(),
		IsActivated: s.IsActived,
		CreatedAt:   s.Created,
	}
}

// NoticeResponse represents a system notice in the API v2.
type NoticeResponse struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

func ToNoticeResponse(n *database.Notice) *NoticeResponse {
	if n == nil {
		return nil
	}
	return &NoticeResponse{
		ID:          n.ID,
		Type:        n.Type.TrStr(),
		Description: n.Description,
		CreatedAt:   n.CreatedUnix,
	}
}
