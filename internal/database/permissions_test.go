package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/nef734/ggit/internal/domain"
)

func TestPerms(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()

	ctx := context.Background()
	s := &PermissionsStore{
		db: newTestDB(t, "PermissionsStore"),
	}

	for _, tc := range []struct {
		name string
		test func(t *testing.T, ctx context.Context, s *PermissionsStore)
	}{
		{"AccessMode", permsAccessMode},
		{"Authorize", permsAuthorize},
		{"SetRepoPerms", permsSetRepoPerms},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := clearTables(t, s.db)
				require.NoError(t, err)
			})
			tc.test(t, ctx, s)
		})
		if t.Failed() {
			break
		}
	}
}

func permsAccessMode(t *testing.T, ctx context.Context, s *PermissionsStore) {
	// Set up permissions
	err := s.SetRepoPerms(ctx, 1,
		map[int64]domain.AccessMode{
			2: domain.AccessModeWrite,
			3: domain.AccessModeAdmin,
		},
	)
	require.NoError(t, err)
	err = s.SetRepoPerms(ctx, 2,
		map[int64]domain.AccessMode{
			1: domain.AccessModeRead,
		},
	)
	require.NoError(t, err)

	publicRepoID := int64(1)
	publicRepoOpts := domain.AccessModeOptions{
		OwnerID: 98,
	}

	privateRepoID := int64(2)
	privateRepoOpts := domain.AccessModeOptions{
		OwnerID: 99,
		Private: true,
	}

	tests := []struct {
		name           string
		userID         int64
		repoID         int64
		opts           domain.AccessModeOptions
		wantAccessMode domain.AccessMode
	}{
		{
			name:           "nil repository",
			wantAccessMode: domain.AccessModeNone,
		},

		{
			name:           "anonymous user has read access to public repository",
			repoID:         publicRepoID,
			opts:           publicRepoOpts,
			wantAccessMode: domain.AccessModeRead,
		},
		{
			name:           "anonymous user has no access to private repository",
			repoID:         privateRepoID,
			opts:           privateRepoOpts,
			wantAccessMode: domain.AccessModeNone,
		},

		{
			name:           "user is the owner",
			userID:         98,
			repoID:         publicRepoID,
			opts:           publicRepoOpts,
			wantAccessMode: domain.AccessModeOwner,
		},
		{
			name:           "user 1 has read access to public repo",
			userID:         1,
			repoID:         publicRepoID,
			opts:           publicRepoOpts,
			wantAccessMode: domain.AccessModeRead,
		},
		{
			name:           "user 2 has write access to public repo",
			userID:         2,
			repoID:         publicRepoID,
			opts:           publicRepoOpts,
			wantAccessMode: domain.AccessModeWrite,
		},
		{
			name:           "user 3 has admin access to public repo",
			userID:         3,
			repoID:         publicRepoID,
			opts:           publicRepoOpts,
			wantAccessMode: domain.AccessModeAdmin,
		},

		{
			name:           "user 1 has read access to private repo",
			userID:         1,
			repoID:         privateRepoID,
			opts:           privateRepoOpts,
			wantAccessMode: domain.AccessModeRead,
		},
		{
			name:           "user 2 has no access to private repo",
			userID:         2,
			repoID:         privateRepoID,
			opts:           privateRepoOpts,
			wantAccessMode: domain.AccessModeNone,
		},
		{
			name:           "user 3 has no access to private repo",
			userID:         3,
			repoID:         privateRepoID,
			opts:           privateRepoOpts,
			wantAccessMode: domain.AccessModeNone,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mode := s.AccessMode(ctx, test.userID, test.repoID, test.opts)
			assert.Equal(t, test.wantAccessMode, mode)
		})
	}
}

