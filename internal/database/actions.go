package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/lazyregexp"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/repox"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/nef734/ggit/internal/strx"
	"github.com/nef734/ggit/internal/tool"
	"gorm.io/gorm"
)

// ActionsStore is the storage layer for actions.
type ActionsStore struct {
	db *gorm.DB
}

func newActionsStore(db *gorm.DB) *ActionsStore {
	return &ActionsStore{db: db}
}

func (s *ActionsStore) listByOrganization(ctx context.Context, orgID, actorID, afterID int64) *gorm.DB {
	/*
		Equivalent SQL for PostgreSQL:

		SELECT * FROM "action"
		WHERE
			user_id = @userID
		AND (@skipAfter OR id < @afterID)
		AND repo_id IN (
			SELECT repository.id FROM "repository"
			JOIN team_repo ON repository.id = team_repo.repo_id
			WHERE team_repo.team_id IN (
					SELECT team_id FROM "team_user"
					WHERE
						team_user.org_id = @orgID AND uid = @actorID)
					OR  (repository.is_private = FALSE AND repository.is_unlisted = FALSE)
			)
		ORDER BY id DESC
		LIMIT @limit
	*/
	return s.db.WithContext(ctx).
		Where("user_id = ?", orgID).
		Where(s.db.
			// Not apply when afterID is not given
			Where("?", afterID <= 0).
			Or("id < ?", afterID),
		).
		Where("repo_id IN (?)", s.db.
			Select("repository.id").
			Table("repository").
			Joins("JOIN team_repo ON repository.id = team_repo.repo_id").
			Where("team_repo.team_id IN (?)", s.db.
				Select("team_id").
				Table("team_user").
				Where("team_user.org_id = ? AND uid = ?", orgID, actorID),
			).
			Or("repository.is_private = ? AND repository.is_unlisted = ?", false, false),
		).
		Limit(conf.UI.User.NewsFeedPagingNum).
		Order("id DESC")
}

// ListByOrganization returns actions of the organization viewable by the actor.
// Results are paginated if `afterID` is given.
func (s *ActionsStore) ListByOrganization(ctx context.Context, orgID, actorID, afterID int64) ([]*domain.Action, error) {
	models := make([]*Action, 0, conf.UI.User.NewsFeedPagingNum)
	err := s.listByOrganization(ctx, orgID, actorID, afterID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	actions := make([]*domain.Action, len(models))
	for i, m := range models {
		actions[i] = m.ToDomain()
	}
	return actions, nil
}

func (s *ActionsStore) listByUser(ctx context.Context, userID, actorID, afterID int64, isProfile bool) *gorm.DB {
	/*
		Equivalent SQL for PostgreSQL:

		SELECT * FROM "action"
		WHERE
			user_id = @userID
		AND (@skipAfter OR id < @afterID)
		AND (@includePrivate OR (is_private = FALSE AND act_user_id = @actorID))
		ORDER BY id DESC
		LIMIT @limit
	*/
	return s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where(s.db.
			// Not apply when afterID is not given
			Where("?", afterID <= 0).
			Or("id < ?", afterID),
		).
		Where(s.db.
			// Not apply when in not profile page or the user is viewing own profile
			Where("?", !isProfile || actorID == userID).
			Or("is_private = ? AND act_user_id = ?", false, userID),
		).
		Limit(conf.UI.User.NewsFeedPagingNum).
		Order("id DESC")
}

// ListByUser returns actions of the user viewable by the actor. Results are
// paginated if `afterID` is given. The `isProfile` indicates whether repository
// permissions should be considered.
func (s *ActionsStore) ListByUser(ctx context.Context, userID, actorID, afterID int64, isProfile bool) ([]*domain.Action, error) {
	models := make([]*Action, 0, conf.UI.User.NewsFeedPagingNum)
	err := s.listByUser(ctx, userID, actorID, afterID, isProfile).Find(&models).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Action, len(models))
	for i, m := range models {
		res[i] = m.ToDomain()
	}
	return res, nil
}

