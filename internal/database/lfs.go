package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
	"github.com/nef734/ggit/internal/fsx"
	"github.com/nef734/ggit/internal/lfsx"
	"gorm.io/gorm"
)

// LFSObject represents an LFS object.
type LFSObject struct {
	ID          int64        `gorm:"primaryKey"`
	OID         lfsx.OID     `gorm:"column:oid;uniqueIndex:lfs_object_oid_repo_unique;index;not null"`
	Size        int64        `gorm:"not null"`
	RepoID      int64        `gorm:"uniqueIndex:lfs_object_oid_repo_unique;not null"`
	Storage     lfsx.Storage `gorm:"not null"`
	Created     time.Time    `json:"-" gorm:"-"`
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

func (o *LFSObject) ToDomain() *domain.LFSObject {
	if o == nil {
		return nil
	}
	return &domain.LFSObject{
		ID:      o.ID,
		OID:     o.OID,
		Size:    o.Size,
		RepoID:  o.RepoID,
		Storage: o.Storage,
		Created: o.Created,
	}
}

func FromDomainLFSObject(o *domain.LFSObject) *LFSObject {
	if o == nil {
		return nil
	}
	return &LFSObject{
		ID:          o.ID,
		OID:         o.OID,
		Size:        o.Size,
		RepoID:      o.RepoID,
		Storage:     o.Storage,
		CreatedUnix: o.Created.Unix(),
	}
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

func (s *LFSStore) CreateObjectAtomic(ctx context.Context, repoID int64, oid lfsx.OID, size int64, storage lfsx.Storage, tempPath, finalPath string) error {
	fsTx := fsx.NewTransaction()
	fsTx.StageCreate(tempPath, finalPath)
	defer fsTx.Rollback()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		object := &LFSObject{
			RepoID:  repoID,
			OID:     oid,
			Size:    size,
			Storage: storage,
		}
		return tx.Create(object).Error
	})
	if err != nil {
		return err
	}

	return fsTx.Commit()
}

type ErrLFSObjectNotExist struct {
	args errx.Args
}

func IsErrLFSObjectNotExist(err error) bool {
	var target ErrLFSObjectNotExist
	return errors.As(err, &target)
}

func (err ErrLFSObjectNotExist) Error() string {
	return fmt.Sprintf("LFS object does not exist: %v", err.args)
}

func (ErrLFSObjectNotExist) NotFound() bool {
	return true
}

// GetObjectByOID returns the LFS object with given OID. It returns
// ErrLFSObjectNotExist when not found.
func (s *LFSStore) GetObjectByOID(ctx context.Context, repoID int64, oid lfsx.OID) (*domain.LFSObject, error) {
	object := new(LFSObject)
	err := s.db.WithContext(ctx).Where("repo_id = ? AND oid = ?", repoID, string(oid)).First(object).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLFSObjectNotExist{args: errx.Args{"repoID": repoID, "oid": oid}}
		}
		return nil, err
	}
	return object.ToDomain(), err
}

// GetObjectsByOIDs returns LFS objects found within "oids". The returned list
// could have fewer elements if some oids were not found.
func (s *LFSStore) GetObjectsByOIDs(ctx context.Context, repoID int64, oids ...lfsx.OID) ([]*domain.LFSObject, error) {
	if len(oids) == 0 {
		return []*domain.LFSObject{}, nil
	}

	strOIDs := make([]string, len(oids))
	for i, oid := range oids {
		strOIDs[i] = string(oid)
	}

	objects := make([]*LFSObject, 0, len(oids))
	err := s.db.WithContext(ctx).Where("repo_id = ? AND oid IN (?)", repoID, strOIDs).Find(&objects).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	res := make([]*domain.LFSObject, len(objects))
	for i, object := range objects {
		res[i] = object.ToDomain()
	}
	return res, nil
}
