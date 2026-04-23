package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/errx"
	"github.com/nef734/ggit/internal/repox"

	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
)

// repoListFields are the columns used for listing repositories in the dashboard and explore pages.
// We avoid 'languages' and other heavy columns here to save RAM.
const repoListFields = "repository.id, repository.owner_id, repository.name, repository.lower_name, repository.description, repository.website, repository.default_branch, repository.size, repository.num_watches, repository.num_stars, repository.num_forks, repository.num_issues, repository.num_closed_issues, repository.num_pulls, repository.num_closed_pulls, repository.is_private, repository.is_fork, repository.is_unlisted, repository.is_bare, repository.enable_wiki, repository.wiki_folder_name, repository.enable_issues, repository.enable_pulls, repository.created_unix, repository.updated_unix"

type RepositoryAPIFormatOptions struct {
	Permission *apiv2types.RepositoryPermission
	Parent     *apiv2types.RepositoryResponse
}

// APIFormat returns the API format of a repository.
func (r *Repository) APIFormat(owner *User, opts ...RepositoryAPIFormatOptions) *apiv2types.RepositoryResponse {
	var opt RepositoryAPIFormatOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	var langs map[string]int64
	if r.Languages != "" {
		_ = json.Unmarshal([]byte(r.Languages), &langs)
	}

	cloneLink := repox.NewCloneLink(owner.Name, r.Name)

	return &apiv2types.RepositoryResponse{
		ID:            r.ID,
		Owner:         owner.APIFormat(),
		Name:          r.Name,
		FullName:      owner.Name + "/" + r.Name,
		Description:   r.Description,
		NumStars:      r.NumStars,
		NumForks:      r.NumForks,
		NumWatches:    r.NumWatches,
		OpenIssues:    r.NumOpenIssues,
		OpenPulls:     r.NumOpenPulls,
		IsPrivate:     r.IsPrivate,
		HTMLURL:       owner.HTMLURL() + "/" + r.Name,
		CloneURL:      cloneLink.HTTPS,
		SSHURL:        cloneLink.SSH,
		DefaultBranch: r.DefaultBranch,
		WikiFolderName: r.WikiFolderName,
		CreatedAt:     time.Unix(r.CreatedUnix, 0),
		UpdatedAt:     time.Unix(r.UpdatedUnix, 0),
		Permissions:   opt.Permission,
		Parent:        opt.Parent,
		Topics:        r.Topics,
	}
}

// RepositoriesStore is the storage layer for repositories.
type RepositoriesStore struct {
	db *gorm.DB
}

func newReposStore(db *gorm.DB) *RepositoriesStore {
	return &RepositoriesStore{db: db}
}

type ErrRepoAlreadyExist struct {
	args errx.Args
}

func IsErrRepoAlreadyExist(err error) bool {
	_, ok := err.(ErrRepoAlreadyExist)
	return ok
}

func (err ErrRepoAlreadyExist) Error() string {
	return fmt.Sprintf("repository already exists: %v", err.args)
}

type CreateRepoOptions struct {
	Name          string
	Description   string
	Website       string
	DefaultBranch string
	Private       bool
	IsUnlisted    bool
	AutoInit      bool
	Gitignores    string
	License       string
	Readme        string
	EnableWiki     bool
	WikiFolderName string
	EnableIssues   bool
	EnablePulls    bool
	Fork           bool
	ForkID         int64
	CreatedAt      time.Time
}

// Create creates a new repository record in the database. It returns
// ErrNameNotAllowed when the repository name is not allowed, or
// ErrRepoAlreadyExist when a repository with same name already exists for the
// owner.
func (s *RepositoriesStore) Create(ctx context.Context, ownerID int64, opts CreateRepoOptions) (*Repository, error) {
	err := isRepoNameAllowed(opts.Name)
	if err != nil {
		return nil, err
	}

	_, err = s.GetByName(ctx, ownerID, opts.Name)
	if err == nil {
		return nil, ErrRepoAlreadyExist{
			args: errx.Args{
				"ownerID": ownerID,
				"name":    opts.Name,
			},
		}
	} else if !IsErrRepoNotExist(err) {
		return nil, err
	}

	repo := &Repository{
		OwnerID:       ownerID,
		LowerName:     strings.ToLower(opts.Name),
		Name:          opts.Name,
		Description:   opts.Description,
		Website:       opts.Website,
		DefaultBranch: opts.DefaultBranch,
		IsPrivate:     opts.Private,
		IsUnlisted:    opts.IsUnlisted,
		EnableWiki:    opts.EnableWiki,
		WikiFolderName: opts.WikiFolderName,
		EnableIssues:  opts.EnableIssues,
		EnablePulls:   opts.EnablePulls,
		IsFork:        opts.Fork,
		ForkID:        opts.ForkID,
	}
	if !opts.CreatedAt.IsZero() {
		repo.CreatedUnix = opts.CreatedAt.Unix()
		repo.UpdatedUnix = repo.CreatedUnix
	}
	return repo, s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Create(repo).Error
		if err != nil {
			return errors.Wrap(err, "create")
		}

		err = newReposStore(tx).Watch(ctx, ownerID, repo.ID)
		if err != nil {
			return errors.Wrap(err, "watch")
		}
		return nil
	})
}

