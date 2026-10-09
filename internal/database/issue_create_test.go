package database

import (
	"context"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbtest"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tablas necesarias para crear issues por la vía completa (NewIssue) y para
// serializar PRs: contadores de repo, IssueUsers, adjuntos y comentarios.
func issueCreateTables() []any {
	return []any{
		new(User), new(Repository), new(Issue), new(IssueUser),
		new(Label), new(IssueLabel), new(Milestone),
		new(Attachment), new(Comment), new(Access),
		new(PullRequest), new(CommitStatus), new(Watch),
	}
}

func seedOwnerRepo(t *testing.T, username, email, reponame string) (*User, *Repository) {
	t.Helper()
	ctx := context.Background()

	domainUser, err := Handle.Users().Create(ctx, username, email, CreateUserOptions{})
	require.NoError(t, err)
	user := FromDomainUser(domainUser)

	repo, err := Handle.Repositories().Create(ctx, user.ID, CreateRepoOptions{Name: reponame})
	require.NoError(t, err)
	return user, FromDomainRepository(repo)
}

// TestNewIssueWithMilestone cubre el 500 ("WHERE conditions required") al crear
// un issue con milestone_id: newIssue llamaba a changeMilestoneAssign antes del
// INSERT, y el updateIssue final no tenía fila (ID == 0) que actualizar.
func TestNewIssueWithMilestone(t *testing.T) {
	db := dbtest.NewDB(t, "NewIssueWithMilestone", issueCreateTables()...)
	SetHandle(db)

	user, repo := seedOwnerRepo(t, "ms-owner", "ms-owner@example.com", "ms-repo")

	milestone := &Milestone{RepoID: repo.ID, Name: "v1.0"}
	require.NoError(t, db.Create(milestone).Error)

	issue := &Issue{
		RepoID:      repo.ID,
		Repo:        repo,
		Title:       "con hito",
		Content:     "cuerpo",
		PosterID:    user.ID,
		Poster:      user,
		MilestoneID: milestone.ID,
	}
	require.NoError(t, NewIssue(repo, issue, nil, nil))

	assert.Equal(t, int64(1), issue.Index)
	assert.Equal(t, milestone.ID, issue.MilestoneID)
	require.NotNil(t, issue.Milestone)

	var m Milestone
	require.NoError(t, db.First(&m, milestone.ID).Error)
	assert.Equal(t, 1, m.NumIssues, "el hito debe contar el issue")
}

// TestIssueListPullRequestFormats cubre el panic del listado de PRs
// (GET /:owner/:repo/pulls): el cargador por lotes dejaba BaseRepo a nil y
// APIFormat() caía en APIFormatLegacy.
func TestIssueListPullRequestFormats(t *testing.T) {
	previousPaging := conf.UI.IssuePagingNum
	conf.UI.IssuePagingNum = 20
	t.Cleanup(func() { conf.UI.IssuePagingNum = previousPaging })

	db := dbtest.NewDB(t, "IssueListPullFormat", issueCreateTables()...)
	SetHandle(db)

	user, repo := seedOwnerRepo(t, "pr-owner", "pr-owner@example.com", "pr-repo")

	issue := &Issue{
		RepoID:   repo.ID,
		Repo:     repo,
		Index:    1,
		Title:    "un pr",
		PosterID: user.ID,
		Poster:   user,
		IsPull:   true,
	}
	require.NoError(t, db.Create(issue).Error)
	require.NoError(t, db.Create(&PullRequest{
		IssueID:      issue.ID,
		Index:        1,
		HeadRepoID:   repo.ID,
		BaseRepoID:   repo.ID,
		HeadUserName: user.Name,
		HeadBranch:   "feature",
		BaseBranch:   "master",
	}).Error)

	issues, err := Issues(&IssuesOptions{RepoID: repo.ID, IsPull: true, Page: 1})
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.NotNil(t, issues[0].PullRequest)

	var formatted *apiv2types.PullRequestResponse
	require.NotPanics(t, func() {
		formatted = issues[0].PullRequest.APIFormat()
	})
	assert.Equal(t, "un pr", formatted.Title)
	assert.Equal(t, "master", formatted.BaseBranch)
	require.NotNil(t, formatted.BaseRepo)
}
