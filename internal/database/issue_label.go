package database

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"

	"gogs.io/gogs/internal/errx"
	"gogs.io/gogs/internal/lazyregexp"
	apiv2types "gogs.io/gogs/internal/route/api/v2/types"
)

var labelColorPattern = lazyregexp.New("#([a-fA-F0-9]{6})")

// GetLabelTemplateFile loads the label template file by given name,
// then parses and returns a list of name-color pairs.
func GetLabelTemplateFile(name string) ([][2]string, error) {
	data, err := getRepoInitFile("label", name)
	if err != nil {
		return nil, errors.Newf("getRepoInitFile: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	list := make([][2]string, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.SplitN(line, " ", 2)
		if len(fields) != 2 {
			return nil, errors.Newf("line is malformed: %s", line)
		}

		if !labelColorPattern.MatchString(fields[0]) {
			return nil, errors.Newf("bad HTML color code in line: %s", line)
		}

		fields[1] = strings.TrimSpace(fields[1])
		list = append(list, [2]string{fields[1], fields[0]})
	}

	return list, nil
}

// Label represents a label of repository for issues.
type Label struct {
	ID              int64  `gorm:"primaryKey"`
	RepoID          int64  `xorm:"INDEX" gorm:"index"`
	Name            string `gorm:"not null"`
	Color           string `xorm:"VARCHAR(7)" gorm:"type:VARCHAR(7);not null"`
	NumIssues       int
	NumClosedIssues int
	NumOpenIssues   int  `xorm:"-" json:"-" gorm:"-"`
	IsChecked       bool `xorm:"-" json:"-" gorm:"-"`
}

// AfterFind implements the GORM query hook.
func (l *Label) AfterFind(_ *gorm.DB) error {
	l.NumOpenIssues = l.NumIssues - l.NumClosedIssues
	return nil
}

func (l *Label) APIFormat() *apiv2types.IssueLabelResponse {
	return &apiv2types.IssueLabelResponse{
		ID:    l.ID,
		Name:  l.Name,
		Color: l.Color,
	}
}

// CalOpenIssues calculates the open issues of label.
func (l *Label) CalOpenIssues() {
	l.NumOpenIssues = l.NumIssues - l.NumClosedIssues
}

// ForegroundColor calculates the text color for labels based
// on their background color.
func (l *Label) ForegroundColor() template.CSS {
	if strings.HasPrefix(l.Color, "#") {
		if color, err := strconv.ParseUint(l.Color[1:], 16, 64); err == nil {
			r := float32(0xFF & (color >> 16))
			g := float32(0xFF & (color >> 8))
			b := float32(0xFF & color)
			luminance := (0.2126*r + 0.7152*g + 0.0722*b) / 255

			if luminance < 0.66 {
				return template.CSS("#fff")
			}
		}
	}

	// default to black
	return template.CSS("#000")
}

// NewLabels creates new label(s) for a repository.
func NewLabels(labels ...*Label) error {
	return Handle.db.Create(labels).Error
}

var _ errx.NotFound = (*ErrLabelNotExist)(nil)

type ErrLabelNotExist struct {
	args map[string]any
}

func IsErrLabelNotExist(err error) bool {
	_, ok := err.(ErrLabelNotExist)
	return ok
}

func (err ErrLabelNotExist) Error() string {
	return fmt.Sprintf("label does not exist: %v", err.args)
}

func (ErrLabelNotExist) NotFound() bool {
	return true
}

// getLabelOfRepoByName returns a label by Name in given repository.
// If pass repoID as 0, then ORM will ignore limitation of repository
// and can return arbitrary label with any valid ID.
func getLabelOfRepoByName(tx *gorm.DB, repoID int64, labelName string) (*Label, error) {
	if len(labelName) <= 0 {
		return nil, ErrLabelNotExist{args: map[string]any{"repoID": repoID}}
	}

	l := new(Label)
	err := tx.Where("name = ? AND repo_id = ?", labelName, repoID).First(l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLabelNotExist{args: map[string]any{"repoID": repoID}}
		}
		return nil, err
	}
	return l, nil
}

