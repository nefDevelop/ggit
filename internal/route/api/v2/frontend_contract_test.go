package v2

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// This file guards the frontend/backend contract at the routing level.
//
// It exists because the mismatches it catches were all real and all silent: the
// repository delete sent a DELETE with no body while the handler demanded a JSON
// one, and the issue editor sent PATCH to a route registered as PUT, so the whole
// sidebar answered 405. A response of 200 or 405 told the user nothing useful and
// nothing broke a test.
//
// The route table is read from RegisterRoutes itself rather than parsed from the
// source, so group prefixes and nested paths are exact.

type routeKey struct {
	method string
	path   string
}

// dynamicPrefixes are frontend URLs whose last segment is chosen at runtime from
// a closed set that the routes do cover. They are listed with the reason instead
// of being matched, so a new one has to be justified.
var dynamicPrefixes = map[string]string{
	"/api/v2/:p/:p/:p":    "kind is typed 'watchers' | 'stars' | 'forks'",
	"/api/v2/users/:p/:p": "tab comes from the tab list of the profile page",
}

func registeredRoutes(t *testing.T) map[routeKey]bool {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterRoutes(router.Group("/api/v2"))

	routes := make(map[routeKey]bool)
	for _, route := range router.Routes() {
		routes[routeKey{route.Method, route.Path}] = true
	}
	return routes
}

// apiCall is one apiFetch invocation found in the frontend.
type apiCall struct {
	method string
	url    string
	file   string
	line   int
}

var (
	apiFetchPattern = regexp.MustCompile("apiFetch\\(\\s*[`\"']([^`\"']+)[`\"']")
	methodPattern   = regexp.MustCompile(`method:\s*["']([A-Za-z]+)["']`)
	placeholder     = regexp.MustCompile(`\$\{[^}]*\}`)
)

// nonAPICalls are endpoints that are not part of the versioned API.
var nonAPICalls = map[string]string{
	"/health": "health check, registered in cmd/ggit/web.go",
}

