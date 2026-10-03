package database

import (
	"os"
	"path/filepath"
	"testing"
)

// The console commands call InitLogging(true), which returns before creating
// the log root, and then reach SetEngine with HookMode false. If SetEngine does
// not create the directory itself, "ggit admin create-user" aborts on a fresh
// install and the documented bootstrap command cannot bootstrap.
func TestEnsureLogDirCreatesMissingDirectories(t *testing.T) {
	root := filepath.Join(t.TempDir(), "data", "log")

	if err := ensureLogDir(filepath.Join(root, "gorm.log")); err != nil {
		t.Fatalf("ensureLogDir on a missing tree: %v", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("log root was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("log root is not a directory")
	}
}

func TestEnsureLogDirIsIdempotent(t *testing.T) {
	root := filepath.Join(t.TempDir(), "log")
	path := filepath.Join(root, "hooks", "gorm.log")

	if err := ensureLogDir(path); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if err := ensureLogDir(path); err != nil {
		t.Fatalf("second call on an existing directory: %v", err)
	}
}
