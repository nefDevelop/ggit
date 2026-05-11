package auth

import (
	"context"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	db := dbtest.NewDB(t, "AuthService_Register", database.Tables...)
	database.SetHandle(db)
	Init(database.Handle)

	ctx := context.Background()

	t.Run("first user becomes admin and activated", func(t *testing.T) {
		// Mock config
		conf.Auth.RequireEmailConfirmation = true // Should be ignored for first user

		user, err := Auth.Register(ctx, SignupOptions{
			Username: "myadmin",
			Email:    "admin@example.com",
			Password: "password",
		})
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, user.IsAdmin)
		assert.True(t, user.IsActive)
	})

	t.Run("second user is not admin and not activated if required", func(t *testing.T) {
		conf.Auth.RequireEmailConfirmation = true

		user, err := Auth.Register(ctx, SignupOptions{
			Username: "user2",
			Email:    "user2@example.com",
			Password: "password",
		})
		require.NoError(t, err)
		assert.False(t, user.IsAdmin)
		assert.False(t, user.IsActive)
	})

	t.Run("third user is activated if not required", func(t *testing.T) {
		conf.Auth.RequireEmailConfirmation = false

		user, err := Auth.Register(ctx, SignupOptions{
			Username: "user3",
			Email:    "user3@example.com",
			Password: "password",
		})
		require.NoError(t, err)
		assert.False(t, user.IsAdmin)
		assert.True(t, user.IsActive)
	})
}
