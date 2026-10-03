package database

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbtest"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// countingLogger counts the statements GORM executes. It is the measurement tool
// behind the N+1 assertions in this file: the regression guards are written as
// "the count must not grow with the number of issues", which holds regardless of
// the absolute number of queries a future refactor happens to need.
type countingLogger struct {
	queries atomic.Int64
	inner   logger.Interface
}

func (l *countingLogger) LogMode(mode logger.LogLevel) logger.Interface { return l }

func (l *countingLogger) Info(ctx context.Context, msg string, args ...any) {
	l.inner.Info(ctx, msg, args...)
}

func (l *countingLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.inner.Warn(ctx, msg, args...)
}

func (l *countingLogger) Error(ctx context.Context, msg string, args ...any) {
	l.inner.Error(ctx, msg, args...)
}

func (l *countingLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.queries.Add(1)
	l.inner.Trace(ctx, begin, fc, err)
}

// seedIssues creates one repository, one user and count issues, and returns the
// repository ID.
func seedIssues(t *testing.T, db *gorm.DB, count int) int64 {
	t.Helper()
	ctx := context.Background()

	domainUser, err := Handle.Users().Create(ctx, "issues-owner", "owner@example.com", CreateUserOptions{})
	require.NoError(t, err)
	user := FromDomainUser(domainUser)

	repo, err := Handle.Repositories().Create(ctx, user.ID, CreateRepoOptions{
		Name: "issues-repo",
	})
	require.NoError(t, err)

	domain := FromDomainRepository(repo)
	for i := 1; i <= count; i++ {
		issue := &Issue{
			RepoID:   domain.ID,
			Index:    int64(i),
			Title:    "issue " + strconv.Itoa(i),
			PosterID: domain.OwnerID,
			Poster:   user,
			Repo:     domain,
		}
		require.NoError(t, db.Create(issue).Error)
	}

	return domain.ID
}

// TestIssuesQueryCountIsIndependentOfPageSize is the regression guard for the
// N+1 in Issues. It used to call LoadAttributes per row, which fans out into
// getRepositoryByID, the poster, the labels, the milestone, the assignee, the
// pull request, the attachments and the comments: about a dozen queries for
// every issue on the page.
//
// The assertion is on the growth, not on an absolute count: doubling the page
// must not double the statements.
func TestIssuesQueryCountIsIndependentOfPageSize(t *testing.T) {
	previousPaging := conf.UI.IssuePagingNum
	t.Cleanup(func() { conf.UI.IssuePagingNum = previousPaging })

	base := dbtest.NewDB(t, "IssuesQueryCount",
		new(User), new(Repository), new(Issue), new(Label), new(IssueLabel),
		new(Milestone), new(Attachment), new(Comment), new(PullRequest), new(Watch),
	)
	SetHandle(base)
	seedIssues(t, base, 40)

	count := func(pageSize int) int64 {
		conf.UI.IssuePagingNum = pageSize
		SetHandle(base)

		counter := &countingLogger{inner: logger.Discard}
		counting := base.Session(&gorm.Session{Logger: counter, NewDB: true})
		SetHandle(counting)

		issues, err := Issues(&IssuesOptions{RepoID: 1, Page: 1})
		require.NoError(t, err)
		assert.Len(t, issues, pageSize, "la pagina deberia traer %d issues", pageSize)

		return counter.queries.Load()
	}

	small := count(5)
	large := count(20)

	t.Logf("5 issues -> %d consultas, 20 issues -> %d consultas", small, large)
	assert.LessOrEqual(t, large, small,
		"las consultas no deben crecer con el numero de issues: %d con 5 frente a %d con 20", small, large)
}

