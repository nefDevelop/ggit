package database

import (
	"context"
	"testing"

	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsUsernameUsedReportsFoundNames covers the ordinary cases.
func TestIsUsernameUsedReportsFoundNames(t *testing.T) {
	db := dbtest.NewDB(t, "IsUsernameUsed", new(User))
	SetHandle(db)
	ctx := context.Background()
	store := NewUsersStore(db)

	created, err := Handle.Users().Create(ctx, "taken-name", "taken@example.com", CreateUserOptions{})
	require.NoError(t, err)

	assert.True(t, store.IsUsernameUsed(ctx, "taken-name", 0), "un nombre existente debe estar en uso")
	assert.True(t, store.IsUsernameUsed(ctx, "TAKEN-NAME", 0), "la comparacion debe ser insensible a mayusculas")
	assert.False(t, store.IsUsernameUsed(ctx, "free-name", 0), "un nombre libre no debe estar en uso")
	assert.False(t, store.IsUsernameUsed(ctx, "", 0), "un nombre vacio no debe estar en uso")
	assert.False(t, store.IsUsernameUsed(ctx, "taken-name", created.ID),
		"excluir el propio id debe reportar el nombre como libre")
}

// TestIsUsernameUsedSurvivesAFailedLookup is the regression guard.
//
// IsUsernameUsed used to end in "Error != gorm.ErrRecordNotFound", so any failure
// other than "no rows" counted as "the name is taken". A missing table, a locked
// SQLite file or a dropped connection all became "User already exists": a
// duplicate-name message for something that is not a duplicate, hiding the real
// cause.
//
// A failed lookup must report "not in use" instead. That is safe because
// user.lower_name carries a unique index, so a real duplicate is still rejected
// by the INSERT, with an accurate error.
func TestIsUsernameUsedSurvivesAFailedLookup(t *testing.T) {
	// Deliberately no tables: every lookup against "user" fails.
	db := dbtest.NewDB(t, "IsUsernameUsedNoSchema")
	SetHandle(db)
	ctx := context.Background()
	store := NewUsersStore(db)

	assert.NotPanics(t, func() {
		assert.False(t, store.IsUsernameUsed(ctx, "cualquier", 0),
			"una consulta fallida no debe informarse como nombre ocupado")
	}, "una consulta fallida no debe provocar un panic")

	// With the table in place the duplicate is caught, and it is caught by the
	// unique index rather than by the pre-check.
	require.NoError(t, db.AutoMigrate(new(User)))

	created, err := Handle.Users().Create(ctx, "first", "first@example.com", CreateUserOptions{})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	_, err = Handle.Users().Create(ctx, "first", "second@example.com", CreateUserOptions{})
	require.Error(t, err, "el indice unico debe seguir rechazando el duplicado")
	assert.True(t, IsErrUserAlreadyExist(err), "el error debe seguir siendo de nombre duplicado")
}
