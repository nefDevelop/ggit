package database

import (
	"net/url"
	"strings"

	"gogs.io/gogs/internal/pathx"
)

// ToWikiPageURL formats a string to corresponding wiki URL name.
func ToWikiPageURL(name string) string {
	return url.QueryEscape(name)
}

// ToWikiPageName formats a URL back to corresponding wiki page name. It enforces
// single-level hierarchy by replacing all "/" with spaces.
func ToWikiPageName(urlString string) string {
	name, _ := url.QueryUnescape(urlString)
	name = pathx.Clean(name)
	return strings.ReplaceAll(name, "/", " ")
}

// HasWiki returns true if repository has wiki.
// NOTE: This now checks for the existence of a "/docs" folder in the main branch.
func (r *Repository) HasWiki() bool {
	// TODO: Implement detection of "/docs" folder in the main branch using gitx.
	return false
}