// ListAuditLog returns all actions for the audit log. It only returns one copy per action
// by filtering where user_id == act_user_id.
func (s *ActionsStore) ListAuditLog(ctx context.Context, page int) ([]*domain.Action, error) {
	if page <= 0 {
		page = 1
	}
	actions := make([]*Action, 0, conf.UI.Admin.NoticePagingNum)
	err := s.db.WithContext(ctx).
		Where("user_id = act_user_id").
		Limit(conf.UI.Admin.NoticePagingNum).
		Offset((page - 1) * conf.UI.Admin.NoticePagingNum).
		Order("id DESC").
		Find(&actions).Error
	if err != nil {
		return nil, err
	}

	dactions := make([]*domain.Action, len(actions))
	for i := range actions {
		dactions[i] = actions[i].ToDomain()
	}
	return dactions, nil
}

// notifyWatchers creates rows in action table for watchers who are able to see the action.
func (s *ActionsStore) notifyWatchers(ctx context.Context, act *Action) error {
	watches, err := NewRepositoriesStore(s.db).ListWatches(ctx, act.RepoID)
	if err != nil {
		return fmt.Errorf("list watches"+": %w", err)
	}

	// Clone returns a deep copy of the action with UserID assigned
	clone := func(userID int64) *Action {
		tmp := *act
		tmp.UserID = userID
		return &tmp
	}

	// Plus one for the actor
	actions := make([]*Action, 0, len(watches)+1)
	actions = append(actions, clone(act.ActUserID))

	for _, watch := range watches {
		if act.ActUserID == watch.ID {
			continue
		}
		actions = append(actions, clone(watch.ID))
	}

	return s.db.Create(actions).Error
}

// UpdateConfig creates an action for updating system configuration.
func (s *ActionsStore) UpdateConfig(ctx context.Context, doer *domain.User, section, content string) error {
	return s.db.Create(&Action{
		UserID:      doer.ID,
		ActUserID:   doer.ID,
		ActUserName: doer.Name,
		OpType:      ActionUpdateConfig,
		Content:     fmt.Sprintf("%s|%s", section, content),
	}).Error
}

// NewRepo creates an action for creating a new repository. The action type
// could be ActionCreateRepo or ActionForkRepo based on whether the repository
// is a fork.
func (s *ActionsStore) NewRepo(ctx context.Context, doer, owner *User, repo *Repository) error {
	opType := ActionCreateRepo
	if repo.IsFork {
		opType = ActionForkRepo
	}

	return s.notifyWatchers(ctx,
		&Action{
			ActUserID:    doer.ID,
			ActUserName:  doer.Name,
			OpType:       opType,
			RepoID:       repo.ID,
			RepoUserName: owner.Name,
			RepoName:     repo.Name,
			IsPrivate:    repo.IsPrivate || repo.IsUnlisted,
		},
	)
}

// MergePullRequest creates an action for merging a pull request.
func (s *ActionsStore) MergePullRequest(ctx context.Context, doer, owner *domain.User, repo *domain.Repository, pull *domain.Issue) error {
	return s.notifyWatchers(ctx,
		&Action{
			ActUserID:    doer.ID,
			ActUserName:  doer.Name,
			OpType:       ActionMergePullRequest,
			Content:      fmt.Sprintf("%d|%s", pull.Index, pull.Title),
			RepoID:       repo.ID,
			RepoUserName: owner.Name,
			RepoName:     repo.Name,
			IsPrivate:    repo.IsPrivate,
		},
	)
}

// TransferRepo creates an action for transferring a repository to a new owner.
func (s *ActionsStore) TransferRepo(ctx context.Context, doer, oldOwner, newOwner *domain.User, repo *domain.Repository) error {
	return s.notifyWatchers(ctx,
		&Action{
			ActUserID:    doer.ID,
			ActUserName:  doer.Name,
			OpType:       ActionTransferRepo,
			RepoID:       repo.ID,
			RepoUserName: newOwner.Name,
			RepoName:     repo.Name,
			IsPrivate:    repo.IsPrivate,
			Content:      oldOwner.Name + "/" + repo.Name,
		},
	)
}

// RenameRepo creates an action for renaming a repository.
func (s *ActionsStore) RenameRepo(ctx context.Context, doer, owner *domain.User, oldRepoName string, repo *domain.Repository) error {
	return s.notifyWatchers(ctx,
		&Action{
			ActUserID:    doer.ID,
			ActUserName:  doer.Name,
			OpType:       ActionRenameRepo,
			RepoID:       repo.ID,
			RepoUserName: owner.Name,
			RepoName:     repo.Name,
			IsPrivate:    repo.IsPrivate,
			Content:      oldRepoName,
		},
	)
}

