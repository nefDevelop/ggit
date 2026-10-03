package repo

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/nef734/ggit/internal/repox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryService_Create(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	db := dbtest.NewDB(t, "RepositoryService_Create", database.Tables...)
	database.SetHandle(db)
	Init(database.Handle)

	ctx := context.Background()

	// Setup: Create a user to own the repo
	owner, err := database.Handle.Users().Create(ctx, "un_pro", "pro@ggit.io", database.CreateUserOptions{})
	require.NoError(t, err)
	ownerUser := database.FromDomainUser(owner)

	t.Run("successful creation with auto-init", func(t *testing.T) {
		opts := CreateOptions{
			Name:     "mi-gran-proyecto",
			AutoInit: true,
			Readme:   "Default",
		}

		repo, err := Repos.Create(ctx, ownerUser, ownerUser, opts)
		require.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, opts.Name, repo.Name)
		assert.Equal(t, ownerUser.ID, repo.OwnerID)

		// Verificar que el registro existe en DB
		dbRepo, err := database.Handle.Repositories().GetByName(ctx, ownerUser.ID, opts.Name)
		require.NoError(t, err)
		assert.Equal(t, repo.ID, dbRepo.ID)
	})

	t.Run("rollback on duplicate name", func(t *testing.T) {
		opts := CreateOptions{Name: "duplicate"}

		// First creation
		_, err := Repos.Create(ctx, ownerUser, ownerUser, opts)
		require.NoError(t, err)

		// Second creation with same name
		_, err = Repos.Create(ctx, ownerUser, ownerUser, opts)
		assert.Error(t, err, "Should fail due to unique constraint")
	})
}

func TestRepositoryService_Create_Concurrency(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Este test intenta romper el sistema creando el mismo repo simultáneamente
	db := dbtest.NewDB(t, "RepositoryService_Concurrency", database.Tables...)
	database.SetHandle(db)
	Init(database.Handle)

	ctx := context.Background()
	owner, _ := database.Handle.Users().Create(ctx, "race_condition_tester", "race@ggit.io", database.CreateUserOptions{})
	ownerUser := database.FromDomainUser(owner)

	const numGoroutines = 5
	var wg sync.WaitGroup
	errs := make(chan error, numGoroutines)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			// Intentamos crear el MISMO repo en paralelo
			_, err := Repos.Create(ctx, ownerUser, ownerUser, CreateOptions{
				Name: "concurrent-repo",
			})
			if err != nil {
				errs <- err
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	// Los perdedores tienen que fallar, no crear filas duplicadas.
	var rejected int
	for err := range errs {
		rejected++
		assert.Error(t, err)
	}
	assert.Equal(t, numGoroutines-1, rejected,
		"solo una creacion concurrent deberia prosperar; el resto debe recibir un error")

	// PageSize es obligatorio: SearchRepositoryByName termina en
	// "Limit(opts.PageSize)", asi que sin el la consulta devuelve LIMIT 0 y este
	// recuento nunca podia dar 1.
	repos, _, err := database.Handle.Repositories().List(ctx, database.ListRepoOptions{
		OwnerID:  ownerUser.ID,
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Len(t, repos, 1, "el indice unico (owner_id, lower_name) debe dejar una sola fila")

	// Y el directorio en disco tambien debe quedar en uno.
	matches, err := filepath.Glob(filepath.Join(
		repox.RepositoryPath(ownerUser.Name, "concurrent-repo") + "*"))
	assert.NoError(t, err)
	assert.Len(t, matches, 1,
		"solo debe quedar el repositorioGanador, sin staging dirs huerfanos")
}

func BenchmarkRepositoryService_Create(b *testing.B) {
	// Setup para el benchmark
	db := dbtest.NewDB(b, "Benchmark_RepoCreate", database.Tables...)
	database.SetHandle(db)
	Init(database.Handle)
	ctx := context.Background()
	owner, _ := database.Handle.Users().Create(ctx, "bench-user", "bench@ggit.io", database.CreateUserOptions{})
	ownerUser := database.FromDomainUser(owner)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Repos.Create(ctx, ownerUser, ownerUser, CreateOptions{
			Name: fmt.Sprintf("bench-repo-%d", i),
		})
	}
}
