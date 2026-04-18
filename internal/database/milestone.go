package database

import (
	"fmt"
	"time"

	log "unknwon.dev/clog/v2"
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/errx"
	apiv2types "gogs.io/gogs/internal/route/api/v2/types"
)

// Milestone represents a milestone of repository.
type Milestone struct {
	ID              int64  `gorm:"primaryKey"`
	RepoID          int64  `xorm:"INDEX" gorm:"index"`
	Name            string `gorm:"not null"`
	Content         string `xorm:"TEXT" gorm:"type:TEXT"`
	RenderedContent string `xorm:"-" json:"-" gorm:"-"`
	IsClosed        bool
	NumIssues       int
	NumClosedIssues int
	NumOpenIssues   int `xorm:"-" json:"-" gorm:"-"`
	Completeness    int // Percentage(1-100).
	IsOverDue       bool `xorm:"-" json:"-" gorm:"-"`

	DeadlineString string    `xorm:"-" json:"-" gorm:"-"`
	Deadline       time.Time `xorm:"-" json:"-" gorm:"-"`
	DeadlineUnix   int64
	ClosedDate     time.Time `xorm:"-" json:"-" gorm:"-"`
	ClosedDateUnix int64
}

// BeforeCreate implements the GORM create hook.
func (m *Milestone) BeforeCreate(tx *gorm.DB) error {
	if m.DeadlineUnix == 0 {
		m.DeadlineUnix = m.Deadline.Unix()
	}
	return nil
}

// BeforeUpdate implements the GORM update hook.
func (m *Milestone) BeforeUpdate(tx *gorm.DB) error {
	if m.NumIssues > 0 {
		m.Completeness = m.NumClosedIssues * 100 / m.NumIssues
	} else {
		m.Completeness = 0
	}

	m.DeadlineUnix = m.Deadline.Unix()
	m.ClosedDateUnix = m.ClosedDate.Unix()
	return nil
}

// AfterFind implements the GORM query hook.
func (m *Milestone) AfterFind(_ *gorm.DB) error {
	m.NumOpenIssues = m.NumIssues - m.NumClosedIssues

	m.Deadline = time.Unix(m.DeadlineUnix, 0).Local()
	if m.Deadline.Year() != 9999 {
		m.DeadlineString = m.Deadline.Format("2006-01-02")
		if time.Now().Local().After(m.Deadline) {
			m.IsOverDue = true
		}
	}

	m.ClosedDate = time.Unix(m.ClosedDateUnix, 0).Local()
	return nil
}


// State returns string representation of milestone status.
func (m *Milestone) State() apiv2types.IssueStateType {
	if m.IsClosed {
		return apiv2types.IssueStateClosed
	}
	return apiv2types.IssueStateOpen
}

func (m *Milestone) ChangeStatus(isClosed bool) error {
	return ChangeMilestoneStatus(m, isClosed)
}

func (m *Milestone) APIFormat() *apiv2types.IssueMilestoneResponse {
	return &apiv2types.IssueMilestoneResponse{
		ID:           m.ID,
		Title:        m.Name,
		Description:  m.Content,
		State:        m.State(),
		OpenIssues:   m.NumOpenIssues,
		ClosedIssues: m.NumClosedIssues,
		Deadline:     &m.Deadline,
	}
}

func (m *Milestone) CountIssues(isClosed, includePulls bool) int64 {
	query := Handle.db.Model(new(Issue)).Where("milestone_id = ? AND is_closed = ?", m.ID, isClosed)
	if !includePulls {
		query = query.Where("is_pull = ?", false)
	}
	var count int64
	query.Count(&count)
	return count
}

// NewMilestone creates new milestone of repository.
func NewMilestone(m *Milestone) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(m).Error; err != nil {
			return err
		}

		return tx.Model(new(Repository)).Where("id = ?", m.RepoID).Update("num_milestones", gorm.Expr("num_milestones + 1")).Error
	})
}

var _ errx.NotFound = (*ErrMilestoneNotExist)(nil)

type ErrMilestoneNotExist struct {
	args map[string]any
}

func IsErrMilestoneNotExist(err error) bool {
	_, ok := err.(ErrMilestoneNotExist)
	return ok
}