type UpdateRepoOptions struct {
	PrimaryLanguage *string
	Languages       *string
	WikiFolderName  *string
}

// Update updates the repository record in the database.
func (s *RepositoriesStore) Update(ctx context.Context, repoID int64, opts UpdateRepoOptions) error {
	updates := make(map[string]any)
	if opts.PrimaryLanguage != nil {
		updates["primary_language"] = *opts.PrimaryLanguage
	}
	if opts.Languages != nil {
		updates["languages"] = *opts.Languages
	}
	if opts.WikiFolderName != nil {
		updates["wiki_folder_name"] = *opts.WikiFolderName
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).Model(&Repository{}).Where("id = ?", repoID).Updates(updates).Error
}

// GetByCollaboratorID returns a list of repositories that the given
// collaborator has access to. Results are limited to the given limit and sorted
// by the given order (e.g. "updated_unix DESC"). Repositories that are owned
// directly by the given collaborator are not included.
func (s *RepositoriesStore) GetByCollaboratorID(ctx context.Context, collaboratorID int64, limit int, orderBy string) ([]*Repository, error) {
	/*
		Equivalent SQL for PostgreSQL:

		SELECT * FROM repository
		JOIN access ON access.repo_id = repository.id AND access.user_id = @collaboratorID
		WHERE access.mode >= @accessModeRead
		ORDER BY @orderBy
		LIMIT @limit
	*/
	var repos []*Repository
	return repos, s.db.WithContext(ctx).
		Select(repoListFields).
		Joins("JOIN access ON access.repo_id = repository.id AND access.user_id = ?", collaboratorID).
		Where("access.mode >= ?", AccessModeRead).
		Order(orderBy).
		Limit(limit).
		Find(&repos).
		Error
}

// GetByCollaboratorIDWithAccessMode returns a list of repositories and
// corresponding access mode that the given collaborator has access to.
// Repositories that are owned directly by the given collaborator are not
// included.
func (s *RepositoriesStore) GetByCollaboratorIDWithAccessMode(ctx context.Context, collaboratorID int64) (map[*Repository]AccessMode, error) {
	/*
		Equivalent SQL for PostgreSQL:

		SELECT
			repository.*,
			access.mode
		FROM repository
		JOIN access ON access.repo_id = repository.id AND access.user_id = @collaboratorID
		WHERE access.mode >= @accessModeRead
	*/
	var reposWithAccessMode []*struct {
		*Repository
		Mode AccessMode
	}
	err := s.db.WithContext(ctx).
		Select("repository.id, repository.owner_id, repository.name, repository.lower_name, repository.description, repository.is_private, repository.updated_unix", "access.mode").
		Table("repository").
		Joins("JOIN access ON access.repo_id = repository.id AND access.user_id = ?", collaboratorID).
		Where("access.mode >= ?", AccessModeRead).
		Find(&reposWithAccessMode).
		Error
	if err != nil {
		return nil, err
	}

	repos := make(map[*Repository]AccessMode, len(reposWithAccessMode))
	for _, repoWithAccessMode := range reposWithAccessMode {
		repos[repoWithAccessMode.Repository] = repoWithAccessMode.Mode
	}
	return repos, nil
}

var _ errx.NotFound = (*ErrRepoNotExist)(nil)

type ErrRepoNotExist struct {
	args errx.Args
}

func IsErrRepoNotExist(err error) bool {
	_, ok := err.(ErrRepoNotExist)
	return ok
}

func (err ErrRepoNotExist) Error() string {
	return fmt.Sprintf("repository does not exist: %v", err.args)
}

func (ErrRepoNotExist) NotFound() bool {
	return true
}

// GetByID returns the repository with given ID. It returns ErrRepoNotExist when
// not found.
func (s *RepositoriesStore) GetByID(ctx context.Context, id int64) (*Repository, error) {
	repo := new(Repository)
	err := s.db.WithContext(ctx).Where("id = ?", id).First(repo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRepoNotExist{errx.Args{"repoID": id}}
		}
		return nil, err
	}
	return repo, nil
}

// GetByName returns the repository with given owner and name. It returns
// ErrRepoNotExist when not found.
func (s *RepositoriesStore) GetByName(ctx context.Context, ownerID int64, name string) (*Repository, error) {
	repo := new(Repository)
	err := s.db.WithContext(ctx).
		Where("owner_id = ? AND lower_name = ?", ownerID, strings.ToLower(name)).
		First(repo).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRepoNotExist{
				args: errx.Args{
					"ownerID": ownerID,
					"name":    name,
				},
			}
		}
		return nil, err
	}
	return repo, nil
}

