package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/dbtest"
)

func TestPullReviews(t *testing.T) {
	db := dbtest.NewDB(t, "PullReviewsStore", new(PullReview), new(PullReviewComment), new(User), new(PullRequest), new(Issue), new(Repository), new(IssueLabel), new(Label), new(Attachment), new(Comment), new(Milestone))
	SetHandle(db)
	ctx := context.Background()
	store := newPullReviewsStore(db)

	t.Run("Superación: Crear y listar aprobaciones", func(t *testing.T) {
		// Need to create a user, repository, issue and pull request record first
		user := &User{
			LowerName: "test-user",
			Name:      "test-user",
			Email:     "test@example.com",
		}
		require.NoError(t, db.Create(user).Error)

		repo := &Repository{
			OwnerID:   user.ID,
			LowerName: "test-repo",
			Name:      "test-repo",
		}
		require.NoError(t, db.Create(repo).Error)

		issue := &Issue{
			RepoID: repo.ID,
			Index:  1,
			Title:  "test-issue",
		}
		require.NoError(t, db.Create(issue).Error)

		err := db.Create(&PullRequest{
			IssueID:    issue.ID,
			BaseRepoID: repo.ID,
			HeadRepoID: repo.ID,
		}).Error
		require.NoError(t, err)

		review := &PullReview{
			IssueID:    issue.ID,
			ReviewerID: user.ID,
			Type:       PullReviewTypeApprove,
		}
		err = store.Create(ctx, review)
		assert.NoError(t, err)

		reviews, err := store.ListByIssueID(ctx, issue.ID)
		assert.NoError(t, err)
		assert.Len(t, reviews, 1)
		assert.Equal(t, PullReviewTypeApprove, reviews[0].Type)

		stats, err := store.GetReviewStatsByIssueID(ctx, issue.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats.Approvals)
		assert.Equal(t, int64(0), stats.Rejections)
	})

	t.Run("Superación: Comentarios en línea (Inline Comments)", func(t *testing.T) {
		commentStore := newPullReviewCommentsStore(db)
		reviewID := int64(1)
		
		comment := &PullReviewComment{
			ReviewID: reviewID,
			IssueID:  100,
			Content:  "Excelente cambio en esta línea",
			Line:     42,
			TreePath: "main.go",
		}
		err := commentStore.Create(ctx, comment)
		assert.NoError(t, err)

		comments, err := commentStore.ListByReviewID(ctx, reviewID)
		assert.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, "main.go", comments[0].TreePath)
	})

	t.Run("Seguridad: Integridad del revisor", func(t *testing.T) {
		// Verificar que no se puede crear una revisión sin tipo válido o sin issue asociado
		badReview := &PullReview{
			IssueID: 0, // Inválido
			Type:    999, // Tipo inexistente
		}
		err := store.Create(ctx, badReview)
		assert.Error(t, err, "Se debe rechazar una revisión con IssueID inválido o tipo inexistente")
	})
}
