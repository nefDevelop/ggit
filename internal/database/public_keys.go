package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/osx"
	"gorm.io/gorm"
)

// PublicKeysStore is the storage layer for public keys.
type PublicKeysStore struct {
	db *gorm.DB
}

func newPublicKeysStore(db *gorm.DB) *PublicKeysStore {
	return &PublicKeysStore{db: db}
}

func authorizedKeysPath() string {
	return filepath.Join(conf.SSH.RootPath, "authorized_keys")
}

// RewriteAuthorizedKeys rewrites the "authorized_keys" file under the SSH root
// path with all public keys stored in the database.
func (s *PublicKeysStore) RewriteAuthorizedKeys() error {
	sshOpLocker.Lock()
	defer sshOpLocker.Unlock()

	err := os.MkdirAll(conf.SSH.RootPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create SSH root path"+": %v", err)
	}
	fpath := authorizedKeysPath()
	tempPath := fpath + ".tmp"
	f, err := os.OpenFile(tempPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("create temporary file"+": %v", err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(tempPath)
	}()

	// NOTE: More recently updated keys are more likely to be used more frequently,
	// putting them in the earlier lines could speed up the key lookup by SSHD.
	rows, err := s.db.Model(&PublicKey{}).Order("updated_unix DESC").Rows()
	if err != nil {
		return fmt.Errorf("iterate public keys"+": %v", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var key PublicKey
		err = s.db.ScanRows(rows, &key)
		if err != nil {
			return fmt.Errorf("scan rows"+": %v", err)
		}

		_, err = f.WriteString(key.AuthorizedString())
		if err != nil {
			return fmt.Errorf("write key %d: %v", key.ID, err)
		}
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("check rows.Err"+": %v", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("close temporary file"+": %v", err)
	}
	if osx.Exist(fpath) {
		err = os.Remove(fpath)
		if err != nil {
			return fmt.Errorf("remove"+": %v", err)
		}
	}
	err = os.Rename(tempPath, fpath)
	if err != nil {
		return fmt.Errorf("rename"+": %v", err)
	}
	return nil
}

// AddPublicKey adds new public key to database and authorized_keys file.
func (s *PublicKeysStore) AddPublicKey(ctx context.Context, ownerID int64, name, content string) (*domain.PublicKey, error) {
	if err := checkKeyContent(content); err != nil {
		return nil, err
	}

	// Key name of same user cannot be duplicated.
	var count int64
	err := s.db.WithContext(ctx).Model(new(PublicKey)).Where("owner_id = ? AND name = ?", ownerID, name).Count(&count).Error
	if err != nil {
		return nil, err
	} else if count > 0 {
		return nil, ErrKeyNameAlreadyUsed{ownerID, name}
	}

	key := &PublicKey{
		OwnerID: ownerID,
		Name:    name,
		Content: content,
		Mode:    domain.AccessModeWrite,
		Type:    domain.KeyTypeUser,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return addKey(tx, key)
	})
	if err != nil {
		return nil, err
	}

	return key.ToDomain(), nil
}

// GetPublicKeyByID returns public key by given ID.
func (s *PublicKeysStore) GetPublicKeyByID(ctx context.Context, id int64) (*domain.PublicKey, error) {
	key := new(PublicKey)
	err := s.db.WithContext(ctx).Where("id = ?", id).First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKeyNotExist{id}
		}
		return nil, err
	}
	return key.ToDomain(), nil
}

// SearchPublicKeyByContent searches a public key using the content as prefix
// (i.e. ignore the email part). It returns ErrKeyNotExist if no such key
// exists.
func (s *PublicKeysStore) SearchPublicKeyByContent(ctx context.Context, content string) (*domain.PublicKey, error) {
	key := new(PublicKey)
	err := s.db.WithContext(ctx).Where("content like ?", content+"%").First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKeyNotExist{}
		}
		return nil, err
	}
	return key.ToDomain(), nil
}

// ListPublicKeys returns a list of public keys belongs to given user.
func (s *PublicKeysStore) ListPublicKeys(ctx context.Context, uid int64) ([]*domain.PublicKey, error) {
	keys := make([]*PublicKey, 0, 5)
	err := s.db.WithContext(ctx).Where("owner_id = ?", uid).Find(&keys).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.PublicKey, len(keys))
	for i, key := range keys {
		res[i] = key.ToDomain()
	}
	return res, nil
}

// UpdatePublicKey updates given public key.
func (s *PublicKeysStore) UpdatePublicKey(ctx context.Context, key *domain.PublicKey) error {
	m := FromDomainPublicKey(key)
	return s.db.WithContext(ctx).Model(m).Updates(m).Error
}

