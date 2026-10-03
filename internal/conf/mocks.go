package conf

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func SetMockApp(t *testing.T, opts AppOpts) {
	before := App
	App = opts
	t.Cleanup(func() {
		App = before
	})
}

func SetMockAuth(t *testing.T, otps AuthOpts) {
	before := Auth
	Auth = otps
	t.Cleanup(func() {
		Auth = before
	})
}

var mockServer sync.Mutex

func SetMockServer(t *testing.T, opts ServerOpts) {
	mockServer.Lock()
	before := Server
	Server = opts
	t.Cleanup(func() {
		Server = before
		mockServer.Unlock()
	})
}

var mockSSH sync.Mutex

func SetMockSSH(t *testing.T, opts SSHOpts) {
	mockSSH.Lock()
	before := SSH
	SSH = opts
	t.Cleanup(func() {
		SSH = before
		mockSSH.Unlock()
	})
}

var mockRepository sync.Mutex

func SetMockRepository(t *testing.T, opts RepositoryOpts) {
	mockRepository.Lock()
	before := Repository
	Repository = opts
	t.Cleanup(func() {
		Repository = before
		mockRepository.Unlock()
	})
}

func SetMockUI(t *testing.T, opts UIOpts) {
	before := UI
	UI = opts
	t.Cleanup(func() {
		UI = before
	})
}

var mockPicture sync.Mutex

func SetMockPicture(t *testing.T, opts PictureOpts) {
	mockPicture.Lock()
	before := Picture
	Picture = opts
	t.Cleanup(func() {
		Picture = before
		mockPicture.Unlock()
	})
}

// UseTempRepositoryRoot points Repository.Root at a directory outside the source
// tree and returns its path.
//
// It exists because Repository.Root is empty until conf.Init runs, and test
// binaries do not call it. repox.UserPath joins the user name straight onto that
// empty root, so every repository fixture created by a test lands relative to the
// package directory, inside the source tree. Running the suite then rewrites
// tracked hook files with the build path and leaves stray git objects behind,
// which is enough to corrupt an unrelated commit.
//
// Meant for TestMain, where no *testing.T is available yet; call RemoveTempRoot
// from TestMain's exit path if the contents need to survive a failed run.
func UseTempRepositoryRoot() (string, error) {
	dir, err := os.MkdirTemp("", "ggit-test-repositories-")
	if err != nil {
		return "", fmt.Errorf("create temporary repository root: %w", err)
	}
	Repository.Root = dir
	return dir, nil
}
