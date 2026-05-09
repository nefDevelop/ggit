package repo

import (
	"path/filepath"
	"strings"
)

var extensionToLanguage = map[string]string{
	".go":            "Go",
	".js":            "JavaScript",
	".jsx":           "JavaScript",
	".mjs":           "JavaScript",
	".cjs":           "JavaScript",
	".ts":            "TypeScript",
	".tsx":           "TypeScript",
	".svelte":        "Svelte",
	".vue":           "Vue",
	".py":            "Python",
	".pyw":           "Python",
	".java":          "Java",
	".class":         "Java",
	".c":             "C",
	".cpp":           "C++",
	".cc":            "C++",
	".cxx":           "C++",
	".h":             "C",
	".hpp":           "C++",
	".hh":            "C++",
	".cs":            "C#",
	".rb":            "Ruby",
	".php":           "PHP",
	".html":          "HTML",
	".htm":           "HTML",
	".css":           "CSS",
	".scss":          "SCSS",
	".sass":          "Sass",
	".less":          "Less",
	".sh":            "Shell",
	".bash":          "Shell",
	".zsh":           "Shell",
	".fish":          "Shell",
	".ps1":           "PowerShell",
	".sql":           "SQL",
	".yaml":          "YAML",
	".yml":           "YAML",
	".json":          "JSON",
	".jsonc":         "JSON",
	".xml":           "XML",
	".md":            "Markdown",
	".markdown":      "Markdown",
	".rs":            "Rust",
	".kt":            "Kotlin",
	".kts":           "Kotlin",
	".swift":         "Swift",
	".pl":            "Perl",
	".pm":            "Perl",
	".lua":           "Lua",
	".r":             "R",
	".dart":          "Dart",
	".ex":            "Elixir",
	".exs":           "Elixir",
	".erl":           "Erlang",
	".hrl":           "Erlang",
	".hs":            "Haskell",
	".lhs":           "Haskell",
	".clj":           "Clojure",
	".cljs":          "Clojure",
	".cljc":          "Clojure",
	".edn":           "Clojure",
	".scala":         "Scala",
	".sc":            "Scala",
	".v":             "Verilog",
	".sv":            "SystemVerilog",
	".dockerfile":    "Dockerfile",
	"Dockerfile":     "Dockerfile",
	".f":             "Fortran",
	".f90":           "Fortran",
	".f95":           "Fortran",
	".ada":           "Ada",
	".ads":           "Ada",
	".adb":           "Ada",
	".asm":           "Assembly",
	".s":             "Assembly",
	".coffee":        "CoffeeScript",
	".elm":           "Elm",
	".fs":            "F#",
	".fsi":           "F#",
	".fsx":           "F#",
	".fsscript":      "F#",
	".groovy":        "Groovy",
	".lisp":          "Lisp",
	".lsp":           "Lisp",
	".cl":            "Lisp",
	".matlab":        "MATLAB",
	".m":             "Objective-C",
	".mm":            "Objective-C++",
	".nim":           "Nim",
	".pas":           "Pascal",
	".p":             "Pascal",
	".pp":            "Pascal",
	".raku":          "Raku",
	".sol":           "Solidity",
	".tcl":           "Tcl",
	".zig":           "Zig",
	".cmake":         "CMake",
	"CMakeLists.txt": "CMake",
	".makefile":      "Makefile",
	"Makefile":       "Makefile",
	"Gulpfile.js":    "JavaScript",
	"Gruntfile.js":   "JavaScript",
	"package.json":   "JSON",
	".gitignore":     "Git Config",
	".dockerignore":  "Git Config",
	".editorconfig":  "Config",
	".env":           "Config",
	".toml":          "TOML",
	".ini":           "INI",
	".vbs":           "VBScript",
	".bat":           "Batch",
	".cmd":           "Batch",
	".graphql":       "GraphQL",
	".gql":           "GraphQL",
	".proto":         "Protocol Buffers",
	".thrift":        "Thrift",
	".f#":            "F#",
	".jl":            "Julia",
	".ipynb":         "Jupyter Notebook",
	".awk":           "AWK",
	".sed":           "Sed",
	".d":             "D",
	".tex":           "TeX",
	".bib":           "TeX",
	".ahk":           "AutoHotkey",
	".au3":           "AutoIt",
	".feature":       "Gherkin",
	".gradle":        "Gradle",
	".properties":    "Java Properties",
	".scheme":        "Scheme",
	".scm":           "Scheme",
	".ss":            "Scheme",
	".rkt":           "Racket",
	".vhd":           "VHDL",
	".vhdl":          "VHDL",
	".glsl":          "GLSL",
	".vert":          "GLSL",
	".frag":          "GLSL",
	".geom":          "GLSL",
	".comp":          "GLSL",
	".coq":           "Coq",
	".lean":          "Lean",
	".idr":           "Idris",
	".lidr":          "Idris",
	".agda":          "Agda",
	".mq4":           "MQL4",
	".mq5":           "MQL5",
	".mqh":           "MQL",
}

func getLanguageByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// Check for specific filenames like "Dockerfile"
		if lang, ok := extensionToLanguage[filename]; ok {
			return lang
		}
		return ""
	}
	return extensionToLanguage[ext]
}

func isVendor(path string) bool {
	parts := strings.Split(strings.ToLower(path), "/")
	for _, part := range parts {
		if part == "node_modules" || part == "vendor" || part == "bower_components" || part == "dist" || part == "build" {
			return true
		}
	}
	return false
}

func isDotFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func isDocumentation(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasPrefix(lower, "readme") || strings.HasPrefix(lower, "license") || strings.HasPrefix(lower, "contributing") || strings.HasPrefix(lower, "changelog")
}

func isConfiguration(filename string) bool {
	lower := strings.ToLower(filename)
	ext := filepath.Ext(lower)
	return ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml" || ext == ".ini" ||
		strings.Contains(lower, "config") || strings.HasPrefix(lower, "package-lock") || strings.HasPrefix(lower, "yarn.lock")
}
