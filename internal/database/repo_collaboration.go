package database

import (
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// Collaboration represent the relation between an individual and a repository.
type Collaboration struct {
	ID     int64      `gorm:"primaryKey"`
	UserID int64      `gorm:"uniqueIndex:collaboration_user_repo_unique;index;not null"`
	RepoID int64      `gorm:"uniqueIndex:collaboration_user_repo_unique;index;not null"`
	Mode   AccessMode `gorm:"not null;default:2"`
}

func (c *Collaboration) ModeI18nKey() string {
	switch c.Mode {
	case AccessModeRead:
		return "repo.settings.collaboration.read"
	case AccessModeWrite:
		return "repo.settings.collaboration.write"
	case AccessModeAdmin:
		return "repo.settings.collaboration.admin"
	default:
		return "repo.settings.collaboration.undefined"
	}
}

// IsCollaborator returns true if the user is a collaborator of the repository.
func IsCollaborator(repoID, userID int64) bool {
	var count int64
	err := Handle.db.Model(new(Collaboration)).Where("repo_id = ? AND user_id = ?", repoID, userID).Count(&count).Error
	if err != nil {
		log.Error("get collaboration [repo_id: %d, user_id: %d]: %v", repoID, userID, err)
		return false
	}
	return count > 0
}

func (r *Repository) IsCollaborator(userID int64) bool {
	return IsCollaborator(r.ID, userID)
}

// AddCollaborator adds new collaboration to a repository with default access mode.
func (r *Repository) AddCollaborator(u *User) error {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		collaboration := &Collaboration{
			RepoID: r.ID,
			UserID: u.ID,
		}

		err := tx.Where("repo_id = ? AND user_id = ?", r.ID, u.ID).First(collaboration).Error
		if err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		collaboration.Mode = AccessModeWrite
		if err = tx.Create(collaboration).Error; err != nil {
			return err
		}

		if err = r.recalculateAccesses(tx); err != nil {
			return errors.Newf("recalculateAccesses [repo_id: %v]: %v", r.ID, err)
		}
		return nil
	})
}

func (r *Repository) getCollaborations(tx *gorm.DB) ([]*Collaboration, error) {
	collaborations := make([]*Collaboration, 0, 5)
	err := tx.Where("repo_id = ?", r.ID).Find(&collaborations).Error
	return collaborations, err
}

// Collaborator represents a user with collaboration details.
type Collaborator struct {
	*User
	Collaboration *Collaboration
}

func (c *Collaborator) APIFormat() *apiv2types.RepositoryCollaboratorResponse {
	return &apiv2types.RepositoryCollaboratorResponse{
		UserResponse: c.User.APIFormat(),
		Permissions: apiv2types.RepositoryPermission{
			Admin: c.Collaboration.Mode >= AccessModeAdmin,
			Push:  c.Collaboration.Mode >= AccessModeWrite,
			Pull:  c.Collaboration.Mode >= AccessModeRead,
		},
	}
}

func (r *Repository) getCollaborators(tx *gorm.DB) ([]*Collaborator, error) {
	collaborations, err := r.getCollaborations(tx)
	if err != nil {
		return nil, errors.Newf("getCollaborations: %v", err)
	}

	collaborators := make([]*Collaborator, len(collaborations))
	for i, c := range collaborations {
		user, err := getUserByID(tx, c.UserID)
		if err != nil {
			return nil, err
		}
		collaborators[i] = &Collaborator{
			User:          user,
			Collaboration: c,
		}
	}
	return collaborators, nil
}

// GetCollaborators returns the collaborators for a repository
func (r *Repository) GetCollaborators() ([]*Collaborator, error) {
	return r.getCollaborators(Handle.db)
}

// ChangeCollaborationAccessMode sets new access mode for the collaboration.
func (r *Repository) ChangeCollaborationAccessMode(userID int64, mode AccessMode) error {
	// Discard invalid input
	if mode <= AccessModeNone || mode > AccessModeOwner {
		return nil
	}

	collaboration := new(Collaboration)
	err := Handle.db.Where("repo_id = ? AND user_id = ?", r.ID, userID).First(collaboration).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return errors.Newf("get collaboration: %v", err)
	}

	if collaboration.Mode == mode {
		return nil
	}
	collaboration.Mode = mode

	// If it's an organizational repository, merge with team access level for highest permission
	if r.Owner.IsOrganization() {
		teams, err := GetUserTeams(r.OwnerID, userID)
		if err != nil {
			return errors.Newf("GetUserTeams: [org_id: %d, user_id: %d]: %v", r.OwnerID, userID, err)
		}
		for i := range teams {
			if mode < teams[i].Authorize {
				mode = teams[i].Authorize
			}
		}
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(collaboration).Updates(collaboration).Error; err != nil {
			return errors.Newf("update collaboration: %v", err)
		}

		access := &Access{
			UserID: userID,
			RepoID: r.ID,
		}
		err = tx.Where("user_id = ? AND repo_id = ?", userID, r.ID).First(access).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Newf("get access record: %v", err)
		}

		if err == nil {
			return tx.Model(access).Update("mode", mode).Error
		}

		access.Mode = mode
		return tx.Create(access).Error
	})
}

// DeleteCollaboration removes collaboration relation between the user and repository.
func DeleteCollaboration(repo *Repository, userID int64) (err error) {
	if !IsCollaborator(repo.ID, userID) {
		return nil
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("repo_id = ? AND user_id = ?", repo.ID, userID).Delete(new(Collaboration)).Error; err != nil {
			return err
		}

		if err = repo.recalculateAccesses(tx); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) DeleteCollaboration(userID int64) error {
	return DeleteCollaboration(r, userID)
}
