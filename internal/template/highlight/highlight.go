package highlight

import (
	"bytes"
	"path"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/go-enry/go-enry/v2"
	"github.com/nef734/ggit/internal/conf"
)

var highlightMapping = map[string]string{
	".txt": "nohighlight",
}

func NewContext() {
	keys := conf.File.Section("highlight.mapping").Keys()
	for i := range keys {
		highlightMapping[keys[i].Name()] = keys[i].Value()
	}
}

// FileNameToHighlightClass returns the best match for highlight class name.
func FileNameToHighlightClass(fname string) string {
	fname = strings.ToLower(fname)
	ext := path.Ext(fname)
	if name, ok := highlightMapping[ext]; ok {
		return name
	}
	lang, _ := enry.GetLanguageByExtension(fname)
	if lang != "" {
		return strings.ToLower(lang)
	}
	return ""
}

// Code returns the highlighted HTML lines using Chroma.
func Code(fileName, code string) []string {
	lang := enry.GetLanguage(fileName, []byte(code))
	if lang == "" {
		lang = "fallback"
	}

	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(false), html.PreventSurroundingPre(true))
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return splitLines(code)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return splitLines(code)
	}

	return splitLines(buf.String())
}

func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
