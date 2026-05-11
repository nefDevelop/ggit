package database

import (
	"context"
	"testing"
	"time"

	"github.com/nef734/ggit/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserAndCollaborativeRepositories(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	db := newTestDB(t, "get_user_and_collaborative_repositories")
	Handle = &DB{db: db}

	ctx := context.Background()
	usersStore := NewUsersStore(db)

	// Create users
	alice, err := usersStore.Create(ctx, "alice", "alice@example.com", CreateUserOptions{})
	require.NoError(t, err)
	bob, err := usersStore.Create(ctx, "bob", "bob@example.com", CreateUserOptions{})
	require.NoError(t, err)

	// Create repositories
	reposStore := NewRepositoriesStore(db)
	repo1, err := reposStore.Create(ctx, alice.ID, CreateRepoOptions{Name: "repo1"})
	require.NoError(t, err)
	repo2, err := reposStore.Create(ctx, bob.ID, CreateRepoOptions{Name: "repo2"})
	require.NoError(t, err)
	_, err = reposStore.Create(ctx, bob.ID, CreateRepoOptions{Name: "repo3"})
	require.NoError(t, err)

	// alice owns repo1.
	// bob owns repo2 and repo3.

	// Make alice a collaborator of repo2 via Access table (e.g. simulate team/explicit access)
	permissionsStore := newPermissionsStore(db)
	err = permissionsStore.SetRepoPerms(ctx, repo2.ID, map[int64]domain.AccessMode{alice.ID: domain.AccessModeWrite})
	require.NoError(t, err)

	// Test GetUserAndCollaborativeRepositories for alice
	got, err := GetUserAndCollaborativeRepositories(alice.ID)
	require.NoError(t, err)
	assert.Len(t, got, 2)

	foundRepo1 := false
	foundRepo2 := false
	for _, r := range got {
		if r.ID == repo1.ID {
			foundRepo1 = true
		}
		if r.ID == repo2.ID {
			foundRepo2 = true
		}
	}
	assert.True(t, foundRepo1, "Should find owned repo1")
	assert.True(t, foundRepo2, "Should find collaborative repo2 (via Access table)")
}

func TestRepositoryList_SortByUpdatedDesc(t *testing.T) {
	now := time.Now().Unix()
	repos := RepositoryList{
		{ID: 1, UpdatedUnix: now - 100},
		{ID: 2, UpdatedUnix: now},
		{ID: 3, UpdatedUnix: now - 50},
	}

	repos.SortByUpdatedDesc()

	assert.Equal(t, int64(2), repos[0].ID)
	assert.Equal(t, int64(3), repos[1].ID)
	assert.Equal(t, int64(1), repos[2].ID)
}