var (
	// Same as GitHub, see https://docs.github.com/en/free-pro-team@latest/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue
	issueCloseKeywords  = []string{"close", "closes", "closed", "fix", "fixes", "fixed", "resolve", "resolves", "resolved"}
	issueReopenKeywords = []string{"reopen", "reopens", "reopened"}

	issueCloseKeywordsPattern  = lazyregexp.New(assembleKeywordsPattern(issueCloseKeywords))
	issueReopenKeywordsPattern = lazyregexp.New(assembleKeywordsPattern(issueReopenKeywords))
	issueReferencePattern      = lazyregexp.New(`(?i)(?:)(^| )\S*#\d+`)
)

func assembleKeywordsPattern(words []string) string {
	return fmt.Sprintf(`(?i)(?:%s) \S+`, strings.Join(words, "|"))
}

// updateCommitReferencesToIssues checks if issues are manipulated by commit message.
func updateCommitReferencesToIssues(doer *User, repo *Repository, commits []*PushCommit) error {
	trimRightNonDigits := func(c rune) bool {
		return !unicode.IsDigit(c)
	}

	// Commits are appended in the reverse order.
	for i := len(commits) - 1; i >= 0; i-- {
		c := commits[i]

		refMarked := make(map[int64]bool)
		for _, ref := range issueReferencePattern.FindAllString(c.Message, -1) {
			ref = strings.TrimSpace(ref)
			ref = strings.TrimRightFunc(ref, trimRightNonDigits)

			if ref == "" {
				continue
			}

			// Add repo name if missing
			if ref[0] == '#' {
				ref = fmt.Sprintf("%s%s", repo.FullName(), ref)
			} else if !strings.Contains(ref, "/") {
				// FIXME: We don't support User#ID syntax yet
				continue
			}

			issue, err := GetIssueByRef(ref)
			if err != nil {
				if IsErrIssueNotExist(err) {
					continue
				}
				return err
			}

			if refMarked[issue.ID] {
				continue
			}
			refMarked[issue.ID] = true

			msgLines := strings.Split(c.Message, "\n")
			shortMsg := msgLines[0]
			if len(msgLines) > 2 {
				shortMsg += "..."
			}
			message := fmt.Sprintf(`<a href="%s/commit/%s">%s</a>`, repo.Link(), c.Sha1, shortMsg)
			if err = CreateRefComment(doer, repo, issue, message, c.Sha1); err != nil {
				return err
			}
		}

		refMarked = make(map[int64]bool)
		// FIXME: Can merge this and the next for loop to a common function.
		for _, ref := range issueCloseKeywordsPattern.FindAllString(c.Message, -1) {
			ref = ref[strings.IndexByte(ref, byte(' '))+1:]
			ref = strings.TrimRightFunc(ref, trimRightNonDigits)

			if ref == "" {
				continue
			}

			// Add repo name if missing
			if ref[0] == '#' {
				ref = fmt.Sprintf("%s%s", repo.FullName(), ref)
			} else if !strings.Contains(ref, "/") {
				// FIXME: We don't support User#ID syntax yet
				continue
			}

			issue, err := GetIssueByRef(ref)
			if err != nil {
				if IsErrIssueNotExist(err) {
					continue
				}
				return err
			}

			if refMarked[issue.ID] {
				continue
			}
			refMarked[issue.ID] = true

			if issue.RepoID != repo.ID || issue.IsClosed {
				continue
			}

			if err = issue.ChangeStatus(context.TODO(), doer, repo, true); err != nil {
				return err
			}
		}

		// It is conflict to have close and reopen at same time, so refsMarkd doesn't need to reinit here.
		for _, ref := range issueReopenKeywordsPattern.FindAllString(c.Message, -1) {
			ref = ref[strings.IndexByte(ref, byte(' '))+1:]
			ref = strings.TrimRightFunc(ref, trimRightNonDigits)

			if ref == "" {
				continue
			}

			// Add repo name if missing
			if ref[0] == '#' {
				ref = fmt.Sprintf("%s%s", repo.FullName(), ref)
			} else if !strings.Contains(ref, "/") {
				// We don't support User#ID syntax yet
				// return ErrNotImplemented
				continue
			}

			issue, err := GetIssueByRef(ref)
			if err != nil {
				if IsErrIssueNotExist(err) {
					continue
				}
				return err
			}

			if refMarked[issue.ID] {
				continue
			}
			refMarked[issue.ID] = true

			if issue.RepoID != repo.ID || !issue.IsClosed {
				continue
			}

			if err = issue.ChangeStatus(context.TODO(), doer, repo, false); err != nil {
				return err
			}
		}
	}
	return nil
}

