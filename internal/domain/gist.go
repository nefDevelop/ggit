package domain

import (
	"time"
)

// Gist represents a snippet of code or text.
type Gist struct {
	ID          int64
	UUID        string
	UserID      int64
	Description string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Files associated with this Gist.
	Files []*GistFile
}

// GistFile represents a single file within a Gist.
type GistFile struct {
	ID        int64
	GistID    int64
	Filename  string
	Content   string
	Language  string
	CreatedAt time.Time
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
