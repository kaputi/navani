package utils

import (
	"fmt"
	"path/filepath"
	"sync"
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
	if _, exists := extensionMap[extension]; exists {
		return
	}
	extensionMap[extension] = fileType
	reverseLookupFt[fileType] = extension
}

func MatchExtension(fileName, extension string) bool {
	if len(fileName) > len(extension) &&
		fileName[len(fileName)-len(extension):] == extension {
		return true
	}

	return false
}
