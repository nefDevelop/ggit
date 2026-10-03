package gitx

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. Setup: Crear un repositorio real en una carpeta temporal
	tmpDir, err := os.MkdirTemp("", "ggit-git-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repoPath := filepath.Join(tmpDir, "repo.git")
	err = os.MkdirAll(repoPath, 0755)
	require.NoError(t, err)

	// Inicializar repo real usando el comando git del sistema
	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = repoPath
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(out))
	}

	runGit("init")
	runGit("config", "user.email", "test@ggit.io")
	runGit("config", "user.name", "GGit Tester")

	// Crear un archivo y un commit
	err = os.WriteFile(filepath.Join(repoPath, "README.md"), []byte("# Test Repo"), 0644)
	require.NoError(t, err)
	runGit("add", "README.md")
	runGit("commit", "-m", "Initial commit")

	// 2. Ejecutar tests usando nuestra abstracción 'Module' (la implementación real)
	realModule := module{}

	t.Run("Log returns real commits", func(t *testing.T) {
		commits, err := realModule.Log(repoPath, "master")
		if err != nil {
			// Intentar con 'main' si 'master' falla por configuración local de git
			commits, err = realModule.Log(repoPath, "main")
		}
		require.NoError(t, err)
		assert.Len(t, commits, 1)
		assert.Equal(t, "Initial commit", commits[0].Summary())
	})

	t.Run("Tags management", func(t *testing.T) {
		// Crear un tag real
		runGit("tag", "-a", "v1.0.0", "-m", "Release v1.0.0")

		tags, err := realModule.RepoTags(repoPath)
		require.NoError(t, err)
		assert.Contains(t, tags, "v1.0.0")
	})

	t.Run("DiffNameOnly finds changes", func(t *testing.T) {
		// Crear un cambio
		err = os.WriteFile(filepath.Join(repoPath, "TODO.md"), []byte("- [ ] Fix bugs"), 0644)
		require.NoError(t, err)
		runGit("add", "TODO.md")
		runGit("commit", "-m", "Add TODO")

		// Comparar HEAD con el commit anterior (HEAD^)
		diffs, err := realModule.DiffNameOnly(repoPath, "HEAD^", "HEAD")
		require.NoError(t, err)
		assert.Contains(t, diffs, "TODO.md")
		assert.NotContains(t, diffs, "README.md")
	})

	t.Run("Remotes management", func(t *testing.T) {
		remoteURL := "https://github.com/nef734/ggit.git"
		err := realModule.RemoteAdd(repoPath, "origin", remoteURL)
		require.NoError(t, err)

		remotes, err := realModule.Remotes(repoPath)
		require.NoError(t, err)
		assert.Contains(t, remotes, "origin")

		err = realModule.RemoteRemove(repoPath, "origin")
		require.NoError(t, err)

		remotes, _ = realModule.Remotes(repoPath)
		assert.NotContains(t, remotes, "origin")
	})

	t.Run("MergeBase and PullRequestMeta with Fork simulation", func(t *testing.T) {
		// 1. Crear Repo Base
		basePath := filepath.Join(tmpDir, "base.git")
		err = os.MkdirAll(basePath, 0755)
		require.NoError(t, err)

		runGit := func(path string, args ...string) string {
			cmd := exec.Command("git", args...)
			cmd.Dir = path
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "git %v failed: %s", args, string(out))
			return string(out)
		}

		runGit(basePath, "init")
		runGit(basePath, "config", "user.email", "base@ggit.io")
		runGit(basePath, "config", "user.name", "Base User")
		os.WriteFile(filepath.Join(basePath, "base.txt"), []byte("base"), 0644)
		runGit(basePath, "add", ".")
		runGit(basePath, "commit", "-m", "Base commit")

		// Guardar el hash del merge base esperado
		expectedMergeBase := runGit(basePath, "rev-parse", "HEAD")

		// 2. Simular un Fork (Clonar base a head)
		headPath := filepath.Join(tmpDir, "head.git")
		runGit(tmpDir, "clone", basePath, headPath)
		runGit(headPath, "config", "user.email", "fork@ggit.io")
		runGit(headPath, "config", "user.name", "Fork User")

		// Crear una rama de feature en el fork
		runGit(headPath, "checkout", "-b", "feature")
		os.WriteFile(filepath.Join(headPath, "feature.txt"), []byte("feature content"), 0644)
		runGit(headPath, "add", ".")
		runGit(headPath, "commit", "-m", "Feature commit 1")

		os.WriteFile(filepath.Join(headPath, "feature2.txt"), []byte("more content"), 0644)
		runGit(headPath, "add", ".")
		runGit(headPath, "commit", "-m", "Feature commit 2")

		// 3. Verificar MergeBase
		// Para esto necesitamos que el fork tenga conocimiento de la rama del base
		// PullRequestMeta lo hace automáticamente añadiendo un remoto temporal.
		// Aquí probamos primero la lógica interna de PullRequestMeta.

		meta, err := realModule.PullRequestMeta(headPath, basePath, "feature", "master")
		if err != nil {
			// Fallback if default branch is main
			meta, err = realModule.PullRequestMeta(headPath, basePath, "feature", "main")
		}
		require.NoError(t, err)

		assert.Contains(t, meta.MergeBase, expectedMergeBase[:10]) // Hash corto suele bastar
		assert.Len(t, meta.Commits, 2, "Debería haber 2 commits en el PR")
		assert.Equal(t, 2, meta.NumFiles, "Deberían haber cambiado 2 archivos (feature.txt y feature2.txt)")
	})

	t.Run("ListTagsAfter with real pagination", func(t *testing.T) {
		repoPath := filepath.Join(tmpDir, "repo.git")

		// git-module lists tags with "--sort=-creatordate", so the creation date is
		// what decides the order. Lightweight tags all inherit the commit date, and
		// creating five of them back to back gives every one the same creatordate:
		// git then returns them in an order that is stable per repository but not
		// across repositories, which is why this test failed from time to time.
		// Annotated tags with an explicit committer date pin the ordering.
		tags := []string{"v1.0.0", "v1.1.0", "v1.2.0", "v2.0.0", "v2.1.0"}

		// An earlier subtest shares this repository and leaves these tag names
		// behind, so they are cleared first rather than silently reused.
		for _, tag := range tags {
			del := exec.Command("git", "tag", "-d", tag)
			del.Dir = repoPath
			_ = del.Run()
		}

		base := time.Now().Add(-time.Hour)
		for i, tag := range tags {
			cmd := exec.Command("git", "tag", "-a", tag, "-m", tag)
			cmd.Dir = repoPath
			cmd.Env = append(os.Environ(),
				"GIT_COMMITTER_DATE="+base.Add(time.Duration(i)*time.Minute).Format(time.RFC3339),
				"GIT_AUTHOR_DATE="+base.Add(time.Duration(i)*time.Minute).Format(time.RFC3339),
			)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git tag %s: %v\n%s", tag, err, out)
			}
		}

		// Probar paginación: Dame 2 tags después de v2.1.0 (en orden descendente)
		page, err := realModule.ListTagsAfter(repoPath, "", 2)
		require.NoError(t, err)
		assert.Len(t, page.Tags, 2)
		assert.True(t, page.HasNext)

		// El orden real depende de cómo git-module ejecute 'git tag'.
		// Generalmente es alfabético descendente en Gogs.
		assert.Equal(t, "v2.1.0", page.Tags[0])
		assert.Equal(t, "v2.0.0", page.Tags[1])

		// Siguiente página
		page2, err := realModule.ListTagsAfter(repoPath, "v2.0.0", 2)
		require.NoError(t, err)
		assert.Len(t, page2.Tags, 2)
		assert.Equal(t, "v1.2.0", page2.Tags[0])
		assert.Equal(t, "v1.1.0", page2.Tags[1])
	})
}
