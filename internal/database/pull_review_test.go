package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nef734/ggit/internal/dbtest"
)

func TestPullReviews(t *testing.T) {
	db := dbtest.NewDB(t, "PullReviewsStore", new(PullReview), new(PullReviewComment), new(User))
	ctx := context.Background()
	store := newPullReviewsStore(db)

	t.Run("Superación: Crear y listar aprobaciones", func(t *testing.T) {
		issueID := int64(100)
		review := &PullReview{
			IssueID:    issueID,
			ReviewerID: 1,
			Type:       PullReviewTypeApprove,
		}
		err := store.Create(ctx, review)
		assert.NoError(t, err)

		reviews, err := store.ListByIssueID(ctx, issueID)
		assert.NoError(t, err)
		assert.Len(t, reviews, 1)
		assert.Equal(t, PullReviewTypeApprove, reviews[0].Type)

		stats, err := store.GetReviewStatsByIssueID(ctx, issueID)
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
