package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbx"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
	"github.com/nef734/ggit/internal/fsx"
	"github.com/nef734/ggit/internal/lazyregexp"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/markup"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/repox"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/nef734/ggit/internal/strx"
	"github.com/nef734/ggit/internal/userx"
	"gorm.io/gorm"
)

// UserListFields are the columns used for listing Users/organizations.
// We avoid 'passwd', 'salt', 'rands', etc. for security and to save RAM.
const UserListFields = "User.id, User.lower_name, User.name, User.full_name, User.email, User.type, User.location, User.website, User.avatar, User.avatar_email, User.use_custom_avatar, User.num_stars, User.num_repos, User.description, User.num_teams, User.num_members, User.created_unix, User.updated_unix, User.is_admin, User.is_active, User.prohibit_login"

var (
	// alphaDashDotPattern represents a regex for valid Username/email characters.
	alphaDashDotPattern = lazyregexp.New(`[^\d\w\-_.]`)
)

// UsersStore is the storage layer for Users.
type UsersStore struct {
	db *gorm.DB
}

func NewUsersStore(db *gorm.DB) *UsersStore {
	return &UsersStore{db: db}
}

type ErrLoginSourceMismatch struct {
	args errx.Args
}

// IsErrLoginSourceMismatch returns true if the underlying error has the type
// ErrLoginSourceMismatch.
func IsErrLoginSourceMismatch(err error) bool {
	var target ErrLoginSourceMismatch
	return errors.As(err, &target)
}

func (err ErrLoginSourceMismatch) Error() string {
	return fmt.Sprintf("login source mismatch: %v", err.args)
}

// Authenticate validates Username and password via given login source ID. It
// returns ErrUserNotExist when the User was not found.
//
// When the "loginSourceID" is negative, it aborts the process and returns
// ErrUserNotExist if the User was not found in the database.
//
// When the "loginSourceID" is non-negative, it returns ErrLoginSourceMismatch
// if the User has different login source ID than the "loginSourceID".
//
// When the "loginSourceID" is positive, it tries to authenticate via given
// login source and creates a new User when not yet exists in the database.
func (s *UsersStore) Authenticate(ctx context.Context, login, password string, loginSourceID int64) (*domain.User, error) {
	login = strings.ToLower(login)

	query := s.db.WithContext(ctx)
	if strings.Contains(login, "@") {
		query = query.Where("email = ?", login)
	} else {
		query = query.Where("lower_name = ?", login)
	}

	u := new(User)
	err := query.First(u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get User"+": %v", err)
	}

	var authSourceID int64 // The login source ID will be used to authenticate the User
	createNewUser := false // Whether to create a new User after successful authentication

	// User found in the database
	if err == nil {
		// Note: This check is unnecessary but to reduce User confusion at login page
		// and make it more consistent from User's perspective.
		if loginSourceID >= 0 && u.LoginSource != loginSourceID {
			return nil, ErrLoginSourceMismatch{args: errx.Args{"expect": loginSourceID, "actual": u.LoginSource}}
		}

		// Validate password hash fetched from database for local accounts.
		if u.LoginSource <= 0 { // IsLocal
			if userx.ValidatePassword(u.Password, u.Salt, password) {
				return u.ToDomain(), nil
			}

			return nil, auth.ErrBadCredentials{Args: map[string]any{"login": login, "userID": u.ID}}
		}

		authSourceID = u.LoginSource

	} else {
		// Non-local login source is always greater than 0.
		if loginSourceID <= 0 {
			return nil, auth.ErrBadCredentials{Args: map[string]any{"login": login}}
		}

		authSourceID = loginSourceID
		createNewUser = true
	}

	source, err := newLoginSourcesStore(s.db, loadedLoginSourceFilesStore).getByID(ctx, authSourceID)

	if err != nil {
		return nil, fmt.Errorf("get login source"+": %v", err)
	}

	if !source.IsActived {
		return nil, fmt.Errorf("login source %d is not activated", source.ID)
	}

	extAccount, err := source.Provider.Authenticate(login, password)
	if err != nil {
		return nil, err
	}

	if !createNewUser {
		return u.ToDomain(), nil
	}

	// Validate Username make sure it satisfies requirement.
	if alphaDashDotPattern.MatchString(extAccount.Name) {
		return nil, fmt.Errorf("invalid pattern for attribute 'Username' [%s]: must be valid alpha or numeric or dash(-_) or dot characters", extAccount.Name)
	}

	createdUser, err := s.Create(ctx, extAccount.Name, extAccount.Email,
		CreateUserOptions{
			FullName:    extAccount.FullName,
			LoginSource: authSourceID,
			LoginName:   extAccount.Login,
			Location:    extAccount.Location,
			Website:     extAccount.Website,
			Activated:   true,
			Admin:       extAccount.Admin,
		})
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// ChangeUsername changes the Username of the given User and updates all
// references to the old Username. It returns ErrNameNotAllowed if the given
// name or pattern of the name is not allowed as a Username, or
// ErrUserAlreadyExist when another User with same name already exists.
func (s *UsersStore) ChangeUsername(ctx context.Context, UserID int64, newUsername string) error {
	err := isUsernameAllowed(newUsername)
	if err != nil {
		return err
	}

	if s.IsUsernameUsed(ctx, newUsername, UserID) {
		return ErrUserAlreadyExist{
			args: errx.Args{
				"name": newUsername,
			},
		}
	}

	u, err := s.getByID(ctx, UserID)
	if err != nil {
		return fmt.Errorf("get User"+": %v", err)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&User{}).
			Where("id = ?", u.ID).
			Updates(map[string]any{
				"lower_name":   strings.ToLower(newUsername),
				"name":         newUsername,
				"updated_unix": tx.NowFunc().Unix(),
			}).Error
		if err != nil {
			return fmt.Errorf("update User name"+": %v", err)
		}

		// Stop here if it's just a case-change of the Username
		if strings.EqualFold(u.Name, newUsername) {
			return nil
		}

		// Update all references to the User name in pull requests
		err = tx.Model(&PullRequest{}).
			Where("head_User_name = ?", u.LowerName).
			Update("head_User_name", strings.ToLower(newUsername)).
			Error
		if err != nil {
			return fmt.Errorf(`update "pull_request.head_User_name"`+": %v", err)
		}

		// Delete local copies of repositories and their wikis that are owned by the User
		rows, err := tx.Model(&Repository{}).Where("owner_id = ?", u.ID).Rows()
		if err != nil {
			return fmt.Errorf("iterate repositories"+": %v", err)
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var repo struct {
				ID int64
			}
			err = tx.ScanRows(rows, &repo)
			if err != nil {
				return fmt.Errorf("scan rows"+": %v", err)
			}

			deleteRepoLocalCopy(repo.ID)
		}
		if err = rows.Err(); err != nil {
			return fmt.Errorf("check rows.Err"+": %v", err)
		}

		// Rename User directory if exists
		UserPath := repox.UserPath(u.Name)
		if osx.Exist(UserPath) {
			newUserPath := repox.UserPath(newUsername)
			err = os.Rename(UserPath, newUserPath)
			if err != nil {
				return fmt.Errorf("rename User directory"+": %v", err)
			}
		}
		return nil
	})
}

