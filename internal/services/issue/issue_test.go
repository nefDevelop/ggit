package issue

import (
	"context"
	"testing"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueService_CreateIssue(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	db := dbtest.NewDB(t, "IssueService_CreateIssue", database.Tables...)
	database.SetHandle(db)
	Init(database.Handle)

	ctx := context.Background()

	// Setup: Owner and Repo
	owner, _ := database.Handle.Users().Create(ctx, "issue_owner", "owner@ggit.io", database.CreateUserOptions{})
	repo, _ := database.Handle.Repositories().Create(ctx, owner.ID, database.CreateRepoOptions{Name: "issue-repo"})
	poster, _ := database.Handle.Users().Create(ctx, "poster", "poster@ggit.io", database.CreateUserOptions{})

	repoModel := database.FromDomainRepository(repo)
	repoModel.Owner = database.FromDomainUser(owner)
	posterModel := database.FromDomainUser(poster)

	t.Run("create basic issue", func(t *testing.T) {
		issue := &database.Issue{
			RepoID:   repoModel.ID,
			PosterID: posterModel.ID,
			Poster:   posterModel,
			Title:    "Un bug crítico",
			Content:  "El sistema explota al hacer click",
		}

		err := Issues.CreateIssue(ctx, repoModel, issue, CreateIssueOptions{})
		require.NoError(t, err)

		// Verificar en DB
		dbIssue, err := database.Handle.Issues().GetByIndex(ctx, repoModel.ID, 1)
		require.NoError(t, err)
		assert.Equal(t, issue.Title, dbIssue.Title)
		assert.Equal(t, int64(1), dbIssue.Index)
	})

	t.Run("issue index increments", func(t *testing.T) {
		issue := &database.Issue{
			RepoID:   repoModel.ID,
			PosterID: posterModel.ID,
			Poster:   posterModel,
			Title:    "Otro issue",
		}

		err := Issues.CreateIssue(ctx, repoModel, issue, CreateIssueOptions{})
		require.NoError(t, err)

		dbIssue, _ := database.Handle.Issues().GetByIndex(ctx, repoModel.ID, 2)
		assert.Equal(t, int64(2), dbIssue.Index)
	})

	// The existing cases above re-read from the database, so they stayed green
	// while the caller's struct kept index zero and the API answered
	// number: 0. What matters to a handler is the struct it still holds.
	t.Run("caller receives the assigned index", func(t *testing.T) {
		issue := &database.Issue{
			RepoID:   repoModel.ID,
			PosterID: posterModel.ID,
			Poster:   posterModel,
			Title:    "Tercer issue",
		}

		require.NoError(t, Issues.CreateIssue(ctx, repoModel, issue, CreateIssueOptions{}))

		assert.Equal(t, int64(3), issue.Index,
			"the struct the caller holds must carry the index assigned during creation")
		assert.NotZero(t, issue.ID, "the caller must receive the new id")
	})
}
