package utils

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/kaputi/navani/internal/config"
)

var (
	extensionMap = map[string]string{
		".asm":    "assembly",
		".bash":   "bash",
		".c":      "c",
		".clj":    "clojure",
		".cljs":   "clojurescript",
		".coffee": "coffeescript",
		".cpp":    "cpp",
		".cs":     "csharp",
		".css":    "css",
		".dart":   "dart",
		".el":     "emacs-lisp",
		".elm":    "elm",
		".erl":    "erlang",
		".fs":     "fsharp",
		".go":     "go",
		".groovy": "groovy",
		".h":      "c-header",
		".hpp":    "cpp-header",
		".hs":     "haskell",
		".html":   "html",
		".java":   "java",
		".jl":     "julia",
		".js":     "javascript",
		".json":   "json",
		".jsx":    "javascript-react",
		".kt":     "kotlin",
		".lua":    "lua",
		".m":      "objective-c",
		".md":     "markdown",
		".php":    "php",
		".pl":     "perl",
		".pm":     "perl-module",
		".py":     "python",
		".r":      "r",
		".rb":     "ruby",
		".rs":     "rust",
		".scala":  "scala",
		".scm":    "scheme",
		".sh":     "shell",
		".sql":    "sql",
		".swift":  "swift",
		".ts":     "typescript",
		".tsx":    "typescript-react",
		".vb":     "visual-basic",
		".xml":    "xml",
		".yaml":   "yaml",
		".yml":    "yaml",
	}

	fileTypeIcons = map[string]string{
		"assembly":         "\033[1;33m\033[0m", // Yellow
		"bash":             "\033[1;32m\033[0m", // Green
		"c":                "\033[1;34m\033[0m", // Blue
		"clojure":          "\033[1;35m\033[0m", // Magenta
		"clojurescript":    "\033[1;35m\033[0m", // Magenta
		"coffeescript":     "\033[1;33m\033[0m", // Yellow
		"cpp":              "\033[1;34m\033[0m", // Blue
		"csharp":           "\033[1;34m\033[0m", // Blue
		"css":              "\033[1;36m\033[0m", // Cyan
		"dart":             "\033[1;36m\033[0m", // Cyan
		"emacs-lisp":       "\033[1;32m\033[0m", // Green
		"elm":              "\033[1;32m\033[0m", // Green
		"erlang":           "\033[1;31m\033[0m", // Red
		"fsharp":           "\033[1;35m\033[0m", // Magenta
		"go":               "\033[1;36m\033[0m", // Cyan
		"groovy":           "\033[1;35m\033[0m", // Magenta
		"c-header":         "\033[1;34m\033[0m", // Blue
		"cpp-header":       "\033[1;34m\033[0m", // Blue
		"haskell":          "\033[1;35m\033[0m", // Magenta
		"html":             "\033[1;31m\033[0m", // Red
		"java":             "\033[1;31m\033[0m", // Red
		"julia":            "\033[1;32m\033[0m", // Green
		"javascript":       "\033[1;33m\033[0m", // Yellow
		"json":             "\033[1;33m\033[0m", // Yellow
		"javascript-react": "\033[1;36m\033[0m", // Cyan
		"kotlin":           "\033[1;35m\033[0m", // Magenta
		"lua":              "\033[1;34m\033[0m", // Blue
		"objective-c":      "\033[1;34m\033[0m", // Blue
		"markdown":         "\033[1;33m\033[0m", // Yellow
		"php":              "\033[1;35m\033[0m", // Magenta
		"perl":             "\033[1;35m\033[0m", // Magenta
		"perl-module":      "\033[1;35m\033[0m", // Magenta
		"python":           "\033[1;32m\033[0m", // Green
		"r":                "\033[1;34mﳒ\033[0m", // Blue
		"ruby":             "\033[1;31m\033[0m", // Red
		"rust":             "\033[1;31m\033[0m", // Red
		"scala":            "\033[1;31m\033[0m", // Red
		"scheme":           "\033[1;32m\033[0m", // Green
		"shell":            "\033[1;32m\033[0m", // Green
		"sql":              "\033[1;34m\033[0m", // Blue
		"swift":            "\033[1;31m\033[0m", // Red
		"typescript":       "\033[1;36m\033[0m", // Cyan
		"typescript-react": "\033[1;36m\033[0m", // Cyan
		"visual-basic":     "\033[1;34m\033[0m", // Blue
		"xml":              "\033[1;33m\033[0m", // Yellow
		"yaml":             "\033[1;33m\033[0m", // Yellow
		"directory":        "\033[1;34m\033[0m", // Blue folder
		"openDirectory":    "\033[1;34m\033[0m", // Blue open folder
		"emptyDirectory":   "\033[1;34m\033[0m", // Blue empty folder
		"unknown":          "\033[1;30m\033[0m", // Grey
	}
	reverseLookupFt = make(map[string]string)
	mu              sync.RWMutex
)

func init() {
	for ext, ft := range extensionMap {
		reverseLookupFt[ft] = ext
	}
}

func FTbyFileName(fileName string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	extension := filepath.Ext(fileName)
	ft, exists := extensionMap[extension]

	if exists {
		return ft, nil
	}

	return "unknown", fmt.Errorf("file type not recognized for extension: %s", extension)
}

func FTbyExtension(extension string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	ft, exists := extensionMap[extension]

	if exists {
		return ft, nil
	}

	return "unknown", fmt.Errorf("file type not recognized for extension: %s", extension)
}

func ExtensionByFT(fileType string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	extension, exists := reverseLookupFt[fileType]
	if exists {
		return extension, nil
	}

	return "unknown", fmt.Errorf("extension not recognized for file type: %s", fileType)
}

func RegisterFileType(extension, fileType string) {
	mu.Lock()
	defer mu.Unlock()
	extensionMap[extension] = fileType
	reverseLookupFt[fileType] = extension
}

func RegisterIcons(fileType, icon string) {
	mu.Lock()
	defer mu.Unlock()
	fileTypeIcons[fileType] = icon
}

func MatchExtension(fileName, extension string) bool {
	if len(fileName) > len(extension) &&
		fileName[len(fileName)-len(extension):] == extension {
		return true
	}

	return false
}

func GetFtIcon(fileType string) string {
	if icon, exists := fileTypeIcons[fileType]; exists {
		return icon
	}
	return fileTypeIcons["unknown"]
}

func GetExtension(fileName string) string {
	if MatchExtension(fileName, config.MetaExtension) {
		return config.MetaExtension
	}
	return filepath.Ext(fileName)
}