// Count returns the total number of Users.
func (s *UsersStore) Count(ctx context.Context) int64 {
	var count int64
	s.db.WithContext(ctx).Model(&User{}).Where("type = ?", UserTypeIndividual).Count(&count)
	return count
}

type CreateUserOptions struct {
	FullName         string
	Password         string
	LoginSource      int64
	LoginName        string
	Location         string
	Website          string
	Activated        bool
	Admin            bool
	ProhibitLogin    bool
	AllowGitHook     bool
	AllowImportLocal bool
}

type ErrUserAlreadyExist struct {
	args errx.Args
}

// IsErrUserAlreadyExist returns true if the underlying error has the type
// ErrUserAlreadyExist.
func IsErrUserAlreadyExist(err error) bool {
	var target ErrUserAlreadyExist
	return errors.As(err, &target)
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("User already exists: %v", err.args)
}

type ErrEmailAlreadyUsed struct {
	args errx.Args
}

// IsErrEmailAlreadyUsed returns true if the underlying error has the type
// ErrEmailAlreadyUsed.
func IsErrEmailAlreadyUsed(err error) bool {
	var target ErrEmailAlreadyUsed
	return errors.As(err, &target)
}

func (err ErrEmailAlreadyUsed) Email() string {
	email, ok := err.args["email"].(string)
	if ok {
		return email
	}
	return "<email not found>"
}

func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("email has been used: %v", err.args)
}

func (s *UsersStore) Create(ctx context.Context, Username, email string, opts any) (*domain.User, error) {
	var options CreateUserOptions
	if o, ok := opts.(CreateUserOptions); ok {
		options = o
	} else {
		return nil, errors.New("invalid options type")
	}

	err := isUsernameAllowed(Username)
	if err != nil {
		return nil, err
	}

	if s.IsUsernameUsed(ctx, Username, 0) {
		return nil, ErrUserAlreadyExist{
			args: errx.Args{
				"name": Username,
			},
		}
	}

	email = strings.ToLower(strings.TrimSpace(email))
	_, err = s.getByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyUsed{
			args: errx.Args{
				"email": email,
			},
		}
	} else if !IsErrUserNotExist(err) {
		return nil, err
	}

	isAdmin := options.Admin
	isActivated := options.Activated

	// The first User created must be an administrator and activated
	if s.Count(ctx) == 0 {
		isAdmin = true
		isActivated = true
	}

	u := &User{
		LowerName:        strings.ToLower(Username),
		Name:             Username,
		FullName:         options.FullName,
		Email:            email,
		Password:         options.Password,
		LoginSource:      options.LoginSource,
		LoginName:        options.LoginName,
		Location:         options.Location,
		Website:          options.Website,
		MaxRepoCreation:  -1,
		IsActive:         isActivated,
		IsAdmin:          isAdmin,
		AllowGitHook:     options.AllowGitHook,
		AllowImportLocal: options.AllowImportLocal,
		ProhibitLogin:    options.ProhibitLogin,
		Avatar:           "",
		AvatarEmail:      email,
	}

	u.Rands, err = userx.RandomSalt()
	if err != nil {
		return nil, err
	}
	u.Salt, err = userx.RandomSalt()
	if err != nil {
		return nil, err
	}
	u.Password = userx.EncodePassword(u.Password, u.Salt)

	err = s.db.WithContext(ctx).Create(u).Error
	return u.ToDomain(), err
}

// DeleteCustomAvatar deletes the current User custom avatar and falls back to
// use look up avatar by email.
func (s *UsersStore) DeleteCustomAvatar(ctx context.Context, UserID int64) error {
	_ = os.Remove(userx.CustomAvatarPath(UserID))
	return s.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", UserID).
		Updates(map[string]any{
			"use_custom_avatar": false,
			"updated_unix":      s.db.NowFunc().Unix(),
		}).
		Error
}

type ErrUserOwnRepos struct {
	args errx.Args
}

// IsErrUserOwnRepos returns true if the underlying error has the type
// ErrUserOwnRepos.
func IsErrUserOwnRepos(err error) bool {
	var target ErrUserOwnRepos
	return errors.As(err, &target)
}

func (err ErrUserOwnRepos) Error() string {
	return fmt.Sprintf("User still has repository ownership: %v", err.args)
}

type ErrUserHasOrgs struct {
	args errx.Args
}

// IsErrUserHasOrgs returns true if the underlying error has the type
// ErrUserHasOrgs.
func IsErrUserHasOrgs(err error) bool {
	var target ErrUserHasOrgs
	return errors.As(err, &target)
}

func (err ErrUserHasOrgs) Error() string {
	return fmt.Sprintf("User still has organization membership: %v", err.args)
}