// getLabelInRepoByID returns a label by ID in given repository.
// If pass repoID as 0, then ORM will ignore limitation of repository
// and can return arbitrary label with any valid ID.
func getLabelOfRepoByID(tx *gorm.DB, repoID, labelID int64) (*Label, error) {
	if labelID <= 0 {
		return nil, ErrLabelNotExist{args: map[string]any{"repoID": repoID, "labelID": labelID}}
	}

	l := new(Label)
	query := tx.Where("id = ?", labelID)
	if repoID > 0 {
		query = query.Where("repo_id = ?", repoID)
	}

	err := query.First(l).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLabelNotExist{args: map[string]any{"repoID": repoID, "labelID": labelID}}
		}
		return nil, err
	}
	return l, nil
}

// GetLabelByID returns a label by given ID.
func GetLabelByID(id int64) (*Label, error) {
	return getLabelOfRepoByID(Handle.db, 0, id)
}

// GetLabelOfRepoByID returns a label by ID in given repository.
func GetLabelOfRepoByID(repoID, labelID int64) (*Label, error) {
	return getLabelOfRepoByID(Handle.db, repoID, labelID)
}

// GetLabelOfRepoByName returns a label by name in given repository.
func GetLabelOfRepoByName(repoID int64, labelName string) (*Label, error) {
	return getLabelOfRepoByName(Handle.db, repoID, labelName)
}

// GetLabelsInRepoByIDs returns a list of labels by IDs in given repository,
// it silently ignores label IDs that are not belong to the repository.
func GetLabelsInRepoByIDs(repoID int64, labelIDs []int64) ([]*Label, error) {
	labels := make([]*Label, 0, len(labelIDs))
	err := Handle.db.Where("repo_id = ? AND id IN (?)", repoID, labelIDs).Order("name ASC").Find(&labels).Error
	return labels, err
}

// GetLabelsByRepoID returns all labels that belong to given repository by ID.
func GetLabelsByRepoID(repoID int64) ([]*Label, error) {
	labels := make([]*Label, 0, 10)
	err := Handle.db.Where("repo_id = ?", repoID).Order("name ASC").Find(&labels).Error
	return labels, err
}

func getLabelsByIssueID(tx *gorm.DB, issueID int64) ([]*Label, error) {
	issueLabels, err := getIssueLabels(tx, issueID)
	if err != nil {
		return nil, errors.Newf("getIssueLabels: %v", err)
	} else if len(issueLabels) == 0 {
		return []*Label{}, nil
	}

	labelIDs := make([]int64, len(issueLabels))
	for i := range issueLabels {
		labelIDs[i] = issueLabels[i].LabelID
	}

	labels := make([]*Label, 0, len(labelIDs))
	err = tx.Where("id IN (?)", labelIDs).Order("name ASC").Find(&labels).Error
	return labels, err
}

// GetLabelsByIssueID returns all labels that belong to given issue by ID.
func GetLabelsByIssueID(issueID int64) ([]*Label, error) {
	return getLabelsByIssueID(Handle.db, issueID)
}

func updateLabel(tx *gorm.DB, l *Label) error {
	return tx.Model(l).Updates(l).Error
}

// UpdateLabel updates label information.
func UpdateLabel(l *Label) error {
	return updateLabel(Handle.db, l)
}

// DeleteLabel delete a label of given repository.
func DeleteLabel(repoID, labelID int64) error {
	_, err := GetLabelOfRepoByID(repoID, labelID)
	if err != nil {
		if IsErrLabelNotExist(err) {
			return nil
		}
		return err
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", labelID).Delete(new(Label)).Error; err != nil {
			return err
		}
		return tx.Where("label_id = ?", labelID).Delete(new(IssueLabel)).Error
	})
}

