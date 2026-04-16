package ginx

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"gogs.io/gogs/internal/conf"
	internalTmpl "gogs.io/gogs/internal/template"
)

// SetupTemplates configures Gin to use the same Go templates as Macaron did.
func SetupTemplates(r *gin.Engine) {
	// 1. Gather all template functions
	funcMap := internalTmpl.FuncMap()[0] // Gogs returns a slice, we take the first map

	// 2. Load templates from directory
	// In Gogs, templates are in WorkDir + "/templates" and can be overridden by CustomDir + "/templates"
	workTmplDir := filepath.Join(conf.WorkDir(), "templates")

	// Gin's LoadHTMLGlob doesn't easily support merging two directories with priority.
	// We'll manually parse the templates to ensure custom overrides work.
	tmpl := template.New("").Funcs(funcMap)
	
	// TODO: Traverse directories and ParseFiles
	// For now, let's use a simple LoadHTMLGlob for the proof of concept
	r.SetHTMLTemplate(tmpl)
	r.LoadHTMLGlob(filepath.Join(workTmplDir, "**", "*.tmpl"))
}
