package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"

	"gogs.io/gogs/internal/errx"
	"gogs.io/gogs/internal/lfsx"
)

// LFSObject represents an LFS object.
type LFSObject struct {
	ID    int64  `gorm:"primaryKey"`
	OID   lfsx.OID `xorm:"UNIQUE(s) INDEX" gorm:"uniqueIndex:lfs_object_oid_repo_unique;index;not null"`
	Size  int64  `xorm:"NOT NULL" gorm:"not null"`
	RepoID      int64        `xorm:"UNIQUE(s)" gorm:"uniqueIndex:lfs_object_oid_repo_unique;not null"`
	Storage     lfsx.Storage `xorm:"NOT NULL" gorm:"not null"`
	Created     time.Time    `xorm:"-" json:"-" gorm:"-"`
	CreatedUnix int64
}

// BeforeCreate implements the GORM create hook.
func (o *LFSObject) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedUnix == 0 {
		o.CreatedUnix = tx.NowFunc().Unix()
	}
	return nil
}
// AfterFind implements the GORM query hook.
func (o *LFSObject) AfterFind(_ *gorm.DB) error {
	o.Created = time.Unix(o.CreatedUnix, 0).Local()
	return nil
}

// LFSStore is the storage layer for LFS objects.
type LFSStore struct {
	db *gorm.DB
}

func newLFSStore(db *gorm.DB) *LFSStore {
	return &LFSStore{db: db}
}

// CreateObject creates an LFS object record in database.
func (s *LFSStore) CreateObject(ctx context.Context, repoID int64, oid lfsx.OID, size int64, storage lfsx.Storage) error {
	object := &LFSObject{
		RepoID:  repoID,
		OID:     oid,
		Size:    size,
		Storage: storage,
	}
	return s.db.WithContext(ctx).Create(object).Error
}

type ErrLFSObjectNotExist struct {
	args errx.Args
}

func IsErrLFSObjectNotExist(err error) bool {
	return errors.As(err, &ErrLFSObjectNotExist{})
}

func (err ErrLFSObjectNotExist) Error() string {
	return fmt.Sprintf("LFS object does not exist: %v", err.args)
}

func (ErrLFSObjectNotExist) NotFound() bool {
	return true
}

// GetObjectByOID returns the LFS object with given OID. It returns
// ErrLFSObjectNotExist when not found.
func (s *LFSStore) GetObjectByOID(ctx context.Context, repoID int64, oid lfsx.OID) (*LFSObject, error) {
	object := new(LFSObject)
	err := s.db.WithContext(ctx).Where("repo_id = ? AND oid = ?", repoID, oid).First(object).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLFSObjectNotExist{args: errx.Args{"repoID": repoID, "oid": oid}}
		}
		return nil, err
	}
	return object, err
}

// GetObjectsByOIDs returns LFS objects found within "oids". The returned list
// could have fewer elements if some oids were not found.
func (s *LFSStore) GetObjectsByOIDs(ctx context.Context, repoID int64, oids ...lfsx.OID) ([]*LFSObject, error) {
	if len(oids) == 0 {
		return []*LFSObject{}, nil
	}

	objects := make([]*LFSObject, 0, len(oids))
	err := s.db.WithContext(ctx).Where("repo_id = ? AND oid IN (?)", repoID, oids).Find(&objects).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return objects, nil
}
