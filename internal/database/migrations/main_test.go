package migrations

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/testx"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
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

	os.Exit(m.Run())
}
