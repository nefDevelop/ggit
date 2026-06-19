package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/markup"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"gorm.io/gorm"
)

// CommentType defines whether a comment is just a simple comment, an action (like close) or a reference.
type CommentType int

const (
	// Plain comment, can be associated with a commit (CommitID > 0) and a line (LineNum > 0)
	CommentTypeComment CommentType = iota
	CommentTypeReopen
	CommentTypeClose

	// References.
	CommentTypeIssueRef
	// Reference from a commit (not part of a pull request)
	CommentTypeCommitRef
	// Reference from a comment
	CommentTypeCommentRef
	// Reference from a pull request
	CommentTypePullRef
)

type CommentTag int

const (
	CommentTagNone CommentTag = iota
	CommentTagPoster
	CommentTagWriter
	CommentTagOwner
)

// Comment represents a comment in commit and issue page.
type Comment struct {
	ID       int64       `gorm:"primaryKey"`
	Type     CommentType `gorm:"not null"`
	PosterID int64       `gorm:"index"`
	Poster   *User       `json:"-" gorm:"-"`
	IssueID  int64       `gorm:"index"`
	Issue    *Issue      `json:"-" gorm:"-"`
	CommitID int64
	Line     int64
	Content  string `gorm:"type:TEXT"`

	Created     time.Time `json:"-" gorm:"-"`
	CreatedUnix int64
	Updated     time.Time `json:"-" gorm:"-"`
	UpdatedUnix int64

	// Reference issue in commit message
	CommitSHA string `gorm:"type:VARCHAR(40)"`

	Attachments []*Attachment `json:"-" gorm:"-"`

	// For view issue page.
	ShowTag CommentTag `json:"-" gorm:"-"`
}

type CommentsStore struct {
	db *gorm.DB
}

func newCommentsStore(db *gorm.DB) *CommentsStore {
	return &CommentsStore{db: db}
}

func (s *CommentsStore) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	c, err := getCommentByID(s.db, id)
	if err != nil {
		return nil, err
	}
	return c.ToDomain(), nil
}

func (s *CommentsStore) ListByIssueID(ctx context.Context, issueID int64) ([]*domain.Comment, error) {
	models, err := getCommentsByIssueID(s.db, issueID)
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Comment, len(models))
	for i, m := range models {
		res[i] = m.ToDomain()
	}
	return res, nil
}

func (c *Comment) ToDomain() *domain.Comment {
	if c == nil {
		return nil
	}
	return &domain.Comment{
		ID:        c.ID,
		Type:      int(c.Type),
		PosterID:  c.PosterID,
		Poster:    c.Poster.ToDomain(),
		IssueID:   c.IssueID,
		Issue:     c.Issue.ToDomain(),
		CommitID:  c.CommitID,
		Line:      c.Line,
		Content:   c.Content,
		Created:   c.Created,
		Updated:   c.Updated,
		CommitSHA: c.CommitSHA,
	}
}

func FromDomainComment(c *domain.Comment) *Comment {
	if c == nil {
		return nil
	}
	return &Comment{
		ID:          c.ID,
		Type:        CommentType(c.Type),
		PosterID:    c.PosterID,
		Poster:      FromDomainUser(c.Poster),
		IssueID:     c.IssueID,
		Issue:       FromDomainIssue(c.Issue),
		CommitID:    c.CommitID,
		Line:        c.Line,
		Content:     c.Content,
		CreatedUnix: c.Created.Unix(),
		UpdatedUnix: c.Updated.Unix(),
		CommitSHA:   c.CommitSHA,
	}
}

// BeforeCreate implements the GORM create hook.
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedUnix == 0 {
		c.CreatedUnix = tx.NowFunc().Unix()
		c.UpdatedUnix = c.CreatedUnix
	}
	return nil
}

// BeforeUpdate implements the GORM update hook.
func (c *Comment) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedUnix = tx.NowFunc().Unix()
	return nil
}

// AfterFind implements the GORM query hook.
func (c *Comment) AfterFind(_ *gorm.DB) error {
	c.Created = time.Unix(c.CreatedUnix, 0).Local()
	c.Updated = time.Unix(c.UpdatedUnix, 0).Local()
	return nil
}

