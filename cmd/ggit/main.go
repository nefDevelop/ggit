// Ggit is a painless self-hosted Git Service.
package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
)

func init() {
	conf.App.Version = "2026.4.0"
}

func main() {
	cmd := &cli.Command{
		Name:    "Ggit",
		Usage:   "A beautiful and ultra-fast self-hosted Git service",
		Version: conf.App.Version,
		Commands: []*cli.Command{
			&webCommand,
			&servCommand,
			&hookCommand,
			&adminCommand,
			&importCommand,
			&backupCommand,
			&restoreCommand,
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