// DeleteByID deletes the given User and all their resources. It returns
// ErrUserOwnRepos when the User still has repository ownership, or returns
// ErrUserHasOrgs when the User still has organization membership. It is more
// performant to skip rewriting the "authorized_keys" file for individual
// deletion in a batch operation.
func (s *UsersStore) DeleteByID(ctx context.Context, UserID int64, skipRewriteAuthorizedKeys bool) error {
	u, err := s.getByID(ctx, UserID)
	if err != nil {
		if IsErrUserNotExist(err) {
			return nil
		}
		return fmt.Errorf("get User"+": %v", err)
	}

	// Double-check the User is not a direct owner of any repository and not a
	// member of any organization.
	var count int64
	err = s.db.WithContext(ctx).Model(&Repository{}).Where("owner_id = ?", UserID).Count(&count).Error
	if err != nil {
		return fmt.Errorf("count repositories"+": %v", err)
	} else if count > 0 {
		return ErrUserOwnRepos{args: errx.Args{"userID": UserID}}
	}

	err = s.db.WithContext(ctx).Model(&OrgUser{}).Where("uid = ?", UserID).Count(&count).Error
	if err != nil {
		return fmt.Errorf("count organization membership"+": %v", err)
	} else if count > 0 {
		return ErrUserHasOrgs{args: errx.Args{"userID": UserID}}
	}

	needsRewriteAuthorizedKeys := false
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		/*
			Equivalent SQL for PostgreSQL:

			UPDATE repository
			SET num_watches = num_watches - 1
			WHERE id IN (
				SELECT repo_id FROM watch WHERE User_id = @UserID
			)
		*/
		err = tx.Table("repository").
			Where("id IN (?)", tx.
				Select("repo_id").
				Table("watch").
				Where("User_id = ?", UserID),
			).
			UpdateColumn("num_watches", gorm.Expr("num_watches - 1")).
			Error
		if err != nil {
			return fmt.Errorf(`decrease "repository.num_watches"`+": %v", err)
		}

		/*
			Equivalent SQL for PostgreSQL:

			UPDATE repository
			SET num_stars = num_stars - 1
			WHERE id IN (
				SELECT repo_id FROM star WHERE uid = @UserID
			)
		*/
		err = tx.Table("repository").
			Where("id IN (?)", tx.
				Select("repo_id").
				Table("star").
				Where("uid = ?", UserID),
			).
			UpdateColumn("num_stars", gorm.Expr("num_stars - 1")).
			Error
		if err != nil {
			return fmt.Errorf(`decrease "repository.num_stars"`+": %v", err)
		}

		if !skipRewriteAuthorizedKeys {
			// We need to rewrite "authorized_keys" file if the User owns any public keys.
			needsRewriteAuthorizedKeys = tx.Where("owner_id = ?", UserID).First(&PublicKey{}).Error != gorm.ErrRecordNotFound
		}

		err = tx.Model(&Issue{}).Where("assignee_id = ?", UserID).Update("assignee_id", 0).Error
		if err != nil {
			return fmt.Errorf("clear assignees"+": %v", err)
		}

		for _, t := range []struct {
			table any
			where string
		}{
			{&Watch{}, "User_id = @UserID"},
			{&Star{}, "uid = @UserID"},
			{&PublicKey{}, "owner_id = @UserID"},

			{&AccessToken{}, "uid = @UserID"},
			{&Collaboration{}, "User_id = @UserID"},
			{&Access{}, "User_id = @UserID"},
			{&Action{}, "User_id = @UserID"},
			{&IssueUser{}, "uid = @UserID"},
			{&EmailAddress{}, "uid = @UserID"},
			{&User{}, "id = @UserID"},
		} {
			err = tx.Where(t.where, sql.Named("UserID", UserID)).Delete(t.table).Error
			if err != nil {
				return fmt.Errorf("clean up table %T: %v", t.table, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	_ = os.RemoveAll(repox.UserPath(u.Name))
	_ = os.Remove(userx.CustomAvatarPath(UserID))

	if needsRewriteAuthorizedKeys {
		err = newPublicKeysStore(s.db).RewriteAuthorizedKeys()
		if err != nil {
			return fmt.Errorf(`rewrite "authorized_keys" file`+": %v", err)
		}
	}
	return nil
}

// DeleteInactivated deletes all inactivated Users.
//
// NOTE: We do not take context.Context here because this operation in practice
// could much longer than the general request timeout (e.g. one minute).
func (s *UsersStore) DeleteInactivated(ctx context.Context) error {
	var UserIDs []int64
	err := s.db.Model(&User{}).Where("is_active = ?", false).Pluck("id", &UserIDs).Error
	if err != nil {
		return fmt.Errorf("get inactivated User IDs"+": %v", err)
	}

	for _, UserID := range UserIDs {
		err = s.DeleteByID(context.Background(), UserID, true)
		if err != nil {
			// Skip Users that may had set to inactivated by admins.
			if IsErrUserOwnRepos(err) || IsErrUserHasOrgs(err) {
				continue
			}
			return fmt.Errorf("delete User with ID %d: %v", UserID, err)
		}
	}
	err = newPublicKeysStore(s.db).RewriteAuthorizedKeys()
	if err != nil {
		return fmt.Errorf(`rewrite "authorized_keys" file`+": %v", err)
	}
	return nil
}

var _ errx.NotFound = (*ErrUserNotExist)(nil)

type ErrUserNotExist struct {
	args errx.Args
}

// IsErrUserNotExist returns true if the underlying error has the type
// ErrUserNotExist.
func IsErrUserNotExist(err error) bool {
	var target ErrUserNotExist
	return errors.As(err, &target)
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("User does not exist: %v", err.args)
}

func (ErrUserNotExist) NotFound() bool {
	return true
}

func (s *UsersStore) getByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, ErrUserNotExist{args: errx.Args{"email": email}}
	}
	email = strings.ToLower(email)

	var u User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist{args: errx.Args{"email": email}}
		}
		return nil, fmt.Errorf("get user by email"+": %v", err)
	}
	return &u, nil
}