func (c *Comment) loadAttributes(tx *gorm.DB) (err error) {
	if c.Poster == nil {
		u, err := Handle.Users().GetByID(context.TODO(), c.PosterID)
		if err != nil {
			if IsErrUserNotExist(err) {
				c.PosterID = -1
				c.Poster = NewGhostUser()
			} else {
				return fmt.Errorf("GetByID.(Poster) [%d]: %v", c.PosterID, err)
			}
		} else {
			c.Poster = FromDomainUser(u)
		}
	}

	if c.Issue == nil {
		c.Issue, err = getRawIssueByID(tx, c.IssueID)
		if err != nil {
			return fmt.Errorf("getIssueByID [%d]: %v", c.IssueID, err)
		}
		if c.Issue.Repo == nil {
			c.Issue.Repo, err = getRepositoryByID(tx, c.Issue.RepoID)
			if err != nil {
				return fmt.Errorf("getRepositoryByID [%d]: %v", c.Issue.RepoID, err)
			}
		}
	}

	if c.Attachments == nil {
		c.Attachments, err = getAttachmentsByCommentID(tx, c.ID)
		if err != nil {
			return fmt.Errorf("getAttachmentsByCommentID [%d]: %v", c.ID, err)
		}
	}

	return nil
}

func (c *Comment) LoadAttributes() error {
	return c.loadAttributes(Handle.db)
}

func (c *Comment) HTMLURL() string {
	return fmt.Sprintf("%s#issuecomment-%d", c.Issue.HTMLURL(), c.ID)
}

// This method assumes following fields have been assigned with valid values:
// Required - Poster, Issue
func (c *Comment) APIFormat() *apiv2types.IssueCommentResponse {
	return &apiv2types.IssueCommentResponse{
		ID:        c.ID,
		HTMLURL:   c.HTMLURL(),
		Poster:    c.Poster.APIFormat(),
		Body:      c.Content,
		CreatedAt: c.Created,
		UpdatedAt: c.Updated,
	}
}

func CommentHashTag(id int64) string {
	return "issuecomment-" + strconv.FormatInt(id, 10)
}

// HashTag returns unique hash tag for comment.
func (c *Comment) HashTag() string {
	return CommentHashTag(c.ID)
}

// EventTag returns unique event hash tag for comment.
func (c *Comment) EventTag() string {
	return "event-" + strconv.FormatInt(c.ID, 10)
}

// mailParticipants sends new comment emails to repository watchers
// and mentioned people.
func (c *Comment) mailParticipants(tx *gorm.DB, opType ActionType, issue *Issue) (err error) {
	mentions := markup.FindAllMentions(c.Content)
	if err = updateIssueMentions(tx, c.IssueID, mentions); err != nil {
		return fmt.Errorf("UpdateIssueMentions [%d]: %v", c.IssueID, err)
	}

	switch opType {
	case ActionCommentIssue:
		issue.Content = c.Content
	case ActionCloseIssue:
		issue.Content = fmt.Sprintf("Closed #%d", issue.Index)
	case ActionReopenIssue:
		issue.Content = fmt.Sprintf("Reopened #%d", issue.Index)
	}
	if err = mailIssueCommentToParticipants(issue, c.Poster, mentions); err != nil {
		log.Error("mailIssueCommentToParticipants: %v", err)
	}

	return nil
}

