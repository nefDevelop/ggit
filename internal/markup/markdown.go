package markup

import (
	"path/filepath"
	"strings"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/lazyregexp"
)

// IsMarkdownFile reports whether name looks like a Markdown file based on its extension.
func IsMarkdownFile(name string) bool {
	extension := strings.ToLower(filepath.Ext(name))
	for _, ext := range conf.Markdown.FileExtensions {
		if strings.ToLower(ext) == extension {
			return true
		}
	}
	return false
}

// validLinksPattern matches strings that are valid URLs.
var validLinksPattern = lazyregexp.New(`^[a-z][\w-]+://|^mailto:`)

// isLink reports whether link fits valid format.
func isLink(link []byte) bool {
	return validLinksPattern.Match(link)
}


// RawMarkdown renders content in Markdown syntax to HTML using Goldmark.
func RawMarkdown(body []byte, urlPrefix string) []byte {
	return goldmarkRender(body, urlPrefix)
}

// Markdown takes a string or []byte and renders to HTML in Markdown syntax with special links.
func Markdown(input any, urlPrefix string, metas map[string]string) []byte {
	return Render(TypeMarkdown, input, urlPrefix, metas)
}
