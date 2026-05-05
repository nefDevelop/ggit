package domain

import (
	"context"
	"time"
)

// KeyType indicates the type of the SSH key.
type KeyType int

const (
	KeyTypeUser KeyType = iota + 1
	KeyTypeDeploy
)

// PublicKey represents a user or deploy SSH public key.
type PublicKey struct {
	ID                int64
	OwnerID           int64
	Name              string
	Fingerprint       string
	Content           string
	Mode              AccessMode
	Type              KeyType
	Created           time.Time
	Updated           time.Time
	HasRecentActivity bool
	HasUsed           bool
}

// DeployKey represents deploy key information and its relation with repository.
type DeployKey struct {
	ID                int64
	KeyID             int64
	RepoID            int64
	Name              string
	Fingerprint       string
	Content           string
	Mode              AccessMode
	Created           time.Time
	Updated           time.Time
	HasRecentActivity bool
	HasUsed           bool
}

// PublicKeysStore defines the interface for managing public and deploy keys.
type PublicKeysStore interface {
	// AddPublicKey adds new public key to database and authorized_keys file.
	AddPublicKey(ctx context.Context, ownerID int64, name, content string) (*PublicKey, error)
	// GetPublicKeyByID returns public key by given ID.
	GetPublicKeyByID(ctx context.Context, id int64) (*PublicKey, error)
	// SearchPublicKeyByContent searches a public key using the content as prefix.
	SearchPublicKeyByContent(ctx context.Context, content string) (*PublicKey, error)
	// ListPublicKeys returns a list of public keys belongs to given user.
	ListPublicKeys(ctx context.Context, userID int64) ([]*PublicKey, error)
	// UpdatePublicKey updates given public key.
	UpdatePublicKey(ctx context.Context, key *PublicKey) error
	// DeletePublicKey deletes SSH key information both in database and authorized_keys file.
	DeletePublicKey(ctx context.Context, doer *User, id int64) error
	// RewriteAuthorizedKeys rewrites the "authorized_keys" file.
	RewriteAuthorizedKeys() error

	// AddDeployKey add new deploy key to database and authorized_keys file.
	AddDeployKey(ctx context.Context, repoID int64, name, content string, mode AccessMode) (*DeployKey, error)
	// GetDeployKeyByID returns deploy key by given ID.
	GetDeployKeyByID(ctx context.Context, id int64) (*DeployKey, error)
	// GetDeployKeyByRepo returns deploy key by given public key ID and repository ID.
	GetDeployKeyByRepo(ctx context.Context, keyID, repoID int64) (*DeployKey, error)
	// UpdateDeployKey updates deploy key information.
	UpdateDeployKey(ctx context.Context, key *DeployKey) error
	// DeleteDeployKey deletes deploy key from its repository authorized_keys file if needed.
	DeleteDeployKey(ctx context.Context, doer *User, id int64) error
	// ListDeployKeys returns all deploy keys by given repository ID.
	ListDeployKeys(ctx context.Context, repoID int64) ([]*DeployKey, error)
}