func (s *RepositoriesStore) recountStars(tx *gorm.DB, userID, repoID int64) error {
	/*
		Equivalent SQL for PostgreSQL:

		UPDATE repository
		SET num_stars = (
			SELECT COUNT(*) FROM star WHERE repo_id = @repoID
		)
		WHERE id = @repoID
	*/
	err := tx.Model(&Repository{}).
		Where("id = ?", repoID).
		Update(
			"num_stars",
			tx.Model(&Star{}).Select("COUNT(*)").Where("repo_id = ?", repoID),
		).
		Error
	if err != nil {
		return errors.Wrap(err, `update "repository.num_stars"`)
	}

	/*
		Equivalent SQL for PostgreSQL:

		UPDATE "user"
		SET num_stars = (
			SELECT COUNT(*) FROM star WHERE uid = @userID
		)
		WHERE id = @userID
	*/
	err = tx.Model(&User{}).
		Where("id = ?", userID).
		Update(
			"num_stars",
			tx.Model(&Star{}).Select("COUNT(*)").Where("uid = ?", userID),
		).
		Error
	if err != nil {
		return errors.Wrap(err, `update "user.num_stars"`)
	}
	return nil
}

// Star marks the user to star the repository.
func (s *RepositoriesStore) Star(ctx context.Context, userID, repoID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		star := &Star{
			UserID: userID,
			RepoID: repoID,
		}
		result := tx.FirstOrCreate(star, star)
		if result.Error != nil {
			return errors.Wrap(result.Error, "upsert")
		} else if result.RowsAffected <= 0 {
			return nil // Relation already exists
		}

		return s.recountStars(tx, userID, repoID)
	})
}

// Touch updates the updated time to the current time and removes the bare state
// of the given repository.
func (s *RepositoriesStore) Touch(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).
		Model(new(Repository)).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_bare":      false,
			"updated_unix": s.db.NowFunc().Unix(),
		}).
		Error
}

// ListWatches returns all watches of the given repository.
func (s *RepositoriesStore) ListWatches(ctx context.Context, repoID int64) ([]*Watch, error) {
	var watches []*Watch
	return watches, s.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&watches).Error
}

func (s *RepositoriesStore) recountWatches(tx *gorm.DB, repoID int64) error {
	/*
		Equivalent SQL for PostgreSQL:

		UPDATE repository
		SET num_watches = (
			SELECT COUNT(*) FROM watch WHERE repo_id = @repoID
		)
		WHERE id = @repoID
	*/
	return tx.Model(&Repository{}).
		Where("id = ?", repoID).
		Update(
			"num_watches",
			tx.Model(&Watch{}).Select("COUNT(*)").Where("repo_id = ?", repoID),
		).
		Error
}

// Watch marks the user to watch the repository.
func (s *RepositoriesStore) Watch(ctx context.Context, userID, repoID int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w := &Watch{
			UserID: userID,
			RepoID: repoID,
		}
		result := tx.FirstOrCreate(w, w)
		if result.Error != nil {
			return errors.Wrap(result.Error, "upsert")
		} else if result.RowsAffected <= 0 {
			return nil // Relation already exists
		}

		return s.recountWatches(tx, repoID)
	})
}

// HasForkedBy returns true if the given repository has forked by the given user.
func (s *RepositoriesStore) HasForkedBy(ctx context.Context, repoID, userID int64) bool {
	var count int64
	s.db.WithContext(ctx).Model(new(Repository)).Where("owner_id = ? AND fork_id = ?", userID, repoID).Count(&count)
	return count > 0
}

// ListStars returns a list of users that starred the given repository.
func (s *RepositoriesStore) ListStars(ctx context.Context, repoID int64, page, pageSize int) ([]*User, error) {
	users := make([]*User, 0, pageSize)
	return users, s.db.WithContext(ctx).
		Joins("JOIN star ON star.uid = user.id").
		Where("star.repo_id = ?", repoID).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Order("star.id DESC").
		Find(&users).
		Error
}

// ListForks returns a list of repositories that are forks of the given repository.
func (s *RepositoriesStore) ListForks(ctx context.Context, repoID int64, page, pageSize int) ([]*Repository, error) {
	repos := make([]*Repository, 0, pageSize)
	return repos, s.db.WithContext(ctx).
		Select(repoListFields).
		Where("fork_id = ?", repoID).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Order("updated_unix DESC").
		Find(&repos).
		Error
}

// ListStarred returns a list of repositories starred by the given user.
func (s *RepositoriesStore) ListStarred(ctx context.Context, userID int64, page, pageSize int) ([]*Repository, error) {
	repos := make([]*Repository, 0, pageSize)
	return repos, s.db.WithContext(ctx).
		Table("repository").
		Select(repoListFields).
		Joins("JOIN star ON star.repo_id = repository.id").
		Where("star.uid = ?", userID).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Order("star.id DESC").
		Find(&repos).
		Error
}
