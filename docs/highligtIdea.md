# Recommended Libraries
- **Chroma (Go port of Pygments)**:
    Chroma is a popular Go library for syntax highlighting. You can use it to highlight code snippets by language and output ANSI color codes for terminal display.​
- **charmbracelet/glamour**:
Glamour is a Markdown renderer for Go that supports syntax highlighting using Chroma. You can use it to render code blocks with proper colors.​
- **Custom Highlighting with Lipgloss**:
    Lipgloss (from Charm Bracelet) lets you style terminal output, and you can combine it with Chroma to apply colors to highlighted code.​

# How to Integrate
1. Parse the snippet content and detect its language (from the file extension).
2. Use Chroma to highlight the code, generating ANSI-colored output.
3. Render the highlighted code in your Bubble Tea view using lipgloss or similar styling.

# Example Workflow
```go
import (
    "github.com/alecthomas/chroma/formatters/terminal"
    "github.com/alecthomas/chroma/lexers"
    "github.com/alecthomas/chroma/styles"
    "github.com/charmbracelet/lipgloss"
)

func highlightCode(code, language string) string {
    lexer := lexers.Get(language)
    if lexer == nil {
        lexer = lexers.Fallback
    }
    style := styles.Get("dracula")
    formatter := terminal.NewFormatter(style)
    iterator, _ := lexer.Tokenise(nil, code)
    var buf bytes.Buffer
    formatter.Format(&buf, iterator)
    return buf.String()
}
```

