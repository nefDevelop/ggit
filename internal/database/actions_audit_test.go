package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbtest"
)

func TestActionsStore_ListAuditLog(t *testing.T) {
	// Asegurar que la configuración de paginación esté inicializada
	conf.UI.Admin.NoticePagingNum = 25

	db := dbtest.NewDB(t, "ListAuditLog", new(Action), new(User))
	ctx := context.Background()
	store := &ActionsStore{db: db}

	// Setup users
	admin := &User{ID: 1, Name: "admin", LowerName: "admin"}
	user := &User{ID: 2, Name: "user", LowerName: "user"}
	require.NoError(t, db.Create([]*User{admin, user}).Error)

	t.Run("Superación: Listar logs de auditoría (Acciones donde UserID == ActUserID)", func(t *testing.T) {
		// Crear acciones de auditoría (donde el usuario es el mismo actor, p.ej. cambios de config global)
		actions := []*Action{
			{UserID: admin.ID, ActUserID: admin.ID, ActUserName: admin.Name, OpType: ActionUpdateConfig, Content: "auth|ldap"},
			{UserID: admin.ID, ActUserID: admin.ID, ActUserName: admin.Name, OpType: ActionUpdateConfig, Content: "server|offline_mode"},
		}
		require.NoError(t, db.Create(actions).Error)

		// Crear acción normal de feed (UserID != ActUserID, p.ej. una notificación de commit para otro usuario)
		require.NoError(t, db.Create(&Action{UserID: user.ID, ActUserID: admin.ID, ActUserName: admin.Name, OpType: ActionCommitRepo}).Error)

		logs, err := store.ListAuditLog(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, logs, 2, "Solo deben devolverse las acciones de auditoría (UserID == ActUserID)")
		assert.Equal(t, int(ActionUpdateConfig), logs[0].OpType)
	})

	t.Run("Fallo: Paginación fuera de rango", func(t *testing.T) {
		logs, err := store.ListAuditLog(ctx, 99)
		assert.NoError(t, err)
		assert.Empty(t, logs, "Una página sin resultados debe devolver lista vacía, no error")
		
		logsZero, err := store.ListAuditLog(ctx, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, logsZero, "Página 0 debe reajustarse automáticamente a página 1")
	})

	t.Run("Seguridad: Aislamiento de contenido", func(t *testing.T) {
		// Verificar que el log no expone información sensible de forma inadvertida
		// (p.ej. si el contenido tuviera un formato específico)
		logs, _ := store.ListAuditLog(ctx, 1)
		for _, log := range logs {
			assert.NotContains(t, log.Content, "PASSWORD", "Los logs de auditoría nunca deben persistir secretos en texto plano")
		}
	})
}
