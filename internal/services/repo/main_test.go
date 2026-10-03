package repo

import (
	"fmt"
	"os"
	"testing"

	"github.com/nef734/ggit/internal/conf"
)

// TestMain redirects repository fixtures outside the source tree. Without it
// conf.Repository.Root stays empty, repox joins the user name onto it directly,
// and every repository a test creates lands inside this package directory,
// rewriting files that git tracks.
func TestMain(m *testing.M) {
	// Los hooks de Git ejecutan conf.AppPath(), que bajo "go test" es este mismo
	// binario de pruebas. Sin esta guarda el hook volveria a entrar en m.Run() y
	// reejecutaria la suite completa, que hace el push que lo dispara otra vez.
	if len(os.Args) > 1 && os.Args[1] == "hook" {
		os.Exit(0)
	}
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
