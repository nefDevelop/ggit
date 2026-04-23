package main

import (
	"fmt"
	"os"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
)

func main() {
	// Setup work dir for testing
	workDir, _ := os.Getwd()
	os.Setenv("GGIT_WORK_DIR", workDir)

	conf.Init("")
	database.LoadRepoConfig()
	fmt.Printf("Gitignores: %v\n", len(database.Gitignores))
	fmt.Printf("GlobalIgnores: %v\n", len(database.GlobalIgnores))
	for _, g := range database.GlobalIgnores {
		if len(g) > 0 && g[0] == '_' || (len(g) > 7 && g[7] == '_') {
			fmt.Printf("Found global: %s\n", g)
		}
	}
}