// GetByLoginSource returns the user with given login source ID and login name.
func (s *UsersStore) GetByLoginSource(ctx context.Context, loginSourceID int64, loginName string) (*domain.User, error) {
	var u User
	err := s.db.WithContext(ctx).Where("login_source = ? AND login_name = ?", loginSourceID, loginName).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist{args: errx.Args{"login_source": loginSourceID, "login_name": loginName}}
		}
		return nil, fmt.Errorf("get user by login source"+": %v", err)
	}
	return u.ToDomain(), nil
}

// GetByEmail returns the User (not organization) with given email. It ignores
// records with unverified emails and returns ErrUserNotExist when not found.
func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, ErrUserNotExist{args: errx.Args{"email": email}}
	}
	email = strings.ToLower(email)

	// 1. Check primary email
	u := new(User)
	err := s.db.WithContext(ctx).
		Where("email = ? AND type = ? AND is_active = ?", email, UserTypeIndividual, true).
		First(u).Error
	if err == nil {
		return u.ToDomain(), nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get user by primary email: %w", err)
	}

	// 2. Check secondary emails
	emailAddr := new(EmailAddress)
	err = s.db.WithContext(ctx).
		Where("email = ? AND is_activated = ?", email, true).
		First(emailAddr).Error
	if err == nil {
		return s.GetByID(ctx, emailAddr.UserID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get user by secondary email: %w", err)
	}

	return nil, ErrUserNotExist{args: errx.Args{"email": email}}
}

func (s *UsersStore) getByID(ctx context.Context, id int64) (*User, error) {
	u := new(User)
	err := s.db.WithContext(ctx).Where("id = ?", id).First(u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist{args: errx.Args{"userID": id}}
		}
		return nil, err
	}
	return u, nil
}

// GetByID returns the User with given ID. It returns ErrUserNotExist when not
// found.
func (s *UsersStore) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	u, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (s *UsersStore) getByUsername(ctx context.Context, username string) (*User, error) {
	u := new(User)
	err := s.db.WithContext(ctx).Where("lower_name = ?", strings.ToLower(username)).First(u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist{args: errx.Args{"name": username}}
		}
		return nil, err
	}
	return u, nil
}

// GetByUsername returns the User with given Username. It returns
// ErrUserNotExist when not found.
func (s *UsersStore) GetByUsername(ctx context.Context, Username string) (*domain.User, error) {
	u, err := s.getByUsername(ctx, Username)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

// GetByKeyID returns the owner of given public key ID. It returns
// ErrUserNotExist when not found.
func (s *UsersStore) GetByKeyID(ctx context.Context, keyID int64) (*domain.User, error) {
	u := new(User)
	err := s.db.WithContext(ctx).
		Joins(dbx.Quote("JOIN public_key ON public_key.owner_id = %s.id", "User")).
		Where("public_key.id = ?", keyID).
		First(u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist{args: errx.Args{"keyID": keyID}}
		}
		return nil, err
	}
	return u.ToDomain(), nil
}

// GetMailableEmailsByUsernames returns a list of verified primary email
// addresses (where email notifications are sent to) of Users with given list of
// Usernames. Non-existing Usernames are ignored.
func (s *UsersStore) GetMailableEmailsByUsernames(ctx context.Context, Usernames []string) ([]string, error) {
	emails := make([]string, 0, len(Usernames))
	return emails, s.db.WithContext(ctx).
		Model(&User{}).
		Select("email").
		Where("lower_name IN (?) AND is_active = ?", Usernames, true).
		Find(&emails).Error
}

// IsUsernameUsed returns true if the given Username has been used other than
// the excluded User (a non-positive ID effectively meaning check against all
// Users).
func (s *UsersStore) IsUsernameUsed(ctx context.Context, Username string, excludeUserID int64) bool {
	if Username == "" {
		return false
	}

	err := s.db.WithContext(ctx).
		Select("id").
		Where("lower_name = ? AND id != ?", strings.ToLower(Username), excludeUserID).
		First(&User{}).Error

	if err == nil {
		return true
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A lookup that failed is not evidence that the name is taken. Reporting
		// it as taken turns any transient database error, a locked SQLite file or
		// a missing table, into "User already exists", which is both wrong and
		// hides the real cause behind a message about duplicates.
		//
		// Answering false here is safe: user.lower_name carries a unique index, so
		// a duplicate is still rejected at INSERT time, with the real error.
		log.Error("IsUsernameUsed(%q): %v", Username, err)
	}
	return false
}

// List returns a list of Users. Results are paginated by given page and page
// size, and sorted by primary key (id) in ascending order.
func (s *UsersStore) List(ctx context.Context, page, pageSize int) ([]*domain.User, error) {
	Users := make([]*User, 0, pageSize)
	err := s.db.WithContext(ctx).
		Select(UserListFields).
		Where("type = ?", UserTypeIndividual).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Order("id ASC").
		Find(&Users).
		Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.User, len(Users))
	for i, u := range Users {
		res[i] = u.ToDomain()
	}
	return res, nil
}

