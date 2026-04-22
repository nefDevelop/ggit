package markup_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gogs.io/gogs/internal/conf"
	. "gogs.io/gogs/internal/markup"
)

func Test_IsMarkdownFile(t *testing.T) {
	// TODO: Refactor to accept a list of extensions
	conf.Markdown.FileExtensions = strings.Split(".md,.markdown,.mdown,.mkd", ",")
	tests := []struct {
		ext    string
		expVal bool
	}{
		{ext: ".md", expVal: true},
		{ext: ".markdown", expVal: true},
		{ext: ".mdown", expVal: true},
		{ext: ".mkd", expVal: true},
		{ext: ".org", expVal: false},
		{ext: ".rst", expVal: false},
		{ext: ".asciidoc", expVal: false},
	}
	for _, test := range tests {
		assert.Equal(t, test.expVal, IsMarkdownFile(test.ext))
	}
}

func Test_Markdown(t *testing.T) {
	conf.Server.ExternalURL = "http://localhost:3000/"

	tests := []struct {
		input  string
		expVal string
	}{
		// Issue URL
		{input: "http://localhost:3000/user/repo/issues/3333", expVal: "<p><a href=\"http://localhost:3000/user/repo/issues/3333\">#3333</a></p>\n"},
		// Commit URL
		{input: "http://localhost:3000/user/project/commit/d8a994ef243349f321568f9e36d5c3f444b99cae", expVal: "<p><a href=\"http://localhost:3000/user/project/commit/d8a994ef243349f321568f9e36d5c3f444b99cae\"><code>d8a994ef24</code></a></p>\n"},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := Markdown([]byte(test.input), "http://localhost:3000/user/repo", nil)
			assert.Equal(t, test.expVal, string(res))
		})
	}
}