type CommitRepoOptions struct {
	Owner       *User
	Repo        *Repository
	PusherName  string
	RefFullName string
	OldCommitID string
	NewCommitID string
	Commits     *PushCommits
}

// CommitRepo creates actions for pushing commits to the repository. An action
// with the type ActionDeleteBranch is created if the push deletes a branch; an
// action with the type ActionCommitRepo is created for a regular push. If the
// regular push also creates a new branch, then another action with type
// ActionCreateBranch is created.
func (s *ActionsStore) CommitRepo(ctx context.Context, opts any) error {
	var options CommitRepoOptions
	if o, ok := opts.(CommitRepoOptions); ok {
		options = o
	} else if o, ok := opts.(*CommitRepoOptions); ok {
		options = *o
	} else {
		return errors.New("invalid options type")
	}
	optionsValue := options

	err := NewRepositoriesStore(s.db).Touch(ctx, optionsValue.Repo.ID)
	if err != nil {
		return fmt.Errorf("touch repository"+": %w", err)
	}

	pusher, err := NewUsersStore(s.db).getByUsername(ctx, optionsValue.PusherName)
	if err != nil {
		return fmt.Errorf("get pusher [name: %s]: %w", optionsValue.PusherName, err)
	}

	isNewRef := optionsValue.OldCommitID == git.EmptyID
	isDelRef := optionsValue.NewCommitID == git.EmptyID

	// If not the first commit, set the compare URL.
	if !isNewRef && !isDelRef {
		optionsValue.Commits.CompareURL = repox.CompareCommitsPath(optionsValue.Owner.Name, optionsValue.Repo.Name, optionsValue.OldCommitID, optionsValue.NewCommitID)
	}

	refName := git.RefShortName(optionsValue.RefFullName)
	action := &Action{
		ActUserID:    pusher.ID,
		ActUserName:  pusher.Name,
		RepoID:       optionsValue.Repo.ID,
		RepoUserName: optionsValue.Owner.Name,
		RepoName:     optionsValue.Repo.Name,
		RefName:      refName,
		IsPrivate:    optionsValue.Repo.IsPrivate || optionsValue.Repo.IsUnlisted,
	}

	apiRepo := optionsValue.Repo.APIFormat(optionsValue.Owner)
	apiPusher := pusher.APIFormat()
	if isDelRef {
		err = PrepareWebhooksWithDB(
			s.db,
			optionsValue.Repo,
			HookEventTypeDelete,
			&apiv2types.WebhookDeletePayload{
				Ref:        refName,
				RefType:    "branch",
				PusherType: apiv2types.WebhookPusherTypeUser,
				Repo:       apiRepo,
				Sender:     apiPusher,
			},
		)
		if err != nil {
			return fmt.Errorf("prepare webhooks for delete branch"+": %w", err)
		}

		action.OpType = ActionDeleteBranch
		err = s.notifyWatchers(ctx, action)
		if err != nil {
			return fmt.Errorf("notify watchers"+": %w", err)
		}

		// Delete branch doesn't have anything to push or compare
		return nil
	}

	// Only update issues via commits when internal issue tracker is enabled
	if optionsValue.Repo.EnableIssues && !optionsValue.Repo.EnableExternalTracker {
		if err = updateCommitReferencesToIssues(pusher, optionsValue.Repo, optionsValue.Commits.Commits); err != nil {
			log.Error("update commit references to issues: %v", err)
		}
	}

	if conf.UI.FeedMaxCommitNum > 0 && len(optionsValue.Commits.Commits) > conf.UI.FeedMaxCommitNum {
		optionsValue.Commits.Commits = optionsValue.Commits.Commits[:conf.UI.FeedMaxCommitNum]
	}

	data, err := json.Marshal(optionsValue.Commits)
	if err != nil {
		return fmt.Errorf("marshal JSON"+": %w", err)
	}
	action.Content = string(data)

	var compareURL string
	if isNewRef {
		err = PrepareWebhooksWithDB(
			s.db,
			optionsValue.Repo,
			HookEventTypeCreate,
			&apiv2types.WebhookCreatePayload{
				Ref:           refName,
				RefType:       "branch",
				DefaultBranch: optionsValue.Repo.DefaultBranch,
				Repo:          apiRepo,
				Sender:        apiPusher,
			},
		)
		if err != nil {
			return fmt.Errorf("prepare webhooks for new branch"+": %w", err)
		}

		action.OpType = ActionCreateBranch
		err = s.notifyWatchers(ctx, action)
		if err != nil {
			return fmt.Errorf("notify watchers"+": %w", err)
		}
	} else {
		compareURL = conf.Server.ExternalURL + optionsValue.Commits.CompareURL
	}

	commits, err := optionsValue.Commits.APIFormat(ctx,
		NewUsersStore(s.db),
		repox.RepositoryPath(optionsValue.Owner.Name, optionsValue.Repo.Name),
		repox.HTMLURL(optionsValue.Owner.Name, optionsValue.Repo.Name),
	)
	if err != nil {
		return fmt.Errorf("convert commits to API format"+": %w", err)
	}

	err = PrepareWebhooks(
		optionsValue.Repo,
		HookEventTypePush,
		&apiv2types.WebhookPushPayload{
			Ref:        optionsValue.RefFullName,
			Before:     optionsValue.OldCommitID,
			After:      optionsValue.NewCommitID,
			CompareURL: compareURL,
			Commits:    commits,
			Repo:       apiRepo,
			Pusher:     apiPusher,
			Sender:     apiPusher,
		},
	)
	if err != nil {
		return fmt.Errorf("prepare webhooks for new commit"+": %w", err)
	}

	action.OpType = ActionCommitRepo
	err = s.notifyWatchers(ctx, action)
	if err != nil {
		return fmt.Errorf("notify watchers"+": %w", err)
	}
	return nil
}

