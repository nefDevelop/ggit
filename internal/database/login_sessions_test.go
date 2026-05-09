package database

import (
	"context"
	"testing"
	"time"

	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
)

func TestLoginSessionsStore(t *testing.T) {
	db := dbtest.NewDB(t, "LoginSessionsStore", new(LoginSession))
	ctx := context.Background()
	store := newLoginSessionsStore(db)

	t.Run("Superación: Flujo completo de sesión", func(t *testing.T) {
		sid := "secure-session-id-123"
		userID := int64(42)

		// Create
		sess, err := store.Create(ctx, userID, sid, "Mozilla/5.0", "127.0.0.1")
		assert.NoError(t, err)
		assert.Equal(t, sid, sess.SID)

		// Get
		got, err := store.GetBySID(ctx, sid)
		assert.NoError(t, err)
		assert.Equal(t, userID, got.UserID)

		// Touch
		oldUpdate := got.UpdatedUnix
		time.Sleep(1 * time.Second)
		err = store.TouchBySID(ctx, sid)
		assert.NoError(t, err)

		updated, _ := store.GetBySID(ctx, sid)
		assert.True(t, updated.UpdatedUnix > oldUpdate, "UpdatedUnix debe incrementarse al hacer Touch")

		// Delete
		err = store.DeleteBySID(ctx, sid)
		assert.NoError(t, err)

		_, err = store.GetBySID(ctx, sid)
		assert.Error(t, err, "La sesión no debe existir tras ser borrada")
	})

	t.Run("Fallo: SID inexistente", func(t *testing.T) {
		_, err := store.GetBySID(ctx, "non-existent")
		assert.Error(t, err)
	})

	t.Run("Seguridad: Colisión de SID y Secuestro", func(t *testing.T) {
		sid := "unique-sid"
		_, err := store.Create(ctx, 1, sid, "ua", "ip")
		assert.NoError(t, err)

		// Intento de duplicar SID (Debe fallar por UniqueIndex)
		_, err = store.Create(ctx, 2, sid, "ua2", "ip2")
		assert.Error(t, err, "El sistema debe impedir la creación de sesiones con SIDs duplicados")

		// Intento de borrar sesión de otro usuario (Security Breach)
		// DeleteByIDAndUserID protege este caso
		err = store.DeleteByIDAndUserID(ctx, 1, 999) // ID 1, pero UserID equivocado
		assert.NoError(t, err)

		// Verificar que la sesión SIGUE existiendo (no fue borrada por el impostor)
		exists, _ := store.GetBySID(ctx, sid)
		assert.NotNil(t, exists, "Un usuario no debe poder borrar las sesiones de otros mediante manipulación de ID")
	})
}
