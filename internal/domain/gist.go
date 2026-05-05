package domain

import (
	"time"
)

// Gist represents a snippet of code or text.
type Gist struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	UserID      int64     `json:"user_id"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Files associated with this Gist.
	Files []*GistFile `json:"files,omitempty"`
}

// GistFile represents a single file within a Gist.
type GistFile struct {
	ID        int64     `json:"id"`
	GistID    int64     `json:"gist_id"`
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}

// GistStore is the persistence layer for Gists.
type GistStore interface {
	// Create saves a new Gist and its files.
	Create(gist *Gist) error
	// GetByUUID returns a Gist by its UUID.
	GetByUUID(uuid string) (*Gist, error)
	// ListByUserID returns public Gists for a user.
	ListByUserID(userID int64, isOwner bool) ([]*Gist, error)
	// ListPublic returns all public Gists.
	ListPublic() ([]*Gist, error)
	// Update updates Gist metadata and creates a new version of files.
	Update(gist *Gist) error
	// Delete removes a Gist and all its files.
	Delete(uuid string) error
}

// ErrGistNotFound represents an error when a Gist cannot be found.
type ErrGistNotFound struct {
	UUID string
}

func (e ErrGistNotFound) Error() string {
	return "gist not found [uuid: " + e.UUID + "]"
}

// IsErrGistNotFound returns true if the error is an ErrGistNotFound.
func IsErrGistNotFound(err error) bool {
	_, ok := err.(ErrGistNotFound)
	return ok
}