func createComment(tx *gorm.DB, opts *CreateCommentOptions) (_ *Comment, err error) {
	comment := &Comment{
		Type:      opts.Type,
		PosterID:  opts.Doer.ID,
		Poster:    opts.Doer,
		IssueID:   opts.Issue.ID,
		CommitID:  opts.CommitID,
		CommitSHA: opts.CommitSHA,
		Line:      opts.LineNum,
		Content:   opts.Content,
	}
	if err = tx.Create(comment).Error; err != nil {
		return nil, err
	}

	// Compose comment action, could be plain comment, close or reopen issue/pull request.
	// This object will be used to notify watchers in the end of function.
	act := &Action{
		ActUserID:    opts.Doer.ID,
		ActUserName:  opts.Doer.Name,
		Content:      fmt.Sprintf("%d|%s", opts.Issue.Index, strings.Split(opts.Content, "\n")[0]),
		RepoID:       opts.Repo.ID,
		RepoUserName: opts.Repo.Owner.Name,
		RepoName:     opts.Repo.Name,
		IsPrivate:    opts.Repo.IsPrivate,
	}

	// Check comment type.
	switch opts.Type {
	case CommentTypeComment:
		act.OpType = ActionCommentIssue

		if err = tx.Model(new(Issue)).Where("id = ?", opts.Issue.ID).Update("num_comments", gorm.Expr("num_comments + 1")).Error; err != nil {
			return nil, err
		}

		// Check attachments
		attachments := make([]*Attachment, 0, len(opts.Attachments))
		for _, uuid := range opts.Attachments {
			attach, err := getAttachmentByUUID(tx, uuid)
			if err != nil {
				if IsErrAttachmentNotExist(err) {
					continue
				}
				return nil, fmt.Errorf("getAttachmentByUUID [%s]: %v", uuid, err)
			}
			attachments = append(attachments, attach)
		}

		for i := range attachments {
			attachments[i].IssueID = opts.Issue.ID
			attachments[i].CommentID = comment.ID
			if err = tx.Model(attachments[i]).Updates(map[string]any{
				"issue_id":   opts.Issue.ID,
				"comment_id": comment.ID,
			}).Error; err != nil {
				return nil, fmt.Errorf("update attachment [%d]: %v", attachments[i].ID, err)
			}
		}

	case CommentTypeReopen:
		act.OpType = ActionReopenIssue
		if opts.Issue.IsPull {
			act.OpType = ActionReopenPullRequest
		}

		if opts.Issue.IsPull {
			err = tx.Model(new(Repository)).Where("id = ?", opts.Repo.ID).Update("num_closed_pulls", gorm.Expr("num_closed_pulls - 1")).Error
		} else {
			err = tx.Model(new(Repository)).Where("id = ?", opts.Repo.ID).Update("num_closed_issues", gorm.Expr("num_closed_issues - 1")).Error
		}
		if err != nil {
			return nil, err
		}

	case CommentTypeClose:
		act.OpType = ActionCloseIssue
		if opts.Issue.IsPull {
			act.OpType = ActionClosePullRequest
		}

		if opts.Issue.IsPull {
			err = tx.Model(new(Repository)).Where("id = ?", opts.Repo.ID).Update("num_closed_pulls", gorm.Expr("num_closed_pulls + 1")).Error
		} else {
			err = tx.Model(new(Repository)).Where("id = ?", opts.Repo.ID).Update("num_closed_issues", gorm.Expr("num_closed_issues + 1")).Error
		}
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Model(new(Issue)).Where("id = ?", opts.Issue.ID).Update("updated_unix", time.Now().Unix()).Error; err != nil {
		return nil, fmt.Errorf("update issue 'updated_unix': %v", err)
	}

	// Notify watchers for whatever action comes in, ignore if no action type.
	if act.OpType > 0 {
		if err = notifyWatchers(tx, act); err != nil {
			log.Error("notifyWatchers: %v", err)
		}
		if err = comment.mailParticipants(tx, act.OpType, opts.Issue); err != nil {
			log.Error("MailParticipants: %v", err)
		}
	}

	return comment, comment.loadAttributes(tx)
}

func createStatusComment(tx *gorm.DB, doer *User, repo *Repository, issue *Issue) (*Comment, error) {
	cmtType := CommentTypeClose
	if !issue.IsClosed {
		cmtType = CommentTypeReopen
	}
	return createComment(tx, &CreateCommentOptions{
		Type:  cmtType,
		Doer:  doer,
		Repo:  repo,
		Issue: issue,
	})
}

type CreateCommentOptions struct {
	Type  CommentType
	Doer  *User
	Repo  *Repository
	Issue *Issue

	CommitID    int64
	CommitSHA   string
	LineNum     int64
	Content     string
	Attachments []string // UUIDs of attachments
}

// CreateComment creates comment of issue or commit.
func CreateComment(opts *CreateCommentOptions) (comment *Comment, err error) {
	err = Handle.db.Transaction(func(tx *gorm.DB) error {
		comment, err = createComment(tx, opts)
		return err
	})
	return comment, err
}

// CreateIssueComment creates a plain issue comment.
func CreateIssueComment(doer *User, repo *Repository, issue *Issue, content string, attachments []string) (*Comment, error) {
	comment, err := CreateComment(&CreateCommentOptions{
		Type:        CommentTypeComment,
		Doer:        doer,
		Repo:        repo,
		Issue:       issue,
		Content:     content,
		Attachments: attachments,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateComment: %v", err)
	}

	comment.Issue = issue
	if err = PrepareWebhooks(repo, HookEventTypeIssueComment, &apiv2types.WebhookIssueCommentPayload{
		Action:     apiv2types.WebhookIssueCommentCreated,
		Issue:      issue.APIFormat(),
		Comment:    comment.APIFormat(),
		Repository: repo.APIFormatLegacy(nil),
		Sender:     doer.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks [comment_id: %d]: %v", comment.ID, err)
	}

	return comment, nil
}

// CreateRefComment creates a commit reference comment to issue.
func CreateRefComment(doer *User, repo *Repository, issue *Issue, content, commitSHA string) error {
	if commitSHA == "" {
		return fmt.Errorf("cannot create reference with empty commit SHA")
	}

	// Check if same reference from same commit has already existed.
	err := Handle.db.Where("type = ? AND issue_id = ? AND commit_sha = ?", CommentTypeCommitRef, issue.ID, commitSHA).First(new(Comment)).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check reference comment: %v", err)
		}
	} else {
		return nil
	}

	_, err = CreateComment(&CreateCommentOptions{
		Type:      CommentTypeCommitRef,
		Doer:      doer,
		Repo:      repo,
		Issue:     issue,
		CommitSHA: commitSHA,
		Content:   content,
	})
	return err
}

var _ errx.NotFound = (*ErrCommentNotExist)(nil)

type ErrCommentNotExist struct {
	args map[string]any
}

func IsErrCommentNotExist(err error) bool {
	_, ok := err.(ErrCommentNotExist)
	return ok
}

func (err ErrCommentNotExist) Error() string {
	return fmt.Sprintf("comment does not exist: %v", err.args)
}

func (ErrCommentNotExist) NotFound() bool {
	return true
}

func getCommentByID(db *gorm.DB, id int64) (*Comment, error) {
	c := new(Comment)
	err := db.Where("id = ?", id).First(c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotExist{args: map[string]any{"commentID": id}}
		}
		return nil, err
	}
	return c, c.loadAttributes(db)
}

