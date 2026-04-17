package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/markup"
)

// MarkdownOptions represents the request body for rendering markdown.
type MarkdownOptions struct {
	Text    string `json:"text" binding:"required"`
	Context string `json:"context"` // Repo link for resolving relative links
}

// Markdown renders markdown text to HTML.
func RenderMarkdown(c *context.APIContext) {
	var f MarkdownOptions
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	if f.Text == "" {
		c.String(http.StatusOK, "")
		return
	}

	rendered := markup.Markdown([]byte(f.Text), f.Context, nil)
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(rendered)
}

// MarkdownRaw renders raw markdown text to HTML.
func RenderMarkdownRaw(c *context.APIContext) {
	body, err := c.RawData()
	if err != nil {
		c.Error(err, "read body")
		return
	}
	rendered := markup.SanitizeBytes(markup.RawMarkdown(body, ""))
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(rendered)
}
