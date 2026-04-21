package ginx

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"gogs.io/gogs/internal/conf"
	internalTmpl "gogs.io/gogs/internal/template"
)

// SetupTemplates configures Gin to use the same Go templates as Macaron did.
func SetupTemplates(r *gin.Engine) {
	fmt.Printf("[DEBUG] SetupTemplates: Initializing on %p\n", r)
	// 1. Gather all template functions
	funcMap := internalTmpl.FuncMap()[0] // Ggit returns a slice, we take the first map

	// 2. Load templates from directory
	// In Gogs, templates are in WorkDir + "/templates" and can be overridden by CustomDir + "/templates"
	workTmplDir := filepath.Join(conf.WorkDir(), "templates")
	tmpl := template.New("").Funcs(funcMap)

	err := filepath.Walk(workTmplDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		rel, err := filepath.Rel(workTmplDir, path)
		if err != nil {
			return err
		}

		// Template name should be "repo/home" for "repo/home.tmpl"
		// Ggit uses "/" as separator even on Windows for template names
		name := filepath.ToSlash(strings.TrimSuffix(rel, ".tmpl"))

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = tmpl.New(name).Parse(string(b))
		return err
	})

	if err != nil {
		fmt.Printf("[DEBUG] Failed to walk template dir: %v\n", err)
		return
	}

	templates := tmpl.Templates()
	fmt.Printf("[DEBUG] Loaded %d HTML templates from %s\n", len(templates), workTmplDir)
	r.SetHTMLTemplate(tmpl)
}