func searchUserByName(ctx context.Context, db *gorm.DB, UserType UserType, keyword string, page, pageSize int, orderBy string) ([]*domain.User, int64, error) {
	if keyword == "" {
		return []*domain.User{}, 0, nil
	}
	keyword = "%" + strings.ToLower(keyword) + "%"

	tx := db.WithContext(ctx).
		Where("type = ? AND (lower_name LIKE ? OR LOWER(full_name) LIKE ?)", UserType, keyword, keyword)

	var count int64
	err := tx.Model(&User{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count"+": %v", err)
	}

	Users := make([]*User, 0, pageSize)
	err = tx.Select(UserListFields).Order(orderBy).Limit(pageSize).Offset((page - 1) * pageSize).Find(&Users).Error
	if err != nil {
		return nil, 0, err
	}

	res := make([]*domain.User, len(Users))
	for i, u := range Users {
		res[i] = u.ToDomain()
	}
	return res, count, nil
}

// SearchByName returns a list of Users whose Username or full name matches the
// given keyword case-insensitively. Results are paginated by given page and
// page size, and sorted by the given order (e.g. "id DESC"). A total count of
// all results is also returned. If the order is not given, it's up to the
// database to decide.
func (s *UsersStore) SearchByName(ctx context.Context, keyword string, page, pageSize int, orderBy string) ([]*domain.User, int64, error) {
	return searchUserByName(ctx, s.db, UserTypeIndividual, keyword, page, pageSize, orderBy)
}

type UpdateUserOptions struct {
	LoginSource *int64
	LoginName   *string

	Password *string
	// GenerateNewRands indicates whether to force generate new rands for the User.
	GenerateNewRands bool

	FullName    *string
	Email       *string
	Website     *string
	Location    *string
	Description *string

	MaxRepoCreation    *int
	LastRepoVisibility *bool

	IsActivated      *bool
	IsAdmin          *bool
	AllowGitHook     *bool
	AllowImportLocal *bool
	ProhibitLogin    *bool

	Avatar      *string
	AvatarEmail *string

	DisableNewsFeed *bool

	IsPrivate *bool
}

// Update updates fields for the given User.
func (s *UsersStore) Update(ctx context.Context, UserID int64, opts any) error {
	var options UpdateUserOptions
	if o, ok := opts.(UpdateUserOptions); ok {
		options = o
	} else {
		return errors.New("invalid options type")
	}

	updates := map[string]any{
		"updated_unix": s.db.NowFunc().Unix(),
	}

	if options.IsPrivate != nil {
		if *options.IsPrivate {
			updates["type"] = UserTypeOrganizationPrivate
		} else {
			updates["type"] = UserTypeOrganization
		}
	}

	if options.LoginSource != nil {
		updates["login_source"] = *options.LoginSource
	}
	if options.LoginName != nil {
		updates["login_name"] = *options.LoginName
	}

	if options.Password != nil {
		salt, err := userx.RandomSalt()
		if err != nil {
			return fmt.Errorf("generate salt"+": %v", err)
		}
		updates["salt"] = salt
		updates["passwd"] = userx.EncodePassword(*options.Password, salt)
		options.GenerateNewRands = true
	}
	if options.GenerateNewRands {
		rands, err := userx.RandomSalt()
		if err != nil {
			return fmt.Errorf("generate rands"+": %v", err)
		}
		updates["rands"] = rands
	}

	if options.FullName != nil {
		updates["full_name"] = strx.Truncate(*options.FullName, 255)
	}
	if options.Email != nil {
		_, err := s.getByEmail(ctx, *options.Email)
		if err == nil {
			return ErrEmailAlreadyUsed{args: errx.Args{"email": *options.Email}}
		} else if !IsErrUserNotExist(err) {
			return err
		}
		updates["email"] = *options.Email
	}
	if options.Website != nil {
		updates["website"] = strx.Truncate(*options.Website, 255)
	}
	if options.Location != nil {
		updates["location"] = strx.Truncate(*options.Location, 255)
	}
	if options.Description != nil {
		updates["description"] = strx.Truncate(*options.Description, 255)
	}

	if options.MaxRepoCreation != nil {
		updates["max_repo_creation"] = *options.MaxRepoCreation
	}
	if options.LastRepoVisibility != nil {
		updates["last_repo_visibility"] = *options.LastRepoVisibility
	}

	if options.IsActivated != nil {
		updates["is_active"] = *options.IsActivated
	}
	if options.IsAdmin != nil {
		updates["is_admin"] = *options.IsAdmin
	}
	if options.AllowGitHook != nil {
		updates["allow_git_hook"] = *options.AllowGitHook
	}
	if options.AllowImportLocal != nil {
		updates["allow_import_local"] = *options.AllowImportLocal
	}
	if options.ProhibitLogin != nil {
		updates["prohibit_login"] = *options.ProhibitLogin
	}

	if options.Avatar != nil {
		updates["avatar"] = strx.Truncate(*options.Avatar, 2048)
	}
	if options.AvatarEmail != nil {
		updates["avatar_email"] = strx.Truncate(*options.AvatarEmail, 255)
	}

	if options.DisableNewsFeed != nil {
		updates["disable_news_feed"] = *options.DisableNewsFeed
	}

	return s.db.WithContext(ctx).Model(&User{}).Where("id = ?", UserID).Updates(updates).Error
}

// UseCustomAvatar uses the given avatar as the User custom avatar.
func (s *UsersStore) UseCustomAvatar(ctx context.Context, UserID int64, avatar []byte) error {
	avatarPath := userx.CustomAvatarPath(UserID)
	tempPath := avatarPath + ".tmp"

	fsTx := fsx.NewTransaction()
	defer fsTx.Rollback()

	if err := userx.SaveAvatarToPath(tempPath, avatar); err != nil {
		return fmt.Errorf("save avatar to temp"+": %v", err)
	}
	fsTx.StageCreate(tempPath, avatarPath)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&User{}).
			Where("id = ?", UserID).
			Updates(map[string]any{
				"use_custom_avatar": true,
				"updated_unix":      tx.NowFunc().Unix(),
			}).
			Error
	})
	if err != nil {
		return err
	}

	return fsTx.Commit()
}

// AddEmail adds a new email address to given User. It returns
// ErrEmailAlreadyUsed if the email has been verified by another User.
func (s *UsersStore) AddEmail(ctx context.Context, UserID int64, email string, isActivated bool) error {
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := s.GetByEmail(ctx, email)
	if err == nil {
		return ErrEmailAlreadyUsed{
			args: errx.Args{
				"email": email,
			},
		}
	} else if !IsErrUserNotExist(err) {
		return fmt.Errorf("check User by email"+": %v", err)
	}

	return s.db.WithContext(ctx).Create(
		&EmailAddress{
			UserID:      UserID,
			Email:       email,
			IsActivated: isActivated,
		},
	).Error
}

var _ errx.NotFound = (*ErrEmailNotExist)(nil)