// TestIssueListMatchesPerIssueLoading is the guard that the batch loader did not
// change any result. It runs the same issues through both paths and compares the
// fields APIFormat serialises.
//
// A faster loader that returns a different answer would be worse than the N+1, so
// this is the test that matters most here.
func TestIssueListMatchesPerIssueLoading(t *testing.T) {
	// suite keeps the database files apart; suffix is part of the seeded names and
	// must be identical for both runs, otherwise the comparison would be between
	// two different datasets rather than two loaders.
	build := func(t *testing.T, suite, suffix string) []*Issue {
		t.Helper()
		db := dbtest.NewDB(t, "IssueListEquivalence"+suite,
			new(User), new(Repository), new(Issue), new(Label), new(IssueLabel),
			new(Milestone), new(Attachment), new(Comment), new(PullRequest), new(Watch))
		SetHandle(db)

		// The creation timestamp is pinned. Left to the clock, two runs straddling a
		// second boundary produce users whose CreatedAt differ by one second, and the
		// comparison below fails for a reason that has nothing to do with the loader.
		seeded := time.Unix(1700000000, 0)

		owner, err := Handle.Users().Create(context.Background(), "eq-owner-"+suffix, "eq-"+suffix+"@example.com", CreateUserOptions{})
		require.NoError(t, err)
		assignee, err := Handle.Users().Create(context.Background(), "eq-assignee-"+suffix, "as-"+suffix+"@example.com", CreateUserOptions{})
		require.NoError(t, err)

		for _, u := range []*User{FromDomainUser(owner), FromDomainUser(assignee)} {
			require.NoError(t, db.Model(u).Updates(map[string]any{
				"created_unix": seeded.Unix(),
				"updated_unix": seeded.Unix(),
			}).Error)
		}

		repo, err := Handle.Repositories().Create(context.Background(), owner.ID, CreateRepoOptions{Name: "eq-repo-" + suffix})
		require.NoError(t, err)
		domainRepo := FromDomainRepository(repo)

		milestone := &Milestone{RepoID: domainRepo.ID, Name: "m"}
		require.NoError(t, db.Create(milestone).Error)

		labels := []*Label{
			{RepoID: domainRepo.ID, Name: "bug", Color: "f00"},
			{RepoID: domainRepo.ID, Name: "chore", Color: "0f0"},
		}
		for _, l := range labels {
			require.NoError(t, db.Create(l).Error)
		}

		issues := make([]*Issue, 0, 4)
		for i := 1; i <= 4; i++ {
			issue := &Issue{
				RepoID:      domainRepo.ID,
				Index:       int64(i),
				Title:       "issue " + suffix,
				PosterID:    owner.ID,
				Poster:      FromDomainUser(owner),
				Repo:        domainRepo,
				MilestoneID: milestone.ID,
				AssigneeID:  assignee.ID,
				NumComments: i,
			}
			require.NoError(t, db.Create(issue).Error)

			// Odd issues get every label and even ones get none, so both the empty
			// and the populated path are covered.
			if i%2 == 1 {
				for _, l := range labels {
					require.NoError(t, db.Create(&IssueLabel{IssueID: issue.ID, LabelID: l.ID}).Error)
				}
			}
			issues = append(issues, issue)
		}
		return issues
	}

	// Per-issue path.
	perIssue := build(t, "PerIssue", "d")
	for i := range perIssue {
		perIssue[i].Labels = nil
		perIssue[i].Milestone = nil
		perIssue[i].Assignee = nil
		require.NoError(t, perIssue[i].LoadAttributes())
	}
	perIssueFormatted := make([]*apiv2types.IssueResponse, len(perIssue))
	for i := range perIssue {
		perIssueFormatted[i] = perIssue[i].APIFormat()
	}

	// Batch path.
	batched := build(t, "Batched", "d")
	for i := range batched {
		batched[i].Labels = nil
		batched[i].Milestone = nil
		batched[i].Assignee = nil
	}
	require.NoError(t, IssueList(batched).LoadAttributes())
	batchFormatted := make([]*apiv2types.IssueResponse, len(batched))
	for i := range batched {
		batchFormatted[i] = batched[i].APIFormat()
	}

	require.Len(t, batchFormatted, len(perIssueFormatted))
	for i := range perIssueFormatted {
		a := perIssueFormatted[i]
		b := batchFormatted[i]

		assert.Equal(t, a.ID, b.ID, "issue %d: ID", i)
		assert.Equal(t, a.Index, b.Index, "issue %d: Index", i)
		assert.Equal(t, a.Title, b.Title, "issue %d: Title", i)
		assert.Equal(t, a.Comments, b.Comments, "issue %d: Comments", i)
		assert.Equal(t, a.Labels, b.Labels, "issue %d: Labels", i)
		assert.Equal(t, a.Poster, b.Poster, "issue %d: Poster", i)
		assert.Equal(t, a.Assignee, b.Assignee, "issue %d: Assignee", i)
		assert.Equal(t, a.Milestone, b.Milestone, "issue %d: Milestone", i)
	}
}