// GetCommentByID returns the comment by given ID.
func GetCommentByID(id int64) (*Comment, error) {
	return getCommentByID(Handle.db, id)
}

// FIXME: use CommentList to improve performance.
func loadCommentsAttributes(tx *gorm.DB, comments []*Comment) (err error) {
	if len(comments) == 0 {
		return nil
	}

	posterIDs := make([]int64, 0, len(comments))
	issueIDs := make([]int64, 0, len(comments))
	commentIDs := make([]int64, len(comments))
	needIssues := false

	for i, c := range comments {
		commentIDs[i] = c.ID
		if c.Poster == nil {
			posterIDs = append(posterIDs, c.PosterID)
		}
		if c.Issue == nil {
			needIssues = true
			issueIDs = append(issueIDs, c.IssueID)
		}
	}

	// Batch load posters
	if len(posterIDs) > 0 {
		posters := make([]User, 0, len(posterIDs))
		if err = tx.Where("id IN ?", posterIDs).Find(&posters).Error; err != nil {
			return fmt.Errorf("batch load posters: %v", err)
		}
		posterMap := make(map[int64]*User, len(posters))
		for i := range posters {
			posterMap[posters[i].ID] = &posters[i]
		}
		for _, c := range comments {
			if c.Poster == nil {
				if p, ok := posterMap[c.PosterID]; ok {
					c.Poster = p
				} else {
					c.PosterID = -1
					c.Poster = NewGhostUser()
				}
			}
		}
	}

	// Batch load issues + repos
	if needIssues && len(issueIDs) > 0 {
		issues := make([]Issue, 0, len(issueIDs))
		if err = tx.Where("id IN ?", issueIDs).Find(&issues).Error; err != nil {
			return fmt.Errorf("batch load issues: %v", err)
		}
		issueMap := make(map[int64]*Issue, len(issues))
		repoIDs := make([]int64, 0, len(issues))
		for i := range issues {
			issueMap[issues[i].ID] = &issues[i]
			repoIDs = append(repoIDs, issues[i].RepoID)
		}

		// Batch load repos
		repos := make([]Repository, 0, len(repoIDs))
		if err = tx.Where("id IN ?", repoIDs).Find(&repos).Error; err != nil {
			return fmt.Errorf("batch load repos: %v", err)
		}
		repoMap := make(map[int64]*Repository, len(repos))
		for i := range repos {
			repoMap[repos[i].ID] = &repos[i]
		}

		for _, c := range comments {
			if c.Issue == nil {
				if issue, ok := issueMap[c.IssueID]; ok {
					issue.Repo = repoMap[issue.RepoID]
					c.Issue = issue
				}
			}
		}
	}

	// Batch load attachments
	allAttachments := make([]Attachment, 0, len(comments))
	if err = tx.Where("comment_id IN ?", commentIDs).Find(&allAttachments).Error; err != nil {
		return fmt.Errorf("batch load attachments: %v", err)
	}
	attachmentsByComment := make(map[int64][]*Attachment, len(comments))
	for i := range allAttachments {
		cid := allAttachments[i].CommentID
		attachmentsByComment[cid] = append(attachmentsByComment[cid], &allAttachments[i])
	}
	for _, c := range comments {
		if c.Attachments == nil {
			c.Attachments = attachmentsByComment[c.ID]
		}
	}

	return nil
}