func (err ErrMilestoneNotExist) Error() string {
	return fmt.Sprintf("milestone does not exist: %v", err.args)
}

func (ErrMilestoneNotExist) NotFound() bool {
	return true
}

func getMilestoneByRepoID(tx *gorm.DB, repoID, id int64) (*Milestone, error) {
	m := new(Milestone)
	err := tx.Where("id = ? AND repo_id = ?", id, repoID).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMilestoneNotExist{args: map[string]any{"repoID": repoID, "milestoneID": id}}
		}
		return nil, err
	}
	return m, nil
}

// GetWebhookByRepoID returns the milestone in a repository.
func GetMilestoneByRepoID(repoID, id int64) (*Milestone, error) {
	return getMilestoneByRepoID(Handle.db, repoID, id)
}

// GetMilestonesByRepoID returns all milestones of a repository.
func GetMilestonesByRepoID(repoID int64) ([]*Milestone, error) {
	miles := make([]*Milestone, 0, 10)
	err := Handle.db.Where("repo_id = ?", repoID).Find(&miles).Error
	return miles, err
}

// GetMilestones returns a list of milestones of given repository and status.
func GetMilestones(repoID int64, page int, isClosed bool) ([]*Milestone, error) {
	miles := make([]*Milestone, 0, conf.UI.IssuePagingNum)
	query := Handle.db.Where("repo_id = ? AND is_closed = ?", repoID, isClosed)
	if page > 0 {
		query = query.Limit(conf.UI.IssuePagingNum).Offset((page - 1) * conf.UI.IssuePagingNum)
	}
	err := query.Find(&miles).Error
	return miles, err
}

func updateMilestone(tx *gorm.DB, m *Milestone) error {
	return tx.Model(m).Updates(m).Error
}

// UpdateMilestone updates information of given milestone.
func UpdateMilestone(m *Milestone) error {
	return updateMilestone(Handle.db, m)
}

func countRepoMilestones(tx *gorm.DB, repoID int64) int64 {
	var count int64
	tx.Model(new(Milestone)).Where("repo_id = ?", repoID).Count(&count)
	return count
}

// CountRepoMilestones returns number of milestones in given repository.
func CountRepoMilestones(repoID int64) int64 {
	return countRepoMilestones(Handle.db, repoID)
}

func countRepoClosedMilestones(tx *gorm.DB, repoID int64) int64 {
	var count int64
	tx.Model(new(Milestone)).Where("repo_id = ? AND is_closed = ?", repoID, true).Count(&count)
	return count
}

// CountRepoClosedMilestones returns number of closed milestones in given repository.
func CountRepoClosedMilestones(repoID int64) int64 {
	return countRepoClosedMilestones(Handle.db, repoID)
}

// MilestoneStats returns number of open and closed milestones of given repository.
func MilestoneStats(repoID int64) (open, closed int64) {
	Handle.db.Model(new(Milestone)).Where("repo_id = ? AND is_closed = ?", repoID, false).Count(&open)
	return open, CountRepoClosedMilestones(repoID)
}

// ChangeMilestoneStatus changes the milestone open/closed status.
// If milestone passes with changed values, those values will be
// updated to database as well.
func ChangeMilestoneStatus(m *Milestone, isClosed bool) (err error) {
	repo, err := GetRepositoryByID(m.RepoID)
	if err != nil {
		return err
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		m.IsClosed = isClosed
		if err = updateMilestone(tx, m); err != nil {
			return err
		}

		repo.NumMilestones = int(countRepoMilestones(tx, repo.ID))
		repo.NumClosedMilestones = int(countRepoClosedMilestones(tx, repo.ID))
		return tx.Model(repo).Updates(repo).Error
	})
}

func changeMilestoneIssueStats(tx *gorm.DB, issue *Issue) error {
	if issue.MilestoneID == 0 {
		return nil
	}

	m, err := getMilestoneByRepoID(tx, issue.RepoID, issue.MilestoneID)
	if err != nil {
		return err
	}

	if issue.IsClosed {
		m.NumOpenIssues--
		m.NumClosedIssues++
	} else {
		m.NumOpenIssues++
		m.NumClosedIssues--
	}

	return updateMilestone(tx, m)
}