type PushTagOptions struct {
	Owner       *User
	Repo        *Repository
	PusherName  string
	RefFullName string
	NewCommitID string
}

// PushTag creates an action for pushing tags to the repository. An action with
// the type ActionDeleteTag is created if the push deletes a tag. Otherwise, an
// action with the type ActionPushTag is created for a regular push.
func (s *ActionsStore) PushTag(ctx context.Context, opts any) error {
	var options PushTagOptions
	if o, ok := opts.(PushTagOptions); ok {
		options = o
	} else {
		return errors.New("invalid options type")
	}

	err := NewRepositoriesStore(s.db).Touch(ctx, options.Repo.ID)
	if err != nil {
		return fmt.Errorf("touch repository"+": %w", err)
	}

	pusher, err := NewUsersStore(s.db).getByUsername(ctx, options.PusherName)
	if err != nil {
		return fmt.Errorf("get pusher [name: %s]: %w", options.PusherName, err)
	}

	refName := git.RefShortName(options.RefFullName)
	action := &Action{
		ActUserID:    pusher.ID,
		ActUserName:  pusher.Name,
		RepoID:       options.Repo.ID,
		RepoUserName: options.Owner.Name,
		RepoName:     options.Repo.Name,
		RefName:      refName,
		IsPrivate:    options.Repo.IsPrivate || options.Repo.IsUnlisted,
	}

	apiRepo := options.Repo.APIFormat(options.Owner)
	apiPusher := pusher.APIFormat()
	if options.NewCommitID == git.EmptyID {
		err = PrepareWebhooks(
			options.Repo,
			HookEventTypeDelete,
			&apiv2types.WebhookDeletePayload{
				Ref:        refName,
				RefType:    "tag",
				PusherType: apiv2types.WebhookPusherTypeUser,
				Repo:       apiRepo,
				Sender:     apiPusher,
			},
		)
		if err != nil {
			return fmt.Errorf("prepare webhooks for delete tag"+": %w", err)
		}

		action.OpType = ActionDeleteTag
		err = s.notifyWatchers(ctx, action)
		if err != nil {
			return fmt.Errorf("notify watchers"+": %w", err)
		}
		return nil
	}

	err = PrepareWebhooksWithDB(
		s.db,
		options.Repo,
		HookEventTypeCreate,
		&apiv2types.WebhookCreatePayload{
			Ref:           refName,
			RefType:       "tag",
			Sha:           options.NewCommitID,
			DefaultBranch: options.Repo.DefaultBranch,
			Repo:          apiRepo,
			Sender:        apiPusher,
		},
	)
	if err != nil {
		return fmt.Errorf("prepare webhooks for new tag: %w", err)
	}

	action.OpType = ActionPushTag
	err = s.notifyWatchers(ctx, action)
	if err != nil {
		return fmt.Errorf("notify watchers"+": %w", err)
	}
	return nil
}

