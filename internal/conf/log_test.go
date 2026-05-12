package conf

import (
	"path/filepath"
	"testing"

	log "github.com/nef734/ggit/internal/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func Test_initLogConf(t *testing.T) {
	t.Run("missing configuration section", func(t *testing.T) {
		f, err := ini.Load([]byte(`
[log]
MODE = console
`))
		if err != nil {
			t.Fatal(err)
		}

		got, hasConsole, err := initLogConf(f, false)
		assert.NotNil(t, err)
		assert.Equal(t, `missing configuration section [log.console] for "console" logger`, err.Error())
		assert.False(t, hasConsole)
		assert.Nil(t, got)
	})

	t.Run("no console logger", func(t *testing.T) {
		f, err := ini.Load([]byte(`
[log]
MODE = file

[log.file]
`))
		if err != nil {
			t.Fatal(err)
		}

		got, hasConsole, err := initLogConf(f, false)
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, hasConsole)
		assert.NotNil(t, got)
	})

	f, err := ini.Load([]byte(`
[log]
ROOT_PATH = log
MODE = console, file
BUFFER_LEN = 50
LEVEL = trace

[log.console]
BUFFER_LEN = 10

[log.file]
LEVEL = INFO
LOG_ROTATE = true
DAILY_ROTATE = true
MAX_SIZE_SHIFT = 20
MAX_LINES = 1000
MAX_DAYS = 3
`))
	if err != nil {
		t.Fatal(err)
	}

	got, hasConsole, err := initLogConf(f, false)
	if err != nil {
		t.Fatal(err)
	}

	logConf := &logConf{
		RootPath: filepath.Join(WorkDir(), "log"),
		Modes: []string{
			log.DefaultConsoleName,
			log.DefaultFileName,
		},
		Configs: []*loggerConf{
			{
				Level: log.LevelTrace,
			}, {
				Level: log.LevelInfo,
				Path:  filepath.Join(WorkDir(), "log", "ggit.log"),
			},
		},
	}
	assert.True(t, hasConsole)
	assert.Equal(t, logConf, got)
}
