package markup

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/nef734/ggit/internal/lazyregexp"
	"gopkg.in/yaml.v3"
)

// Highlight highlights the given code using Chroma.
func Highlight(code, language, styleName string) (string, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(styleName)
	if style == nil {
		style = styles.GitHub
	}

	formatter := formatters.Get("html")
	if formatter == nil {
		formatter = formatters.TTY
	}

	// Use inline styles and no pre wrapper to allow frontend to handle line numbers
	htmlFormatter := html.New(html.WithClasses(false), html.PreventSurroundingPre(true))

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = htmlFormatter.Format(&buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// IsReadmeFile reports whether name looks like a README file based on its extension.
func IsReadmeFile(name string) bool {
	name = strings.ToLower(name)
	if name == "readme" {
		return true
	}
	return strings.HasPrefix(name, "readme.")
}

// IsMarkdownFile reports whether name looks like a Markdown file based on its extension.
func IsMarkdownFile(name string) bool {
	extension := strings.ToLower(strings.TrimPrefix(name, "."))
	for _, ext := range []string{"md", "markdown", "mdown", "mkd"} {
		if extension == ext {
			return true
		}
	}
	return false
}

// IsOrgModeFile reports whether name looks like an Org-mode file based on its extension.
func IsOrgModeFile(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".org")
}

// IsIPythonNotebook reports whether name looks like a IPython notebook based on its extension.
func IsIPythonNotebook(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".ipynb")
}

const (
	IssueNameStyleNumeric      = "numeric"
	IssueNameStyleAlphanumeric = "alphanumeric"
)

type Type string

const (
	TypeUnrecognized    Type = "unrecognized"
	TypeMarkdown        Type = "markdown"
	TypeOrgMode         Type = "orgmode"
	TypeIPythonNotebook Type = "ipynb"
)

// MentionPattern matches string that mentions someone, e.g. @Unknwon
var MentionPattern = lazyregexp.New(`(\s|^|\W)@[0-9a-zA-Z-_\.]+`)

var frontmatterPattern = lazyregexp.New(`(?s)^---\r?\n(.*?)\r?\n---`)

// ExtractMetadata extracts YAML metadata from the markdown content.
func ExtractMetadata(body []byte) (map[string]any, error) {
	match := frontmatterPattern.FindSubmatch(body)
	if len(match) < 2 {
		return nil, nil
	}

	var meta map[string]any
	err := yaml.Unmarshal(match[1], &meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// FindAllMentions matches mention patterns in given content
// and returns a list of found user names without @ prefix.
func FindAllMentions(content string) []string {
	mentions := MentionPattern.FindAllString(content, -1)
	for i := range mentions {
		mentions[i] = mentions[i][strings.Index(mentions[i], "@")+1:] // Strip @ character
	}
	return mentions
}

// Sanitize is a placeholder since sanitization is now done on the frontend.
// It returns the input verbatim for backend compatibility.
func Sanitize(input string) string {
	return input
}

// Detect returns best guess of a markup type based on file name.
func Detect(filename string) Type {
	switch {
	case IsMarkdownFile(filename):
		return TypeMarkdown
	case IsOrgModeFile(filename):
		return TypeOrgMode
	case IsIPythonNotebook(filename):
		return TypeIPythonNotebook
	default:
		return TypeUnrecognized
	}
}
