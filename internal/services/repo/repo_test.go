package repo

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/dbtest"
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

	// Análisis de resultados:
	// Debería haber exactamente 1 éxito y N-1 errores (por nombre duplicado)
	// Si hay más de 1 éxito, hay una condición de carrera grave.
	var successCount int
	for err := range errs {
		fmt.Printf("Expected error during race: %v\n", err)
	}
	
	// Contar cuántos hay en la DB
	repos, _, _ := database.Handle.Repositories().List(ctx, ownerUser.ID, database.ListRepoOptions{})
	successCount = len(repos)

	assert.Equal(t, 1, successCount, "SOLO uno debería haber tenido éxito en la creación concurrente")
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