type ErrEmailNotExist struct {
	args errx.Args
}

// IsErrEmailNotExist returns true if the underlying error has the type
// ErrEmailNotExist.
func IsErrEmailNotExist(err error) bool {
	var target ErrEmailNotExist
	return errors.As(err, &target)
}

func (err ErrEmailNotExist) Error() string {
	return fmt.Sprintf("email address does not exist: %v", err.args)
}

func (ErrEmailNotExist) NotFound() bool {
	return true
}

// GetEmail returns the email address of the given User. If `needsActivated` is
// true, only activated email will be returned, otherwise, it may return
// inactivated email addresses. It returns ErrEmailNotExist when no qualified
// email is not found.
func (s *UsersStore) GetEmail(ctx context.Context, UserID int64, email string, needsActivated bool) (*domain.EmailAddress, error) {
	tx := s.db.WithContext(ctx).Where("uid = ? AND email = ?", UserID, email)
	if needsActivated {
		tx = tx.Where("is_activated = ?", true)
	}

	emailAddress := new(EmailAddress)
	err := tx.First(emailAddress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmailNotExist{
				args: errx.Args{
					"email": email,
				},
			}
		}
		return nil, err
	}
	return emailAddress.ToDomain(), nil
}

// ListEmails returns all email addresses of the given User. It always includes
// a primary email address.
func (s *UsersStore) ListEmails(ctx context.Context, UserID int64) ([]*domain.EmailAddress, error) {
	u, err := s.getByID(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("get User"+": %v", err)
	}

	var emails []*EmailAddress
	err = s.db.WithContext(ctx).Where("uid = ?", UserID).Order("id ASC").Find(&emails).Error
	if err != nil {
		return nil, fmt.Errorf("list emails"+": %v", err)
	}

	isPrimaryFound := false
	for _, email := range emails {
		if email.Email == u.Email {
			isPrimaryFound = true
			email.IsPrimary = true
		}
	}

	res := make([]*domain.EmailAddress, 0, len(emails)+1)
	if !isPrimaryFound {
		res = append(res, &domain.EmailAddress{
			UserID:      u.ID,
			Email:       u.Email,
			IsActivated: u.IsActive,
			IsPrimary:   true,
		})
	}

	for _, email := range emails {
		res = append(res, email.ToDomain())
	}
	return res, nil
}

// MarkEmailActivated marks the email address of the given User as activated,
// and new rands are generated for the User.
func (s *UsersStore) MarkEmailActivated(ctx context.Context, UserID int64, email string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := s.db.WithContext(ctx).
			Model(&EmailAddress{}).
			Where("uid = ? AND email = ?", UserID, email).
			Update("is_activated", true).
			Error
		if err != nil {
			return fmt.Errorf("mark email activated"+": %v", err)
		}

		return NewUsersStore(tx).Update(ctx, UserID, UpdateUserOptions{GenerateNewRands: true})
	})
}

type ErrEmailNotVerified struct {
	args errx.Args
}

// IsErrEmailNotVerified returns true if the underlying error has the type
// ErrEmailNotVerified.
func IsErrEmailNotVerified(err error) bool {
	var target ErrEmailNotVerified
	return errors.As(err, &target)
}

func (err ErrEmailNotVerified) Error() string {
	return fmt.Sprintf("email has not been verified: %v", err.args)
}

// MarkEmailPrimary marks the email address of the given User as primary. It
// returns ErrEmailNotExist when the email is not found for the User, and
// ErrEmailNotActivated when the email is not activated.
func (s *UsersStore) MarkEmailPrimary(ctx context.Context, UserID int64, email string) error {
	var emailAddress EmailAddress
	err := s.db.WithContext(ctx).Where("uid = ? AND email = ?", UserID, email).First(&emailAddress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEmailNotExist{args: errx.Args{"email": email}}
		}
		return fmt.Errorf("get email address"+": %v", err)
	}

	if !emailAddress.IsActivated {
		return ErrEmailNotVerified{args: errx.Args{"email": email}}
	}

	u, err := s.getByID(ctx, UserID)
	if err != nil {
		return fmt.Errorf("get User"+": %v", err)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Make sure the former primary email doesn't disappear.
		err = tx.FirstOrCreate(
			&EmailAddress{
				UserID:      u.ID,
				Email:       u.Email,
				IsActivated: u.IsActive,
			},
			&EmailAddress{
				UserID: u.ID,
				Email:  u.Email,
			},
		).Error
		if err != nil {
			return fmt.Errorf("upsert former primary email address"+": %v", err)
		}

		return tx.Model(&User{}).
			Where("id = ?", u.ID).
			Updates(map[string]any{
				"email":        email,
				"updated_unix": tx.NowFunc().Unix(),
			},
			).Error
	})
}

// DeleteEmail deletes the email address of the given User.
func (s *UsersStore) DeleteEmail(ctx context.Context, UserID int64, email string) error {
	return s.db.WithContext(ctx).Where("uid = ? AND email = ?", UserID, email).Delete(&EmailAddress{}).Error
}

// UserType indicates the type of the User account.
type UserType int

const (
	UserTypeIndividual UserType = iota // NOTE: Historic reason to make it starts at 0.
	UserTypeOrganization
	UserTypeOrganizationPrivate
)

