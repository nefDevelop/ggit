package strx

import (
	"crypto/rand"
	"html"
	"math/big"
	"regexp"
	"strings"
	"unicode"
)

// ToUpperFirst returns s with only the first Unicode letter mapped to its upper case.
func ToUpperFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + s[i+1:]
	}
	return ""
}

// RandomChars returns a generated string in given number of random characters.
func RandomChars(n int) (string, error) {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	randomInt := func(max *big.Int) (int, error) {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}

		return int(r.Int64()), nil
	}

	buffer := make([]byte, n)
	max := big.NewInt(int64(len(alphanum)))
	for i := range n {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = alphanum[index]
	}

	return string(buffer), nil
}

// Ellipsis returns a truncated string and appends "..." to the end of the
// string if the string length is larger than the threshold. Otherwise, the
// original string is returned.
func Ellipsis(str string, threshold int) string {
	if len(str) <= threshold || threshold < 0 {
		return str
	}
	return str[:threshold] + "..."
}

// Truncate returns a truncated string if its length is over the limit.
// Otherwise, it returns the original string.
func Truncate(str string, limit int) string {
	if len(str) < limit {
		return str
	}
	return str[:limit]
}

// ContainsFold reports whether s is within the slice, ignoring case.
func ContainsFold(ss []string, s string) bool {
	for _, v := range ss {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

// ToSnakeCase converts a string into snake_case.
func ToSnakeCase(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 5) // Pre-allocate some extra space

	var lastWasUnderscore bool
	for i, c := range s {
		isUpper := unicode.IsUpper(c)
		isLower := unicode.IsLower(c)
		isDigit := unicode.IsDigit(c)

		if isUpper {
			if i > 0 && !lastWasUnderscore {
				sb.WriteByte('_')
			}
			sb.WriteRune(unicode.ToLower(c))
			lastWasUnderscore = false
		} else if isLower || isDigit {
			sb.WriteRune(c)
			lastWasUnderscore = false
		} else {
			if i > 0 && !lastWasUnderscore {
				sb.WriteByte('_')
				lastWasUnderscore = true
			}
		}
	}

	return strings.TrimSuffix(sb.String(), "_")
}

// Expand replaces {var} format strings with corresponding values from the match map.
func Expand(s string, match map[string]string) string {
	var args []string
	for k, v := range match {
		args = append(args, "{"+k+"}", v)
	}
	return strings.NewReplacer(args...).Replace(s)
}

var (
	reLinks = regexp.MustCompile(`(?i)<a\s+[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	reTags  = regexp.MustCompile(`(?i)<[^>]+>`)
)

// HTMLToText converts simple HTML into plain text, optimized for Slim HTML emails.
func HTMLToText(htmlStr string) string {
	// Replace line breaks and paragraphs with actual newlines
	htmlStr = strings.ReplaceAll(htmlStr, "<br>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "<br/>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "</p>", "\n\n")

	// Format links nicely: <a href="url">text</a> -> text (url)
	htmlStr = reLinks.ReplaceAllString(htmlStr, "$2 ($1)")

	// Strip all remaining HTML tags
	text := reTags.ReplaceAllString(htmlStr, "")

	// Unescape HTML entities (e.g., &amp; -> &) and trim extra spaces
	return strings.TrimSpace(html.UnescapeString(text))
}
