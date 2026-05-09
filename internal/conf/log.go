package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/nef734/ggit/internal/log"
	"gopkg.in/ini.v1"
)

type loggerConf struct {
	Level log.Level
	Path  string
}

type logConf struct {
	RootPath string
	Modes    []string
	Configs  []*loggerConf
}

// Log settings
var Log *logConf

// initLogConf returns parsed logging configuration from given INI file. When the
// argument "hookMode" is true, it only initializes the root path for log files.
func initLogConf(cfg *ini.File, hookMode bool) (_ *logConf, hasConsole bool, _ error) {
	rootPath := cfg.Section("log").Key("ROOT_PATH").MustString(filepath.Join(WorkDir(), "log"))
	if hookMode {
		return &logConf{
			RootPath: ensureAbs(rootPath),
		}, false, nil
	}

	modes := strings.Split(cfg.Section("log").Key("MODE").MustString("console"), ",")
	lc := &logConf{
		RootPath: ensureAbs(rootPath),
		Modes:    make([]string, 0, len(modes)),
		Configs:  make([]*loggerConf, 0, len(modes)),
	}

	// Iterate over [log.*] sections to initialize individual logger.
	levelMappings := map[string]log.Level{
		"trace": log.LevelTrace,
		"info":  log.LevelInfo,
		"warn":  log.LevelWarn,
		"error": log.LevelError,
		"fatal": log.LevelFatal,
	}

	for i := range modes {
		modes[i] = strings.ToLower(strings.TrimSpace(modes[i]))
		secName := "log." + modes[i]
		sec, err := cfg.GetSection(secName)
		if err != nil {
			return nil, hasConsole, fmt.Errorf("missing configuration section [%s] for %q logger", secName, modes[i])
		}

		level := levelMappings[strings.ToLower(sec.Key("LEVEL").MustString("trace"))]
		var c *loggerConf
		switch modes[i] {
		case log.DefaultConsoleName:
			hasConsole = true
			c = &loggerConf{
				Level: level,
			}

		case log.DefaultFileName:
			logPath := filepath.Join(lc.RootPath, "gogs.log")
			c = &loggerConf{
				Level: level,
				Path:  logPath,
			}

		default:
			continue
		}

		lc.Modes = append(lc.Modes, modes[i])
		lc.Configs = append(lc.Configs, c)
	}

	return lc, hasConsole, nil
}

// InitLogging initializes the logging service of the application. When the
// "hookMode" is true, it only initializes the root path for log files without
// creating any logger. It will also not remove the primary logger in "hookMode"
// and is up to the caller to decide when to remove it.
func InitLogging(hookMode bool) {
	logConf, hasConsole, err := initLogConf(File, hookMode)
	if err != nil {
		log.Fatal("Failed to init logging configuration: %v", err)
	}
	defer func() {
		Log = logConf
	}()

	if hookMode {
		return
	}

	err = os.MkdirAll(logConf.RootPath, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create log directory: %v", err)
	}

	for i, mode := range logConf.Modes {
		c := logConf.Configs[i]

		var err error
		switch mode {
		case log.DefaultConsoleName:
			err = log.NewConsole(mode, c.Level)
		case log.DefaultFileName:
			err = log.NewFile(mode, c.Path, c.Level)
		default:
			continue
		}

		if err != nil {
			log.Fatal("Failed to init %s logger: %v", mode, err)
			return
		}
		log.Trace("Log mode: %s (%s)", strings.Title(mode), strings.Title(strings.ToLower(c.Level.String())))
	}

	// ⚠️ WARNING: It is only safe to remove the primary logger until
	// there are other loggers that are initialized. Otherwise, the
	// application will print errors to nowhere.
	if !hasConsole {
		log.Remove(log.DefaultConsoleName)
	}
}
