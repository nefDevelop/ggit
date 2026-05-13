package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nef734/ggit/internal/testx"
	"github.com/stretchr/testify/assert"
)

func TestIsProdMode(t *testing.T) {
	before := App.RunMode
	defer func() {
		App.RunMode = before
	}()

	tests := []struct {
		mode string
		want bool
	}{
		{mode: "dev", want: false},
		{mode: "test", want: false},

		{mode: "prod", want: true},
		{mode: "Prod", want: true},
		{mode: "PROD", want: true},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			App.RunMode = test.mode
			assert.Equal(t, test.want, IsProdMode())
		})
	}
}

func TestWorkDirHelper(_ *testing.T) {
	if !testx.WantHelperProcess() {
		return
	}

	fmt.Fprintln(os.Stdout, WorkDir())
}

func TestWorkDir(t *testing.T) {
	tests := []struct {
		env  string
		want string
	}{
		{env: "GGIT_WORK_DIR=/tmp", want: "/tmp"},
		{env: "", want: WorkDir()},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			out, err := testx.Exec("TestWorkDirHelper", test.env)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.want, out)
		})
	}
}

func TestCustomDirHelper(_ *testing.T) {
	if !testx.WantHelperProcess() {
		return
	}

	fmt.Fprintln(os.Stdout, CustomDir())
}

func TestCustomDirFallbackHelper(_ *testing.T) {
	if !testx.WantHelperProcess() {
		return
	}

	// Create a temporary 'data' directory in the work directory
	workDir := WorkDir()
	dataDir := filepath.Join(workDir, "data")
	_ = os.MkdirAll(dataDir, 0755)
	defer os.RemoveAll(dataDir)

	fmt.Fprintln(os.Stdout, CustomDir())
}

func TestCustomDir(t *testing.T) {
	tests := []struct {
		env  string
		want string
	}{
		{env: "GGIT_CUSTOM=/tmp", want: "/tmp"},
		{env: "", want: CustomDir()},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			out, err := testx.Exec("TestCustomDirHelper", test.env)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.want, out)
		})
	}
}

func TestCustomDirFallback(t *testing.T) {
	t.Run("fallback to data", func(t *testing.T) {
		// We use a temporary work directory to avoid messing with the real one
		tempDir, err := os.MkdirTemp("", "ggit-work-dir")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		env := "GGIT_WORK_DIR=" + tempDir
		
		out, err := testx.Exec("TestCustomDirFallbackHelper", env)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, filepath.Join(tempDir, "data"), out)
	})
}

func TestHomeDirHelper(_ *testing.T) {
	if !testx.WantHelperProcess() {
		return
	}

	fmt.Fprintln(os.Stdout, HomeDir())
}

func TestHomeDir(t *testing.T) {
	tests := []struct {
		envs []string
		want string
	}{
		{envs: []string{"HOME=/tmp"}, want: "/tmp"},
		{envs: []string{`USERPROFILE=C:\Users\Joe`}, want: `C:\Users\Joe`},
		{envs: []string{`HOMEDRIVE=C:`, `HOMEPATH=\Users\Joe`}, want: `C:\Users\Joe`},
		{envs: nil, want: ""},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			out, err := testx.Exec("TestHomeDirHelper", test.envs...)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.want, out)
		})
	}
}
