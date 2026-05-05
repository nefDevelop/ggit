package domain

import (
	"context"
	"time"

	"github.com/nef734/ggit/internal/lfsx"
)

// LFSObject represents an LFS object.
type LFSObject struct {
	ID      int64
	OID     lfsx.OID
	Size    int64
	RepoID  int64
	Storage lfsx.Storage
	Created time.Time
}

// LFSStore defines the interface for managing LFS objects.
type LFSStore interface {
	// CreateObject creates an LFS object record.
	CreateObject(ctx context.Context, repoID int64, oid lfsx.OID, size int64, storage lfsx.Storage) error
	// CreateObjectAtomic creates an LFS object record and moves the file atomically.
	CreateObjectAtomic(ctx context.Context, repoID int64, oid lfsx.OID, size int64, storage lfsx.Storage, tempPath, finalPath string) error
	// GetObjectByOID returns the LFS object with given OID.
	GetObjectByOID(ctx context.Context, repoID int64, oid lfsx.OID) (*LFSObject, error)
	// GetObjectsByOIDs returns LFS objects found within "oids".
	GetObjectsByOIDs(ctx context.Context, repoID int64, oids ...lfsx.OID) ([]*LFSObject, error)
}
