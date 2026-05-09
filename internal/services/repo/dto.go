package repo

import (
	"time"
)

type CreateOptions struct {
	Name          string
	Description   string
	Gitignores    string
	License       string
	Readme        string
	Website       string
	DefaultBranch string
	IsPrivate     bool
	IsUnlisted    bool
	IsArchived    bool
	AutoInit      bool
	CreatedAt     time.Time
}
