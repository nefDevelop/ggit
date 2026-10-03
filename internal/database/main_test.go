package database

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/testx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	// Los hooks de Git ejecutan conf.AppPath(), que bajo "go test" es este mismo
	// binario de pruebas. Sin esta guarda el hook volveria a entrar en m.Run() y
	// reejecutaria la suite completa, que hace el push que lo dispara otra vez.
	if len(os.Args) > 1 && os.Args[1] == "hook" {
		os.Exit(0)
	}

	flag.Parse()

	level := logger.Silent
	if !testing.Verbose() {
		// Remove the primary logger and register a noop logger.
		log.Remove(log.DefaultConsoleName)
		err := log.New("noop", testx.InitNoopLogger)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		level = logger.Info
	}

	// NOTE: AutoMigrate does not respect logger passed in gorm.Config.
	logger.Default = logger.Default.LogMode(level)

	conf.UseSQLite3 = true

	// Repository.Root is empty until conf.Init runs, which test binaries never
	// do. Without this the repository fixtures are created relative to this
	// package directory and rewrite files that are tracked in git.
	// El shebang del hook se compone con ScriptType. Vacio produce
	// "#!/usr/bin/env " sin interprete, y el hook se acaba invocando a si mismo
	// en bucle: el push no termina nunca y el proceso se queda girando.
	conf.Repository.ScriptType = "bash"

	if _, err := conf.UseTempRepositoryRoot(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func newTestDB(t *testing.T, suite string) *gorm.DB {
	db := dbtest.NewDB(t, suite, Tables...)
	SetHandle(db)
	return db
}

func clearTables(t *testing.T, db *gorm.DB) error {
	if t.Failed() {
		return nil
	}

	for _, t := range Tables {
		err := db.Where("TRUE").Delete(t).Error
		if err != nil {
			return err
		}
	}
	return nil
}