// User represents the object of an individual or an organization.
type User struct {
	ID        int64  `gorm:"primaryKey"`
	LowerName string `gorm:"unique;not null"`
	Name      string `gorm:"not null"`
	FullName  string
	// Email is the primary email address (to be used for communication)
	Email       string `gorm:"not null"`
	Password    string `gorm:"column:passwd;not null"`
	LoginSource int64  `gorm:"not null;default:0"`
	LoginName   string
	Type        UserType
	Location    string
	Website     string
	Rands       string `gorm:"type:VARCHAR(10)"`
	Salt        string `gorm:"type:VARCHAR(10)"`

	Created     time.Time `gorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `gorm:"-" json:"-"`
	UpdatedUnix int64

	// Remember visibility choice for convenience, true for private
	LastRepoVisibility bool
	// Maximum repository creation limit, -1 means use global default
	MaxRepoCreation int `gorm:"not null;default:-1"`

	// Permissions
	IsActive         bool // Activate primary email
	IsAdmin          bool
	AllowGitHook     bool
	AllowImportLocal bool // Allow migrate repository by local path
	ProhibitLogin    bool

	// Avatar
	Avatar          string `gorm:"type:VARCHAR(2048);not null"`
	AvatarEmail     string `gorm:"not null"`
	UseCustomAvatar bool

	NumStars int
	NumRepos int

	// For organization
	Description string
	NumTeams    int
	NumMembers  int
	Teams       []*Team `gorm:"-" json:"-"`
	Members     []*User `gorm:"-" json:"-"`

	DisableNewsFeed bool
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) ToDomain() *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:                 u.ID,
		LowerName:          u.LowerName,
		Name:               u.Name,
		FullName:           u.FullName,
		Email:              u.Email,
		Password:           u.Password,
		LoginSource:        u.LoginSource,
		LoginName:          u.LoginName,
		Type:               domain.UserType(u.Type),
		Location:           u.Location,
		Website:            u.Website,
		Rands:              u.Rands,
		Salt:               u.Salt,
		Created:            time.Unix(u.CreatedUnix, 0),
		Updated:            time.Unix(u.UpdatedUnix, 0),
		LastRepoVisibility: u.LastRepoVisibility,
		MaxRepoCreation:    u.MaxRepoCreation,
		IsActive:           u.IsActive,
		IsAdmin:            u.IsAdmin,
		AllowGitHook:       u.AllowGitHook,
		AllowImportLocal:   u.AllowImportLocal,
		ProhibitLogin:      u.ProhibitLogin,
		Avatar:             u.Avatar,
		AvatarEmail:        u.AvatarEmail,
		UseCustomAvatar:    u.UseCustomAvatar,
		NumStars:           u.NumStars,
		NumRepos:           u.NumRepos,
		Description:        u.Description,
		NumTeams:           u.NumTeams,
		NumMembers:         u.NumMembers,
		DisableNewsFeed:    u.DisableNewsFeed,
	}
}

func FromDomainUser(u *domain.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:                 u.ID,
		LowerName:          u.LowerName,
		Name:               u.Name,
		FullName:           u.FullName,
		Email:              u.Email,
		Password:           u.Password,
		LoginSource:        u.LoginSource,
		LoginName:          u.LoginName,
		Type:               UserType(u.Type),
		Location:           u.Location,
		Website:            u.Website,
		Rands:              u.Rands,
		Salt:               u.Salt,
		CreatedUnix:        u.Created.Unix(),
		UpdatedUnix:        u.Updated.Unix(),
		LastRepoVisibility: u.LastRepoVisibility,
		MaxRepoCreation:    u.MaxRepoCreation,
		IsActive:           u.IsActive,
		IsAdmin:            u.IsAdmin,
		AllowGitHook:       u.AllowGitHook,
		AllowImportLocal:   u.AllowImportLocal,
		ProhibitLogin:      u.ProhibitLogin,
		Avatar:             u.Avatar,
		AvatarEmail:        u.AvatarEmail,
		UseCustomAvatar:    u.UseCustomAvatar,
		NumStars:           u.NumStars,
		NumRepos:           u.NumRepos,
		Description:        u.Description,
		NumTeams:           u.NumTeams,
		NumMembers:         u.NumMembers,
		DisableNewsFeed:    u.DisableNewsFeed,
	}
}

// BeforeCreate implements the GORM create hook.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedUnix == 0 {
		u.CreatedUnix = tx.NowFunc().Unix()
		u.UpdatedUnix = u.CreatedUnix
	}
	return nil
}

// AfterFind implements the GORM query hook.
func (u *User) AfterFind(_ *gorm.DB) error {
	u.FullName = markup.Sanitize(u.FullName)
	u.Created = time.Unix(u.CreatedUnix, 0).Local()
	u.Updated = time.Unix(u.UpdatedUnix, 0).Local()
	return nil
}

// IsLocal returns true if the User is created as local account.
func (u *User) IsLocal() bool {
	return u.LoginSource <= 0
}

// IsOrganization returns true if the User is an organization.
func (u *User) IsOrganization() bool {
	return u.Type == UserTypeOrganization || u.Type == UserTypeOrganizationPrivate
}

func (u *User) IsPrivateOrganization() bool {
	return u.Type == UserTypeOrganizationPrivate
}

// APIFormat returns the API format of a User.
func (u *User) APIFormat() *apiv2types.UserResponse {
	return &apiv2types.UserResponse{
		ID:              u.ID,
		Username:        u.Name,
		FullName:        u.FullName,
		Email:           u.Email,
		AvatarURL:       u.AvatarURL(),
		IsAdmin:         u.IsAdmin,
		DisableNewsFeed: u.DisableNewsFeed,
		CreatedAt:       time.Unix(u.CreatedUnix, 0),
	}
}

func (u *User) CanCreateRepo() bool {
	if u.MaxRepoCreation <= -1 {
		return conf.Repository.MaxCreationLimit <= -1 || u.NumRepos < conf.Repository.MaxCreationLimit
	}
	return u.NumRepos < u.MaxRepoCreation
}

func (u *User) MaxNumRepos() int {
	if u.MaxRepoCreation <= -1 {
		return conf.Repository.MaxCreationLimit
	}
	return u.MaxRepoCreation
}

// CanCreateOrganization returns true if User can create organizations.
func (u *User) CanCreateOrganization() bool {
	return !conf.Admin.DisableRegularOrgCreation || u.IsAdmin
}

// CanEditGitHook returns true if User can edit Git hooks.
func (u *User) CanEditGitHook() bool {
	return u.IsAdmin || u.AllowGitHook
}

// CanImportLocal returns true if User can migrate repositories by local path.
func (u *User) CanImportLocal() bool {
	return conf.Repository.EnableLocalPathMigration && (u.IsAdmin || u.AllowImportLocal)
}

// DisplayName returns the full name of the User if it's not empty, returns the
// Username otherwise.
func (u *User) DisplayName() string {
	if len(u.FullName) > 0 {
		return u.FullName
	}
	return u.Name
}

// HomeURLPath returns the URL path to the User or organization home page.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User` and move this to the "userx"
// package.
func (u *User) HomeURLPath() string {
	return conf.Server.Subpath + "/" + u.Name
}

// HTMLURL returns the full URL to the User or organization home page.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User` and move this to the "userx"
// package.
func (u *User) HTMLURL() string {
	return conf.Server.ExternalURL + u.Name
}

// AvatarURLPath returns the URL path to the User or organization avatar.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User` and move this to the "userx"
// package.
func (u *User) AvatarURLPath() string {
	defaultURLPath := conf.UserDefaultAvatarURLPath()
	if u.ID <= 0 {
		return defaultURLPath
	}

	if !osx.IsFile(userx.CustomAvatarPath(u.ID)) {
		return ""
	}
	return fmt.Sprintf("%s/%s/%d", conf.Server.Subpath, conf.UsersAvatarPathPrefix, u.ID)
}

// AvatarURL returns the full URL to the User or organization avatar.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User` and move this to the "userx"
// package.
func (u *User) AvatarURL() string {
	return u.AvatarURLPath()
}

// IsUserOrgOwner returns true if the User is in the owner team of the given
// organization.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User`.
func (u *User) IsUserOrgOwner(orgID int64) bool {
	return IsOrganizationOwner(orgID, u.ID)
}

// IsPublicMember returns true if the User has public membership of the given
// organization.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User`.
func (u *User) IsPublicMember(orgID int64) bool {
	return IsPublicMembership(orgID, u.ID)
}

// GetOrganizationCount returns the count of organization membership that the
// User has.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User`.
func (u *User) GetOrganizationCount() (int64, error) {
	return Handle.Organizations().CountByUser(context.Background(), u.ID)
}

// ShortName truncates and returns the Username at most in given length.
//
// TODO(unknwon): This is also used in templates, which should be fixed by
// having a dedicated type `template.User`.
func (u *User) ShortName(length int) string {
	return strx.Ellipsis(u.Name, length)
}

// NewGhostUser creates and returns a fake User for people who has deleted their
// accounts.
//
// TODO: Once migrated to unknwon.dev/i18n, pass in the `i18n.Locale` to
// translate the text to local language.
func NewGhostUser() *User {
	return &User{
		ID:        -1,
		Name:      "Ghost",
		LowerName: "ghost",
	}
}

var (
	reservedUsernames = map[string]struct{}{
		"-":        {},
		"explore":  {},
		"create":   {},
		"assets":   {},
		"css":      {},
		"img":      {},
		"js":       {},
		"less":     {},
		"plugins":  {},
		"debug":    {},
		"raw":      {},
		"install":  {},
		"api":      {},
		"avatar":   {},
		"user":     {},
		"org":      {},
		"help":     {},
		"stars":    {},
		"issues":   {},
		"pulls":    {},
		"commits":  {},
		"repo":     {},
		"template": {},
		"admin":    {},
		"new":      {},
		".":        {},
		"..":       {},
	}
	reservedUsernamePatterns = []string{"*.keys"}
)

type ErrNameNotAllowed struct {
	args errx.Args
}

// IsErrNameNotAllowed returns true if the underlying error has the type
// ErrNameNotAllowed.
func IsErrNameNotAllowed(err error) bool {
	var target ErrNameNotAllowed
	return errors.As(err, &target)
}

func (err ErrNameNotAllowed) Value() string {
	val, ok := err.args["name"].(string)
	if ok {
		return val
	}

	val, ok = err.args["pattern"].(string)
	if ok {
		return val
	}

	return "<value not found>"
}

func (err ErrNameNotAllowed) Error() string {
	return fmt.Sprintf("name is not allowed: %v", err.args)
}

// isNameAllowed checks if the name is reserved or pattern of the name is not
// allowed based on given reserved names and patterns. Names are exact match,
// patterns can be prefix or suffix match with the wildcard ("*").
func isNameAllowed(names map[string]struct{}, patterns []string, name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrNameNotAllowed{
			args: errx.Args{
				"reason": "empty name",
			},
		}
	}

	if _, ok := names[name]; ok {
		return ErrNameNotAllowed{
			args: errx.Args{
				"reason": "reserved",
				"name":   name,
			},
		}
	}

	for _, pattern := range patterns {
		if pattern[0] == '*' && strings.HasSuffix(name, pattern[1:]) ||
			(pattern[len(pattern)-1] == '*' && strings.HasPrefix(name, pattern[:len(pattern)-1])) {
			return ErrNameNotAllowed{
				args: errx.Args{
					"reason":  "reserved",
					"pattern": pattern,
				},
			}
		}
	}

	return nil
}

// isUsernameAllowed returns ErrNameNotAllowed if the given name or pattern of
// the name is not allowed as a Username.
func isUsernameAllowed(name string) error {
	return isNameAllowed(reservedUsernames, reservedUsernamePatterns, name)
}

// EmailAddress is an email address of a User.
type EmailAddress struct {
	ID          int64  `gorm:"primaryKey"`
	UserID      int64  `gorm:"column:uid;index;uniqueIndex:email_address_User_email_unique;not null"`
	Email       string `gorm:"uniqueIndex:email_address_User_email_unique;not null;size:254"`
	IsActivated bool   `gorm:"not null;default:FALSE"`
	IsPrimary   bool   `gorm:"-" json:"-"`
}

func (e *EmailAddress) ToDomain() *domain.EmailAddress {
	if e == nil {
		return nil
	}
	return &domain.EmailAddress{
		ID:          e.ID,
		UserID:      e.UserID,
		Email:       e.Email,
		IsActivated: e.IsActivated,
		IsPrimary:   e.IsPrimary,
	}
}