// .___                            .____          ___.          .__
// |   | ______ ________ __   ____ |    |   _____ \_ |__   ____ |  |
// |   |/  ___//  ___/  |  \_/ __ \|    |   \__  \ | __ \_/ __ \|  |
// |   |\___ \ \___ \|  |  /\  ___/|    |___ / __ \| \_\ \  ___/|  |__
// |___/____  >____  >____/  \___  >_______ (____  /___  /\___  >____/
//          \/     \/            \/        \/    \/    \/     \/

// IssueLabel represents relation between issue and its label.
type IssueLabel struct {
	ID      int64 `gorm:"primaryKey"`
	IssueID int64 `xorm:"UNIQUE(s)" gorm:"uniqueIndex:issue_label_issue_label_unique"`
	LabelID int64 `xorm:"UNIQUE(s)" gorm:"uniqueIndex:issue_label_issue_label_unique"`
}

func hasIssueLabel(tx *gorm.DB, issueID, labelID int64) bool {
	var count int64
	tx.Model(new(IssueLabel)).Where("issue_id = ? AND label_id = ?", issueID, labelID).Count(&count)
	return count > 0
}

// HasIssueLabel returns true if issue has been labeled.
func HasIssueLabel(issueID, labelID int64) bool {
	return hasIssueLabel(Handle.db, issueID, labelID)
}

func newIssueLabel(tx *gorm.DB, issue *Issue, label *Label) (err error) {
	if err = tx.Create(&IssueLabel{
		IssueID: issue.ID,
		LabelID: label.ID,
	}).Error; err != nil {
		return err
	}

	label.NumIssues++
	if issue.IsClosed {
		label.NumClosedIssues++
	}

	if err = updateLabel(tx, label); err != nil {
		return errors.Newf("updateLabel: %v", err)
	}

	issue.Labels = append(issue.Labels, label)
	return nil
}

// NewIssueLabel creates a new issue-label relation.
func NewIssueLabel(issue *Issue, label *Label) (err error) {
	if HasIssueLabel(issue.ID, label.ID) {
		return nil
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		return newIssueLabel(tx, issue, label)
	})
}

func newIssueLabels(tx *gorm.DB, issue *Issue, labels []*Label) (err error) {
	for i := range labels {
		if hasIssueLabel(tx, issue.ID, labels[i].ID) {
			continue
		}

		if err = newIssueLabel(tx, issue, labels[i]); err != nil {
			return errors.Newf("newIssueLabel: %v", err)
		}
	}

	return nil
}

// NewIssueLabels creates a list of issue-label relations.
func NewIssueLabels(issue *Issue, labels []*Label) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		return newIssueLabels(tx, issue, labels)
	})
}

func getIssueLabels(tx *gorm.DB, issueID int64) ([]*IssueLabel, error) {
	issueLabels := make([]*IssueLabel, 0, 10)
	err := tx.Where("issue_id = ?", issueID).Order("label_id ASC").Find(&issueLabels).Error
	return issueLabels, err
}

// GetIssueLabels returns all issue-label relations of given issue by ID.
func GetIssueLabels(issueID int64) ([]*IssueLabel, error) {
	return getIssueLabels(Handle.db, issueID)
}

func deleteIssueLabel(tx *gorm.DB, issue *Issue, label *Label) (err error) {
	if err = tx.Where("issue_id = ? AND label_id = ?", issue.ID, label.ID).Delete(new(IssueLabel)).Error; err != nil {
		return err
	}

	label.NumIssues--
	if issue.IsClosed {
		label.NumClosedIssues--
	}
	if err = updateLabel(tx, label); err != nil {
		return errors.Newf("updateLabel: %v", err)
	}

	for i := range issue.Labels {
		if issue.Labels[i].ID == label.ID {
			issue.Labels = append(issue.Labels[:i], issue.Labels[i+1:]...)
			break
		}
	}
	return nil
}

// DeleteIssueLabel deletes issue-label relation.
func DeleteIssueLabel(issue *Issue, label *Label) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		return deleteIssueLabel(tx, issue, label)
	})
}
