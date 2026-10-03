package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/nef734/ggit/internal/repox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedRepoForRename creates a repository row and a bare directory for it, which
// is what updateRepository expects to find when a rename arrives.
func seedRepoForRename(t *testing.T, dbOwner string) (*User, *Repository, string) {
	t.Helper()
	ctx := context.Background()

	domainUser, err := Handle.Users().Create(ctx, dbOwner, dbOwner+"@example.com", CreateUserOptions{})
	require.NoError(t, err)
	user := FromDomainUser(domainUser)

	created, err := Handle.Repositories().Create(ctx, user.ID, CreateRepoOptions{Name: "old-name"})
	require.NoError(t, err)
	repo := FromDomainRepository(created)
	repo.Owner = user

	// The directory has to exist, otherwise there is nothing to move and the test
	// would pass without exercising the rename.
	dir := repox.RepositoryPath(user.Name, repo.Name)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "hooks"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "HEAD"), []byte("ref: refs/heads/main\n"), 0o644))

	return user, repo, dir
}

// TestUpdateRepositoryMovesTheDirectory is the regression guard for the rename.
//
// updateRepository used to save the new name and leave the directory where it
// was. Because Repository.RepoPath derives the path from the name, the row then
// pointed at a directory that did not exist, so every later git operation, archive
// and clone URL for that repository failed: renaming a repository from the settings
// page broke it.
func TestUpdateRepositoryMovesTheDirectory(t *testing.T) {
	restore := swapRepositoryRoot(t)
	defer restore()

	db := dbtest.NewDB(t, "UpdateRepositoryRename",
		new(User), new(Repository), new(Access), new(Watch), new(Issue), new(IssueLabel),
		new(Label), new(Milestone), new(Attachment), new(Comment), new(PullRequest))
	SetHandle(db)

	user, repo, oldDir := seedRepoForRename(t, "rename-owner")

	repo.Name = "new-name"
	repo.LowerName = "new-name"
	require.NoError(t, UpdateRepository(repo, false))

	newDir := repox.RepositoryPath(user.Name, "new-name")
	assert.NoDirExists(t, oldDir, "el directorio con el nombre antiguo debe desaparecer")
	assert.DirExists(t, newDir, "el directorio con el nombre nuevo debe existir")
	assert.FileExists(t, filepath.Join(newDir, "HEAD"), "el contenido debe moverse intacto")

	// And the row has to agree with the disk.
	stored, err := GetRepositoryByID(repo.ID)
	require.NoError(t, err)
	assert.Equal(t, "new-name", stored.Name)
	assert.Equal(t, newDir, repox.RepositoryPath(user.Name, stored.Name),
		"la ruta derivada del nombre guardado debe existir")
}

// TestUpdateRepositoryRenameRejectsDuplicates makes sure a rename cannot collide
// with an existing repository, which would otherwise overwrite its directory.
func TestUpdateRepositoryRenameRejectsDuplicates(t *testing.T) {
	restore := swapRepositoryRoot(t)
	defer restore()

	db := dbtest.NewDB(t, "UpdateRepositoryRenameDuplicate",
		new(User), new(Repository), new(Access), new(Watch), new(Issue), new(IssueLabel),
		new(Label), new(Milestone), new(Attachment), new(Comment), new(PullRequest))
	SetHandle(db)

	ctx := context.Background()
	user, repo, oldDir := seedRepoForRename(t, "dup-owner")

	_, err := Handle.Repositories().Create(ctx, user.ID, CreateRepoOptions{Name: "taken-name"})
	require.NoError(t, err)

	repo.Name = "taken-name"
	err = UpdateRepository(repo, false)
	require.Error(t, err, "un renombrado a un nombre existente debe rechazarse")
	assert.True(t, IsErrRepoAlreadyExist(err), "el error debe ser de repositorio ya existente: %v", err)
	assert.DirExists(t, oldDir, "el directorio original debe seguir en su sitio")
}

// TestUpdateRepositoryRenameRejectsTraversal keeps the rename behind the same
// whitelist as creation. Without it the new name decides a path, and a name with a
// separator would move a directory outside conf.Repository.Root.
func TestUpdateRepositoryRenameRejectsTraversal(t *testing.T) {
	restore := swapRepositoryRoot(t)
	defer restore()

	db := dbtest.NewDB(t, "UpdateRepositoryRenameTraversal",
		new(User), new(Repository), new(Access), new(Watch), new(Issue), new(IssueLabel),
		new(Label), new(Milestone), new(Attachment), new(Comment), new(PullRequest))
	SetHandle(db)

	_, repo, oldDir := seedRepoForRename(t, "trav-owner")
	outside := filepath.Join(filepath.Dir(oldDir), "escaped.git")
	require.NoError(t, os.MkdirAll(outside, 0o755))

	for _, hostile := range []string{"../escaped", "a/b", "/abs", `a\b`} {
		repo.Name = hostile
		err := UpdateRepository(repo, false)
		require.Error(t, err, "el renombrado a %q debe rechazarse", hostile)
		assert.DirExists(t, outside, "el directorio de fuera no debe tocarse con %q", hostile)
		assert.DirExists(t, oldDir, "el directorio original debe seguir con %q", hostile)
	}
}

// swapRepositoryRoot points conf.Repository.Root at a temporary directory for the
// duration of the test and restores it afterwards.
func swapRepositoryRoot(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()

	previous := conf.Repository.Root
	conf.Repository.Root = dir
	return func() { conf.Repository.Root = previous }
}