// ChangeMilestoneIssueStats updates the open/closed issues counter and progress
// for the milestone associated with the given issue.
func ChangeMilestoneIssueStats(issue *Issue) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		return changeMilestoneIssueStats(tx, issue)
	})
}

func changeMilestoneAssign(tx *gorm.DB, issue *Issue, oldMilestoneID int64) error {
	if oldMilestoneID > 0 {
		m, err := getMilestoneByRepoID(tx, issue.RepoID, oldMilestoneID)
		if err != nil {
			return err
		}

		m.NumIssues--
		if issue.IsClosed {
			m.NumClosedIssues--
		}

		if err = updateMilestone(tx, m); err != nil {
			return err
		} else if err = tx.Model(new(IssueUser)).Where("issue_id = ?", issue.ID).Update("milestone_id", 0).Error; err != nil {
			return err
		}

		issue.Milestone = nil
	}

	if issue.MilestoneID > 0 {
		m, err := getMilestoneByRepoID(tx, issue.RepoID, issue.MilestoneID)
		if err != nil {
			return err
		}

		m.NumIssues++
		if issue.IsClosed {
			m.NumClosedIssues++
		}

		if err = updateMilestone(tx, m); err != nil {
			return err
		} else if err = tx.Model(new(IssueUser)).Where("issue_id = ?", issue.ID).Update("milestone_id", m.ID).Error; err != nil {
			return err
		}

		issue.Milestone = m
	}

	return updateIssue(tx, issue)
}

// ChangeMilestoneAssign changes assignment of milestone for issue.
func ChangeMilestoneAssign(doer *User, issue *Issue, oldMilestoneID int64) (err error) {
	err = Handle.db.Transaction(func(tx *gorm.DB) error {
		return changeMilestoneAssign(tx, issue, oldMilestoneID)
	})
	if err != nil {
		return err
	}

	var hookAction apiv2types.WebhookIssueAction
	if issue.MilestoneID > 0 {
		hookAction = apiv2types.WebhookIssueMilestoned
	} else {
		hookAction = apiv2types.WebhookIssueDemilestoned
	}

	if issue.IsPull {
		err = issue.PullRequest.LoadIssue()
		if err != nil {
			log.Error("LoadIssue: %v", err)
			return err
		}
		err = PrepareWebhooks(issue.Repo, HookEventTypePullRequest, &apiv2types.WebhookPullRequestPayload{
			Action:      hookAction,
			Index:       issue.Index,
			PullRequest: issue.PullRequest.APIFormat(),
			Repository:  issue.Repo.APIFormatLegacy(nil),
			Sender:      doer.APIFormat(),
		})
	} else {
		err = PrepareWebhooks(issue.Repo, HookEventTypeIssues, &apiv2types.WebhookIssuesPayload{
			Action:     hookAction,
			Index:      issue.Index,
			Issue:      issue.APIFormat(),
			Repository: issue.Repo.APIFormatLegacy(nil),
			Sender:     doer.APIFormat(),
		})
	}
	if err != nil {
		log.Error("PrepareWebhooks [is_pull: %v]: %v", issue.IsPull, err)
	}

	return nil
}

// DeleteMilestoneOfRepoByID deletes a milestone from a repository.
func DeleteMilestoneOfRepoByID(repoID, id int64) error {
	m, err := GetMilestoneByRepoID(repoID, id)
	if err != nil {
		if IsErrMilestoneNotExist(err) {
			return nil
		}
		return err
	}

	repo, err := GetRepositoryByID(m.RepoID)
	if err != nil {
		return err
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("id = ?", m.ID).Delete(new(Milestone)).Error; err != nil {
			return err
		}

		repo.NumMilestones = int(countRepoMilestones(tx, repo.ID))
		repo.NumClosedMilestones = int(countRepoClosedMilestones(tx, repo.ID))
		if err = tx.Model(repo).Updates(repo).Error; err != nil {
			return err
		}

		if err = tx.Model(new(Issue)).Where("milestone_id = ?", m.ID).Update("milestone_id", 0).Error; err != nil {
			return err
		}
		return tx.Model(new(IssueUser)).Where("milestone_id = ?", m.ID).Update("milestone_id", 0).Error
	})
}