func getCommentsByIssueIDSince(tx *gorm.DB, issueID, since int64) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	query := tx.Where("issue_id = ?", issueID).Order("created_unix ASC")
	if since > 0 {
		query = query.Where("updated_unix >= ?", since)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, loadCommentsAttributes(tx, comments)
}

func getCommentsByRepoIDSince(tx *gorm.DB, repoID, since int64) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	query := tx.Table("comment").
		Joins("JOIN issue ON issue.id = comment.issue_id").
		Where("issue.repo_id = ?", repoID).
		Order("comment.created_unix ASC")
	if since > 0 {
		query = query.Where("comment.updated_unix >= ?", since)
	}
	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, loadCommentsAttributes(tx, comments)
}

func getCommentsByIssueID(tx *gorm.DB, issueID int64) ([]*Comment, error) {
	return getCommentsByIssueIDSince(tx, issueID, -1)
}

// GetCommentsByIssueID returns all comments of an issue.
func GetCommentsByIssueID(issueID int64) ([]*Comment, error) {
	return getCommentsByIssueID(Handle.db, issueID)
}

// GetCommentsByIssueIDSince returns a list of comments of an issue since a given time point.
func GetCommentsByIssueIDSince(issueID, since int64) ([]*Comment, error) {
	return getCommentsByIssueIDSince(Handle.db, issueID, since)
}

// GetCommentsByRepoIDSince returns a list of comments for all issues in a repo since a given time point.
func GetCommentsByRepoIDSince(repoID, since int64) ([]*Comment, error) {
	return getCommentsByRepoIDSince(Handle.db, repoID, since)
}

// UpdateComment updates information of comment.
func UpdateComment(doer *User, c *Comment, oldContent string) (err error) {
	if err = Handle.db.Model(c).Updates(c).Error; err != nil {
		return err
	}

	if err = c.Issue.LoadAttributes(); err != nil {
		log.Error("Issue.LoadAttributes [issue_id: %d]: %v", c.IssueID, err)
	} else if err = PrepareWebhooks(c.Issue.Repo, HookEventTypeIssueComment, &apiv2types.WebhookIssueCommentPayload{
		Action:  apiv2types.WebhookIssueCommentEdited,
		Issue:   c.Issue.APIFormat(),
		Comment: c.APIFormat(),
		Changes: &apiv2types.WebhookChangesPayload{
			Body: &apiv2types.WebhookChangesFromPayload{
				From: oldContent,
			},
		},
		Repository: c.Issue.Repo.APIFormatLegacy(nil),
		Sender:     doer.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks [comment_id: %d]: %v", c.ID, err)
	}

	return nil
}

// DeleteCommentByID deletes the comment by given ID.
func DeleteCommentByID(doer *User, id int64) error {
	comment, err := GetCommentByID(id)
	if err != nil {
		if IsErrCommentNotExist(err) {
			return nil
		}
		return err
	}

	err = Handle.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(comment).Error; err != nil {
			return err
		}

		if comment.Type == CommentTypeComment {
			if err = tx.Model(new(Issue)).Where("id = ?", comment.IssueID).Update("num_comments", gorm.Expr("num_comments - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete comment: %v", err)
	}

	_, err = DeleteAttachmentsByComment(comment.ID, true)
	if err != nil {
		log.Error("Failed to delete attachments by comment[%d]: %v", comment.ID, err)
	}

	if err = comment.Issue.LoadAttributes(); err != nil {
		log.Error("Issue.LoadAttributes [issue_id: %d]: %v", comment.IssueID, err)
	} else if err = PrepareWebhooks(comment.Issue.Repo, HookEventTypeIssueComment, &apiv2types.WebhookIssueCommentPayload{
		Action:     apiv2types.WebhookIssueCommentDeleted,
		Issue:      comment.Issue.APIFormat(),
		Comment:    comment.APIFormat(),
		Repository: comment.Issue.Repo.APIFormatLegacy(nil),
		Sender:     doer.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks [comment_id: %d]: %v", comment.ID, err)
	}
	return nil
}