// ActionType is the type of action.
type ActionType int

// ⚠️ WARNING: Only append to the end of list to maintain backward compatibility.
const (
	ActionCreateRepo        ActionType = iota + 1 // 1
	ActionRenameRepo                              // 2
	ActionStarRepo                                // 3
	ActionWatchRepo                               // 4
	ActionCommitRepo                              // 5
	ActionCreateIssue                             // 6
	ActionCreatePullRequest                       // 7
	ActionTransferRepo                            // 8
	ActionPushTag                                 // 9
	ActionCommentIssue                            // 10
	ActionMergePullRequest                        // 11
	ActionCloseIssue                              // 12
	ActionReopenIssue                             // 13
	ActionClosePullRequest                        // 14
	ActionReopenPullRequest                       // 15
	ActionCreateBranch                            // 16
	ActionDeleteBranch                            // 17
	ActionDeleteTag                               // 18
	ActionForkRepo                                // 19
	ActionUpdateConfig                            // 20
)

// Action is a user operation to a repository. It implements template.Actioner
// interface to be able to use it in template rendering.
type Action struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	UserID       int64      `gorm:"index" json:"user_id"` // Receiver user ID
	OpType       ActionType `json:"op_type"`
	ActUserID    int64      `json:"act_user_id"`   // Doer user ID
	ActUserName  string     `json:"act_user_name"` // Doer user name
	ActAvatar    string     `gorm:"-" json:"act_avatar"`
	RepoID       int64      `gorm:"index" json:"repo_id"`
	RepoUserName string     `json:"repo_user_name"`
	RepoName     string     `json:"repo_name"`
	RefName      string     `json:"ref_name"`
	IsPrivate    bool       `gorm:"not null;default:FALSE" json:"is_private"`
	Content      string     `json:"content"`

	Created     time.Time `gorm:"-" json:"created_at"`
	CreatedUnix int64     `json:"created_unix"`
}

func (a *Action) ToDomain() *domain.Action {
	if a == nil {
		return nil
	}
	return &domain.Action{
		ID:           a.ID,
		UserID:       a.UserID,
		OpType:       int(a.OpType),
		ActUserID:    a.ActUserID,
		ActUserName:  a.ActUserName,
		RepoID:       a.RepoID,
		RepoUserName: a.RepoUserName,
		RepoName:     a.RepoName,
		RefName:      a.RefName,
		IsPrivate:    a.IsPrivate,
		Content:      a.Content,
		Created:      time.Unix(a.CreatedUnix, 0),
	}
}

func FromDomainAction(a *domain.Action) *Action {
	if a == nil {
		return nil
	}
	return &Action{
		ID:           a.ID,
		UserID:       a.UserID,
		OpType:       ActionType(a.OpType),
		ActUserID:    a.ActUserID,
		ActUserName:  a.ActUserName,
		RepoID:       a.RepoID,
		RepoUserName: a.RepoUserName,
		RepoName:     a.RepoName,
		RefName:      a.RefName,
		IsPrivate:    a.IsPrivate,
		Content:      a.Content,
		CreatedUnix:  a.Created.Unix(),
	}
}

// BeforeCreate implements the GORM create hook.
func (a *Action) BeforeCreate(tx *gorm.DB) error {
	if a.CreatedUnix <= 0 {
		a.CreatedUnix = tx.NowFunc().Unix()
	}
	return nil
}

// AfterFind implements the GORM query hook.
func (a *Action) AfterFind(_ *gorm.DB) error {
	a.Created = time.Unix(a.CreatedUnix, 0).Local()
	return nil
}

func (a *Action) GetOpType() int {
	return int(a.OpType)
}

func (a *Action) GetActUserName() string {
	return a.ActUserName
}

