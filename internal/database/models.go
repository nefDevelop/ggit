package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database/migrations"
	"github.com/nef734/ggit/internal/dbx"
	"github.com/nef734/ggit/internal/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	HasEngine bool
)

// ensureLogDir creates the directory that will hold the gorm log file. The
// console commands never call InitLogging(false), which is the only path that
// creates the log root, so SetEngine cannot assume the directory is there.
func ensureLogDir(logPath string) error {
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		return fmt.Errorf("create gorm log directory: %v", err)
	}
	return nil
}

func SetEngine() (*gorm.DB, error) {
	var logPath string
	if conf.HookMode {
		logPath = filepath.Join(conf.Log.RootPath, "hooks", "gorm.log")
	} else {
		logPath = filepath.Join(conf.Log.RootPath, "gorm.log")
	}

	// Only "ggit web" reaches InitLogging(false), and that is the single caller
	// that creates the log root. Every console command passes hookMode=true,
	// which returns before the MkdirAll, and then lands here with HookMode still
	// false and no directory to write into. On a fresh install that made
	// "ggit admin create-user" abort with
	// create "gorm.log": open data/log/gorm.log: no such file or directory,
	// so the documented bootstrap command could not bootstrap.
	if err := ensureLogDir(logPath); err != nil {
		return nil, err
	}

	var err error
	var gormLogger logger.Writer
	if conf.HookMode {
		w, err := log.NewFileWriter(
			logPath,
			log.FileRotationConfig{
				Rotate:  true,
				Daily:   true,
				MaxSize: 100 * 1024 * 1024,
				MaxDays: 3,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("create gorm log writer"+": %v", err)
		}
		gormLogger = &dbx.Logger{Writer: w}
	} else {
		gormLogger, err = newLogWriter()
		if err != nil {
			return nil, fmt.Errorf("new log writer"+": %v", err)
		}
	}
	return NewConnection(gormLogger)
}

func NewEngine() error {
	db, err := SetEngine()
	if err != nil {
		return err
	}

	if err = migrations.Migrate(db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// NewTestEngine performs a connection test and migrates the database.
func NewTestEngine() error {
	return NewEngine()
}

type Statistic struct {
	Counter struct {
		User, Org, PublicKey,
		Repo, Watch, Star, Action, Access,
		Issue, Comment, Oauth,
		Release, LoginSource, Webhook,
		Milestone, Label, HookTask,
		Team, UpdateTask, Attachment int64
	}
}

func GetStatistic(ctx context.Context) (stats Statistic) {
	stats.Counter.User = Handle.Users().Count(ctx)
	stats.Counter.Org = CountOrganizations()
	Handle.db.WithContext(ctx).Model(new(PublicKey)).Count(&stats.Counter.PublicKey)
	stats.Counter.Repo = CountRepositories(true)
	Handle.db.WithContext(ctx).Model(new(Watch)).Count(&stats.Counter.Watch)
	Handle.db.WithContext(ctx).Model(new(Star)).Count(&stats.Counter.Star)
	Handle.db.WithContext(ctx).Model(new(Action)).Count(&stats.Counter.Action)
	Handle.db.WithContext(ctx).Model(new(Access)).Count(&stats.Counter.Access)
	Handle.db.WithContext(ctx).Model(new(Issue)).Count(&stats.Counter.Issue)
	Handle.db.WithContext(ctx).Model(new(Comment)).Count(&stats.Counter.Comment)
	stats.Counter.Oauth = 0
	Handle.db.WithContext(ctx).Model(new(Release)).Count(&stats.Counter.Release)
	stats.Counter.LoginSource = Handle.LoginSources().Count(ctx)
	Handle.db.WithContext(ctx).Model(new(Webhook)).Count(&stats.Counter.Webhook)
	Handle.db.WithContext(ctx).Model(new(Milestone)).Count(&stats.Counter.Milestone)
	Handle.db.WithContext(ctx).Model(new(Label)).Count(&stats.Counter.Label)
	Handle.db.WithContext(ctx).Model(new(HookTask)).Count(&stats.Counter.HookTask)
	Handle.db.WithContext(ctx).Model(new(Team)).Count(&stats.Counter.Team)
	Handle.db.WithContext(ctx).Model(new(Attachment)).Count(&stats.Counter.Attachment)
	return stats
}

func Ping() error {
	sqlDB, err := Handle.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Version table for database migrations. Should have only one row with ID == 1.
type Version struct {
	ID      int64 `gorm:"primaryKey"`
	Version int64 `gorm:"not null"`
}
