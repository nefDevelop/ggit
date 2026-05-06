package database

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/errx"
)

// Attachment represent a attachment of issue/comment/release.
type Attachment struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	UUID      string `json:"uuid" gorm:"unique;not null"`
	IssueID   int64  `json:"issue_id" gorm:"index"`
	CommentID int64  `json:"comment_id" gorm:"index"`
	ReleaseID int64  `json:"release_id" gorm:"index"`
	Name      string `json:"name" gorm:"not null"`

	Created     time.Time `json:"-" gorm:"-"`
	CreatedUnix int64     `json:"created_at"`
}

// BeforeCreate implements the GORM create hook.
func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	a.CreatedUnix = tx.NowFunc().Unix()
	return nil
}

// AfterFind implements the GORM query hook.
func (a *Attachment) AfterFind(_ *gorm.DB) error {
	a.Created = time.Unix(a.CreatedUnix, 0).Local()
	return nil
}

// AttachmentLocalPath returns where attachment is stored in local file system based on given UUID.
func AttachmentLocalPath(uuid string) string {
	return path.Join(conf.Attachment.Path, uuid[0:1], uuid[1:2], uuid)
}

// LocalPath returns where attachment is stored in local file system.
func (a *Attachment) LocalPath() string {
	return AttachmentLocalPath(a.UUID)
}

// NewAttachment creates a new attachment object.
func NewAttachment(name string, buf []byte, file multipart.File) (_ *Attachment, err error) {
	attach := &Attachment{
		UUID: uuid.New().String(),
		Name: name,
	}

	localPath := attach.LocalPath()
	if err = os.MkdirAll(path.Dir(localPath), os.ModePerm); err != nil {
		return nil, errors.Newf("MkdirAll: %v", err)
	}

	fw, err := os.Create(localPath)
	if err != nil {
		return nil, errors.Newf("Create: %v", err)
	}
	defer fw.Close()

	if _, err = fw.Write(buf); err != nil {
		return nil, errors.Newf("write: %v", err)
	} else if _, err = io.Copy(fw, file); err != nil {
		return nil, errors.Newf("copy: %v", err)
	}

	if err = Handle.db.Create(attach).Error; err != nil {
		return nil, err
	}

	return attach, nil
}

var _ errx.NotFound = (*ErrAttachmentNotExist)(nil)

type ErrAttachmentNotExist struct {
	args map[string]any
}

func IsErrAttachmentNotExist(err error) bool {
	_, ok := err.(ErrAttachmentNotExist)
	return ok
}

func (err ErrAttachmentNotExist) Error() string {
	return fmt.Sprintf("attachment does not exist: %v", err.args)
}

func (ErrAttachmentNotExist) NotFound() bool {
	return true
}

func getAttachmentByUUID(tx *gorm.DB, uuid string) (*Attachment, error) {
	attach := new(Attachment)
	err := tx.Where("uuid = ?", uuid).First(attach).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttachmentNotExist{args: map[string]any{"uuid": uuid}}
		}
		return nil, err
	}
	return attach, nil
}

func getAttachmentsByUUIDs(tx *gorm.DB, uuids []string) ([]*Attachment, error) {
	if len(uuids) == 0 {
		return []*Attachment{}, nil
	}

	// Silently drop invalid uuids.
	attachments := make([]*Attachment, 0, len(uuids))
	err := tx.Where("uuid IN (?)", uuids).Find(&attachments).Error
	return attachments, err
}

// GetAttachmentByUUID returns attachment by given UUID.
func GetAttachmentByUUID(uuid string) (*Attachment, error) {
	return getAttachmentByUUID(Handle.db, uuid)
}

func getAttachmentsByIssueID(tx *gorm.DB, issueID int64) ([]*Attachment, error) {
	attachments := make([]*Attachment, 0, 5)
	err := tx.Where("issue_id = ? AND comment_id = 0", issueID).Find(&attachments).Error
	return attachments, err
}

// GetAttachmentsByIssueID returns all attachments of an issue.
func GetAttachmentsByIssueID(issueID int64) ([]*Attachment, error) {
	return getAttachmentsByIssueID(Handle.db, issueID)
}

func getAttachmentsByCommentID(tx *gorm.DB, commentID int64) ([]*Attachment, error) {
	attachments := make([]*Attachment, 0, 5)
	err := tx.Where("comment_id = ?", commentID).Find(&attachments).Error
	return attachments, err
}

// GetAttachmentsByCommentID returns all attachments of a comment.
func GetAttachmentsByCommentID(commentID int64) ([]*Attachment, error) {
	return getAttachmentsByCommentID(Handle.db, commentID)
}

func getAttachmentsByReleaseID(tx *gorm.DB, releaseID int64) ([]*Attachment, error) {
	attachments := make([]*Attachment, 0, 10)
	err := tx.Where("release_id = ?", releaseID).Find(&attachments).Error
	return attachments, err
}

// GetAttachmentsByReleaseID returns all attachments of a release.
func GetAttachmentsByReleaseID(releaseID int64) ([]*Attachment, error) {
	return getAttachmentsByReleaseID(Handle.db, releaseID)
}

// DeleteAttachment deletes the given attachment and optionally the associated file.
func DeleteAttachment(a *Attachment, remove bool) error {
	_, err := DeleteAttachments([]*Attachment{a}, remove)
	return err
}

// DeleteAttachments deletes the given attachments and optionally the associated files.
func DeleteAttachments(attachments []*Attachment, remove bool) (int, error) {
	for i, a := range attachments {
		if remove {
			if err := os.Remove(a.LocalPath()); err != nil {
				return i, err
			}
		}

		if err := Handle.db.Where("id = ?", a.ID).Delete(new(Attachment)).Error; err != nil {
			return i, err
		}
	}

	return len(attachments), nil
}

// DeleteAttachmentsByIssue deletes all attachments associated with the given issue.
func DeleteAttachmentsByIssue(issueID int64, remove bool) (int, error) {
	attachments, err := GetAttachmentsByIssueID(issueID)
	if err != nil {
		return 0, err
	}

	return DeleteAttachments(attachments, remove)
}

// DeleteAttachmentsByComment deletes all attachments associated with the given comment.
func DeleteAttachmentsByComment(commentID int64, remove bool) (int, error) {
	attachments, err := GetAttachmentsByCommentID(commentID)
	if err != nil {
		return 0, err
	}

	return DeleteAttachments(attachments, remove)
}

// DeleteOrphanedAttachments deletes attachments that are not linked to any issue, comment, or release
// and are older than the specified duration.
func DeleteOrphanedAttachments(olderThan time.Duration) (int, error) {
	var attachments []*Attachment
	err := Handle.db.Where("issue_id = 0 AND comment_id = 0 AND release_id = 0 AND created_unix < ?",
		time.Now().Add(-olderThan).Unix()).Find(&attachments).Error
	if err != nil {
		return 0, err
	}

	return DeleteAttachments(attachments, true)
}