func (a *Action) ShortActUserName() string {
	return strx.Ellipsis(a.ActUserName, 20)
}

func (a *Action) GetRepoUserName() string {
	return a.RepoUserName
}

func (a *Action) ShortRepoUserName() string {
	return strx.Ellipsis(a.RepoUserName, 20)
}

func (a *Action) GetRepoName() string {
	return a.RepoName
}

func (a *Action) ShortRepoName() string {
	return strx.Ellipsis(a.RepoName, 33)
}

func (a *Action) GetRepoPath() string {
	return path.Join(a.RepoUserName, a.RepoName)
}

func (a *Action) ShortRepoPath() string {
	return path.Join(a.ShortRepoUserName(), a.ShortRepoName())
}

func (a *Action) GetRepoLink() string {
	if conf.Server.Subpath != "" {
		return path.Join(conf.Server.Subpath, a.GetRepoPath())
	}
	return "/" + a.GetRepoPath()
}

func (a *Action) GetBranch() string {
	return a.RefName
}

func (a *Action) GetContent() string {
	return a.Content
}

func (a *Action) GetCreate() time.Time {
	return a.Created
}

func (a *Action) GetIssueInfos() []string {
	return strings.SplitN(a.Content, "|", 2)
}

func (a *Action) GetIssueTitle() string {
	index, _ := strconv.ParseInt(a.GetIssueInfos()[0], 10, 64)
	issue, err := GetIssueByIndex(a.RepoID, index)
	if err != nil {
		log.Error("Failed to get issue title [repo_id: %d, index: %d]: %v", a.RepoID, index, err)
		return "error getting issue"
	}
	return issue.Title
}

func (a *Action) GetIssueContent() string {
	index, _ := strconv.ParseInt(a.GetIssueInfos()[0], 10, 64)
	issue, err := GetIssueByIndex(a.RepoID, index)
	if err != nil {
		log.Error("Failed to get issue content [repo_id: %d, index: %d]: %v", a.RepoID, index, err)
		return "error getting issue"
	}
	return issue.Content
}

// PushCommit contains information of a pushed commit.
type PushCommit struct {
	Sha1           string
	Message        string
	AuthorEmail    string
	AuthorName     string
	CommitterEmail string
	CommitterName  string
	Timestamp      time.Time
}

// PushCommits is a list of pushed commits.
type PushCommits struct {
	Len        int
	Commits    []*PushCommit
	CompareURL string

	avatars map[string]string
}

// NewPushCommits returns a new PushCommits.
func NewPushCommits() *PushCommits {
	return &PushCommits{
		avatars: make(map[string]string),
	}
}

func (pcs *PushCommits) APIFormat(ctx context.Context, usersStore domain.UsersStore, repoPath, repoURL string) ([]*apiv2types.WebhookPayloadCommit, error) {
	commits := make([]*apiv2types.WebhookPayloadCommit, len(pcs.Commits))
	for i := range pcs.Commits {
		commits[i] = &apiv2types.WebhookPayloadCommit{
			ID:      pcs.Commits[i].Sha1,
			Message: pcs.Commits[i].Message,
			URL:     repoURL + "/commit/" + pcs.Commits[i].Sha1,
			Author: &apiv2types.WebhookPayloadUser{
				Name:  pcs.Commits[i].AuthorName,
				Email: pcs.Commits[i].AuthorEmail,
			},
			Committer: &apiv2types.WebhookPayloadUser{
				Name:  pcs.Commits[i].CommitterName,
				Email: pcs.Commits[i].CommitterEmail,
			},
			Timestamp: pcs.Commits[i].Timestamp,
		}
	}
	return commits, nil
}

// AvatarLink tries to match user in database with email in order to show custom
// avatars, and falls back to general avatar link.
//
// FIXME: This method does not belong to PushCommits, should be a pure template
// function.
func (pcs *PushCommits) AvatarLink(email string) string {
	_, ok := pcs.avatars[email]
	if !ok {
		u, err := Handle.Users().GetByEmail(context.Background(), email)
		if err != nil {
			pcs.avatars[email] = tool.AvatarLink(email)
			if !IsErrUserNotExist(err) {
				log.Error("Failed to get user [email: %s]: %v", email, err)
			}
		} else {
			pcs.avatars[email] = u.AvatarURLPath()
		}
	}

	return pcs.avatars[email]
}
