package conf

import (
	"path/filepath"
	"testing"
)

// The generated key must not land in a file git tracks. The default location is
// under the data directory and ignored; anything else usually means -c pointed
// at a tracked config, which is how a generated key got staged for commit
// through "make test-e2e" handing this code conf/app.test.ini.
func TestIsDefaultCustomConf(t *testing.T) {
	saved := CustomConf
	t.Cleanup(func() { CustomConf = saved })

	CustomConf = filepath.Join(CustomDir(), "conf", "app.ini")
	if !isDefaultCustomConf() {
		t.Errorf("%q should be recognised as the default location", CustomConf)
	}

	CustomConf = filepath.Join(t.TempDir(), "app.ini")
	if isDefaultCustomConf() {
		t.Error("a config outside the data directory must not be treated as default")
	}
}