func permsAuthorize(t *testing.T, ctx context.Context, s *PermissionsStore) {
	// Set up permissions
	err := s.SetRepoPerms(ctx, 1,
		map[int64]domain.AccessMode{
			1: domain.AccessModeRead,
			2: domain.AccessModeWrite,
			3: domain.AccessModeAdmin,
		},
	)
	require.NoError(t, err)

	repo := &Repository{
		ID:      1,
		OwnerID: 98,
	}

	tests := []struct {
		name           string
		userID         int64
		desired        domain.AccessMode
		wantAuthorized bool
	}{
		{
			name:           "user 1 has read and wants read",
			userID:         1,
			desired:        domain.AccessModeRead,
			wantAuthorized: true,
		},
		{
			name:           "user 1 has read and wants write",
			userID:         1,
			desired:        domain.AccessModeWrite,
			wantAuthorized: false,
		},

		{
			name:           "user 2 has write and wants read",
			userID:         2,
			desired:        domain.AccessModeRead,
			wantAuthorized: true,
		},
		{
			name:           "user 2 has write and wants write",
			userID:         2,
			desired:        domain.AccessModeWrite,
			wantAuthorized: true,
		},
		{
			name:           "user 2 has write and wants admin",
			userID:         2,
			desired:        domain.AccessModeAdmin,
			wantAuthorized: false,
		},

		{
			name:           "user 3 has admin and wants read",
			userID:         3,
			desired:        domain.AccessModeRead,
			wantAuthorized: true,
		},
		{
			name:           "user 3 has admin and wants write",
			userID:         3,
			desired:        domain.AccessModeWrite,
			wantAuthorized: true,
		},
		{
			name:           "user 3 has admin and wants admin",
			userID:         3,
			desired:        domain.AccessModeAdmin,
			wantAuthorized: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authorized := s.Authorize(ctx, test.userID, repo.ID, test.desired,
				domain.AccessModeOptions{
					OwnerID: repo.OwnerID,
					Private: repo.IsPrivate,
				},
			)
			assert.Equal(t, test.wantAuthorized, authorized)
		})
	}
}

func permsSetRepoPerms(t *testing.T, ctx context.Context, s *PermissionsStore) {
	for _, update := range []struct {
		repoID    int64
		accessMap map[int64]domain.AccessMode
	}{
		{
			repoID: 1,
			accessMap: map[int64]domain.AccessMode{
				1: domain.AccessModeWrite,
				2: domain.AccessModeWrite,
				3: domain.AccessModeAdmin,
				4: domain.AccessModeWrite,
			},
		},
		{
			repoID: 2,
			accessMap: map[int64]domain.AccessMode{
				1: domain.AccessModeWrite,
				2: domain.AccessModeRead,
				4: domain.AccessModeWrite,
				5: domain.AccessModeWrite,
			},
		},
		{
			repoID: 1,
			accessMap: map[int64]domain.AccessMode{
				2: domain.AccessModeWrite,
				3: domain.AccessModeAdmin,
			},
		},
		{
			repoID: 2,
			accessMap: map[int64]domain.AccessMode{
				1: domain.AccessModeWrite,
				2: domain.AccessModeRead,
				5: domain.AccessModeWrite,
			},
		},
	} {
		err := s.SetRepoPerms(ctx, update.repoID, update.accessMap)
		if err != nil {
			t.Fatal(err)
		}
	}

	var accesses []*Access
	err := s.db.Order("user_id, repo_id").Find(&accesses).Error
	require.NoError(t, err)

	// Ignore ID fields
	for _, a := range accesses {
		a.ID = 0
	}

	wantAccesses := []*Access{
		{UserID: 1, RepoID: 2, Mode: domain.AccessModeWrite},
		{UserID: 2, RepoID: 1, Mode: domain.AccessModeWrite},
		{UserID: 2, RepoID: 2, Mode: domain.AccessModeRead},
		{UserID: 3, RepoID: 1, Mode: domain.AccessModeAdmin},
		{UserID: 5, RepoID: 2, Mode: domain.AccessModeWrite},
	}
	assert.Equal(t, wantAccesses, accesses)
}
