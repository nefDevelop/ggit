package database

import (
	"context"
	"testing"

	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPikoCIReposStore(t *testing.T) {
	db := dbtest.NewDB(t, "PikoCIReposStore", new(PikoCIRepo))
	ctx := context.Background()
	store := &PikoCIRepos{db: db}

	_, err := store.GetByRepoID(ctx, 1)
	assert.ErrorIs(t, err, ErrPikoCIRepoNotExist)

	require.NoError(t, store.Upsert(ctx, &PikoCIRepo{
		RepoID:       1,
		PipelineName: "pi-ggit",
		WebhookToken: "token-1",
	}))

	link, err := store.GetByRepoID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "pi-ggit", link.PipelineName)
	assert.Equal(t, "token-1", link.WebhookToken)
	assert.NotZero(t, link.CreatedAt)

	count, err := store.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	links, err := store.List(ctx, 0)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Equal(t, int64(1), links[0].RepoID)

	// Un segundo Upsert actualiza el vínculo en lugar de duplicarlo.
	require.NoError(t, store.Upsert(ctx, &PikoCIRepo{
		RepoID:       1,
		PipelineName: "pi-ggit-2",
		WebhookToken: "token-2",
	}))

	link, err = store.GetByRepoID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "pi-ggit-2", link.PipelineName)
	assert.Equal(t, "token-2", link.WebhookToken)

	require.NoError(t, store.Delete(ctx, 1))
	_, err = store.GetByRepoID(ctx, 1)
	assert.ErrorIs(t, err, ErrPikoCIRepoNotExist)

	count, err = store.Count(ctx)
	require.NoError(t, err)
	assert.Zero(t, count)

	links, err = store.List(ctx, 0)
	require.NoError(t, err)
	assert.Empty(t, links)
}
