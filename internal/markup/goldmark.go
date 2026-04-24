package markup

import (
	"bytes"
	"io"
	"strings"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	chtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"unknwon.dev/clog/v2"
)

// goldmarkRenderer is a Goldmark-based renderer.
var goldmarkRenderer goldmark.Markdown

func init() {
	goldmarkRenderer = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			meta.Meta,
			emoji.Emoji,
			&wikilinkExtension{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&linkTransformer{}, 100),
			),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(),
		),
	)
	goldmarkRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&wikilinkRenderer{}, 0),
		util.Prioritized(&fencedCodeBlockRenderer{}, 0),
	))
}

var urlPrefixKey = parser.NewContextKey()

// goldmarkRender renders markdown using Goldmark.
func goldmarkRender(body []byte, urlPrefix string) []byte {
	var buf bytes.Buffer
	pc := parser.NewContext()
	pc.Set(urlPrefixKey, urlPrefix)

	if err := goldmarkRenderer.Convert(body, &buf, parser.WithContext(pc)); err != nil {
		return body
	}

	return buf.Bytes()
}

type linkTransformer struct{}

func (t *linkTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	urlPrefix, _ := pc.Get(urlPrefixKey).(string)
	if urlPrefix == "" {
		return
	}

	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Link:
			n.Destination = t.maybePrependPrefix(n.Destination, urlPrefix)
		case *ast.Image:
			n.Destination = t.maybePrependPrefix(n.Destination, urlPrefix)
		}
		return ast.WalkContinue, nil
	})
}

func (t *linkTransformer) maybePrependPrefix(destination []byte, urlPrefix string) []byte {
	if len(destination) == 0 || destination[0] == '#' || isLink(destination) {
		return destination
	}
	// Prepend urlPrefix
	return []byte(strings.TrimRight(urlPrefix, "/") + "/" + strings.TrimLeft(string(destination), "/"))
}

// ExtractMetadata extracts YAML/TOML metadata from the markdown content.
func ExtractMetadata(body []byte) (map[string]any, error) {
	pc := parser.NewContext()
	err := goldmarkRenderer.Convert(body, io.Discard, parser.WithContext(pc))
	if err != nil {
		return nil, err
	}
	return meta.Get(pc), nil
}

// --- Wikilinks Extension ---

type wikilinkExtension struct{}

func (e *wikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&wikilinkParser{}, 100),
	))
}

type wikilinkParser struct{}

func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte("[[")) {
		return nil
	}

	// Find the end of the wikilink ]]
	idx := bytes.Index(line, []byte("]]"))
	if idx == -1 {
		return nil
	}

	content := line[2:idx]
	block.Advance(idx + 2)

	link := string(content)
	label := string(content)
	if pipeIdx := strings.Index(link, "|"); pipeIdx != -1 {
		link = strings.TrimSpace(link[:pipeIdx])
		label = strings.TrimSpace(label[pipeIdx+1:])
	}

	// Convert link to Ggit wiki URL style
	link = strings.ReplaceAll(link, " ", "-")

	return &WikiLinkNode{
		Destination: link,
		Label:       label,
	}
}

type WikiLinkNode struct {
	ast.BaseInline
	Destination string
	Label       string
}

func (n *WikiLinkNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindWikiLink = ast.NewNodeKind("WikiLink")

func (n *WikiLinkNode) Kind() ast.NodeKind {
	return KindWikiLink
}

type wikilinkRenderer struct{}

func (r *wikilinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikiLink, r.renderWikiLink)
}

func (r *wikilinkRenderer) renderWikiLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*WikiLinkNode)
	_, _ = w.WriteString(`<a href="`)
	_, _ = w.WriteString(n.Destination)
	_, _ = w.WriteString(`">`)
	_, _ = w.WriteString(n.Label)
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

// --- Fenced Code Block Extension (Chroma) ---

type fencedCodeBlockRenderer struct{}

func (r *fencedCodeBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeBlock, r.renderFencedCodeBlock)
}

func (r *fencedCodeBlockRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	var language string
	var lines *text.Segments
	if fb, ok := node.(*ast.FencedCodeBlock); ok {
		language = string(fb.Language(source))
		lines = fb.Lines()
	} else if cb, ok := node.(*ast.CodeBlock); ok {
		lines = cb.Lines()
	}

	clog.Trace("Rendering code block (lang: %s)", language)

	var code bytes.Buffer
	if lines != nil {
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write(line.Value(source))
		}
	}

	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := chtml.New(chtml.WithClasses(true), chtml.WithLineNumbers(false))

	iterator, err := lexer.Tokenise(nil, code.String())
	if err != nil {
		return ast.WalkStop, err
	}

	err = formatter.Format(w, style, iterator)
	if err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkContinue, nil
}