func frontendCalls(t *testing.T) []apiCall {
	t.Helper()

	// The test runs from the package directory, so the repository root is found by
	// walking up rather than assumed. A wrong guess would skip the test silently,
	// which is exactly what it exists to prevent.
	// filepath.Dir(".") is ".", not the parent, so the walk starts from an
	// absolute path.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("directorio de trabajo: %v", err)
	}
	root := ""
	for {
		if _, err := os.Stat(filepath.Join(dir, "frontend", "src")); err == nil {
			root = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	if root == "" {
		t.Fatal("no se encuentra frontend/src subiendo desde el directorio del paquete")
	}

	var calls []apiCall
	err = filepath.Walk(filepath.Join(root, "frontend", "src"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || (filepath.Ext(path) != ".svelte" && filepath.Ext(path) != ".ts") {
			return nil
		}
		if strings.Contains(path, "__tests__") || strings.HasSuffix(path, ".test.ts") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		src := string(content)
		rel, _ := filepath.Rel(root, path)

		// SubmatchIndex, not StringIndex: the URL is a capture group and the plain
		// variant does not report its position.
		for _, loc := range apiFetchPattern.FindAllStringSubmatchIndex(src, -1) {
			if len(loc) < 4 || loc[2] < 0 {
				continue
			}
			url := src[loc[2]:loc[3]]

			// The verb has to be read from this call's own arguments. Scanning
			// forward for the next "})" picks up the method of a later call when
			// this one takes no options object.
			method := http.MethodGet
			if m := methodPattern.FindStringSubmatch(callArguments(src, loc[1])); m != nil {
				method = strings.ToUpper(m[1])
			}
			calls = append(calls, apiCall{
				method: method,
				url:    url,
				file:   rel,
				line:   strings.Count(src[:loc[0]], "\n") + 1,
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("recorrer el frontend: %v", err)
	}
	return calls
}

// normalise turns a frontend URL into the shape of a registered route: template
// placeholders become parameters and the query string is dropped.
func normalise(url string) string {
	url = placeholder.ReplaceAllString(url, ":p")
	if i := strings.Index(url, "?"); i >= 0 {
		// Optional chaining inside a placeholder was already replaced above, so a
		// remaining "?" always starts the query string.
		url = url[:i]
	}
	// A placeholder glued to the end of a non-empty segment is a query string held
	// in a variable, as in "/releases${query}", and contributes no segment. One
	// that is a segment of its own, as in "/auths/${id}", is a real parameter and
	// has to be kept.
	if strings.HasSuffix(url, ":p") && !strings.HasSuffix(url, "/:p") {
		url = strings.TrimSuffix(url, ":p")
	}
	return url
}

// matchesRoute reports whether a concrete path satisfies a registered pattern.
// It is deliberately simpler than a router: it compares segment counts and
// accepts a parameter where the pattern declares one.
func matchesRoute(pattern, path string) bool {
	pp := strings.Split(strings.Trim(pattern, "/"), "/")
	cp := strings.Split(strings.Trim(path, "/"), "/")

	// A trailing wildcard absorbs any number of segments, so the counts only have
	// to agree when the pattern does not end in one.
	if len(pp) != len(cp) && !strings.HasPrefix(pp[len(pp)-1], "*") {
		return false
	}
	for i := range pp {
		if i >= len(cp) {
			return strings.HasPrefix(pp[i], "*")
		}
		if strings.HasPrefix(pp[i], "*") {
			// A trailing wildcard swallows the rest, including an empty tail.
			return true
		}
		if i >= len(cp) {
			return false
		}
		if strings.HasPrefix(pp[i], ":") {
			if cp[i] == "" {
				return false
			}
			continue
		}
		if pp[i] != cp[i] {
			return false
		}
	}
	return len(pp) == len(cp) || strings.HasPrefix(pp[len(pp)-1], "*")
}

func TestFrontendCallsMatchRegisteredRoutes(t *testing.T) {
	routes := registeredRoutes(t)
	calls := frontendCalls(t)

	if len(calls) == 0 {
		t.Fatal("no se encontro ninguna llamada a apiFetch: el extractor estara roto")
	}

	byMethod := make(map[string][]string)
	for key := range routes {
		byMethod[key.method] = append(byMethod[key.method], key.path)
	}

	var unknownVerb, unknownRoute []string

	for _, call := range calls {
		path := normalise(call.url)
		where := fmt.Sprintf("%s:%d", call.file, call.line)

		if _, ok := nonAPICalls[path]; ok {
			continue
		}
		if _, ok := dynamicPrefixes[path]; ok {
			continue
		}

		if matchesAny(byMethod[call.method], path) {
			continue
		}

		// The path exists but under a different verb. That is the bug that made
		// the issue sidebar answer 405 for every edit.
		found := false
		for method, patterns := range byMethod {
			if method == call.method {
				continue
			}
			if matchesAny(patterns, path) {
				unknownVerb = append(unknownVerb, fmt.Sprintf(
					"%s llama a %s %s, pero esa ruta esta registrada como %s",
					where, call.method, path, method))
				found = true
				break
			}
		}
		if found {
			continue
		}

		unknownRoute = append(unknownRoute, fmt.Sprintf(
			"%s llama a %s %s, que no esta registrada", where, call.method, path))
	}

	sort.Strings(unknownVerb)
	sort.Strings(unknownRoute)

	if len(unknownVerb) > 0 {
		t.Errorf("verbo incorrecto en %d llamadas:\n  %s", len(unknownVerb), strings.Join(unknownVerb, "\n  "))
	}
	if len(unknownRoute) > 0 {
		t.Errorf("ruta inexistente en %d llamadas:\n  %s", len(unknownRoute), strings.Join(unknownRoute, "\n  "))
	}
	t.Logf("comprobadas %d llamadas a apiFetch contra %d rutas", len(calls), len(routes))
}

// callArguments returns the argument list of the apiFetch call whose opening
// parenthesis is at start, stopping at its matching close.
func callArguments(src string, start int) string {
	// start points just past the URL, so the opening parenthesis of the call has
	// already been consumed. A call without options begins with the closing one,
	// hence depth 1 rather than 0.
	depth := 1
	for i := start; i < len(src); i++ {
		switch src[i] {
		case '(', '{', '[':
			depth++
		case ')', '}', ']':
			depth--
			if depth == 0 {
				return src[start:i]
			}
		}
	}
	return src[start:]
}

func matchesAny(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if matchesRoute(pattern, path) {
			return true
		}
	}
	return false
}