// DeletePublicKey deletes SSH key information both in database and authorized_keys file.
func (s *PublicKeysStore) DeletePublicKey(ctx context.Context, doer *domain.User, id int64) (err error) {
	key, err := s.GetPublicKeyByID(ctx, id)
	if err != nil {
		if IsErrKeyNotExist(err) {
			return nil
		}
		return fmt.Errorf("GetPublicKeyByID: %w", err)
	}

	// Check if user has access to delete this key.
	if !doer.IsAdmin && doer.ID != key.OwnerID {
		return ErrKeyAccessDenied{doer.ID, key.ID, "public"}
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return deletePublicKeys(tx, id)
	})
	if err != nil {
		return err
	}

	return s.RewriteAuthorizedKeys()
}

// AddDeployKey add new deploy key to database and authorized_keys file.
func (s *PublicKeysStore) AddDeployKey(ctx context.Context, repoID int64, name, content string, mode domain.AccessMode) (*domain.DeployKey, error) {
	if err := checkKeyContent(content); err != nil {
		return nil, err
	}

	pkey := &PublicKey{
		Content: content,
		Mode:    mode,
		Type:    domain.KeyTypeDeploy,
	}
	err := s.db.WithContext(ctx).Where("content = ? AND mode = ? AND type = ?", content, mode, domain.KeyTypeDeploy).First(pkey).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	has := err == nil

	var dk *DeployKey
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First time use this deploy key.
		if !has {
			if err = addKey(tx, pkey); err != nil {
				return fmt.Errorf("addKey: %w", err)
			}
		}

		dk, err = addDeployKey(tx, pkey.ID, repoID, name, pkey.Fingerprint, mode)
		if err != nil {
			return fmt.Errorf("addDeployKey: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dk.ToDomain(), nil
}

// GetDeployKeyByID returns deploy key by given ID.
func (s *PublicKeysStore) GetDeployKeyByID(ctx context.Context, id int64) (*domain.DeployKey, error) {
	key := new(DeployKey)
	err := s.db.WithContext(ctx).Where("id = ?", id).First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeployKeyNotExist{args: map[string]any{"deployKeyID": id}}
		}
		return nil, err
	}
	return key.ToDomain(), nil
}

// GetDeployKeyByRepo returns deploy key by given public key ID and repository ID.
func (s *PublicKeysStore) GetDeployKeyByRepo(ctx context.Context, keyID, repoID int64) (*domain.DeployKey, error) {
	key := new(DeployKey)
	err := s.db.WithContext(ctx).Where("key_id = ? AND repo_id = ?", keyID, repoID).First(key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeployKeyNotExist{args: map[string]any{"keyID": keyID, "repoID": repoID}}
		}
		return nil, err
	}
	return key.ToDomain(), nil
}

// UpdateDeployKey updates deploy key information.
func (s *PublicKeysStore) UpdateDeployKey(ctx context.Context, key *domain.DeployKey) error {
	m := FromDomainDeployKey(key)
	return s.db.WithContext(ctx).Model(m).Updates(m).Error
}

// DeleteDeployKey deletes deploy key from its repository authorized_keys file if needed.
func (s *PublicKeysStore) DeleteDeployKey(ctx context.Context, doer *domain.User, id int64) error {
	key, err := s.GetDeployKeyByID(ctx, id)
	if err != nil {
		if IsErrDeployKeyNotExist(err) {
			return nil
		}
		return fmt.Errorf("GetDeployKeyByID: %w", err)
	}

	// Check if user has access to delete this key.
	if !doer.IsAdmin {
		repo, err := Handle.Repositories().GetByID(ctx, key.RepoID)
		if err != nil {
			return fmt.Errorf("GetRepositoryByID: %w", err)
		}
		if !Handle.Permissions().Authorize(ctx, doer.ID, repo.ID, domain.AccessModeAdmin,
			domain.AccessModeOptions{
				OwnerID: repo.OwnerID,
				Private: repo.IsPrivate,
			},
		) {
			return ErrKeyAccessDenied{doer.ID, key.ID, "deploy"}
		}
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("id = ?", key.ID).Delete(new(DeployKey)).Error; err != nil {
			return fmt.Errorf("delete deploy key [%d]: %v", key.ID, err)
		}

		// Check if this is the last reference to same key content.
		var count int64
		err = tx.Model(new(DeployKey)).Where("key_id = ?", key.KeyID).Count(&count).Error
		if err != nil {
			return err
		} else if count == 0 {
			if err = deletePublicKeys(tx, key.KeyID); err != nil {
				return err
			}
		}
		return nil
	})
}

// ListDeployKeys returns all deploy keys by given repository ID.
func (s *PublicKeysStore) ListDeployKeys(ctx context.Context, repoID int64) ([]*domain.DeployKey, error) {
	keys := make([]*DeployKey, 0, 5)
	err := s.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&keys).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.DeployKey, len(keys))
	for i, key := range keys {
		res[i] = key.ToDomain()
	}
	return res, nil
}
