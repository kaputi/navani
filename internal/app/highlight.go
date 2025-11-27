package app

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func highlightCode(code, language string) string {
	lexer := lexers.Get(language)

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("dracula")
	if style == nil {

		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, _ := lexer.Tokenise(nil, code)
	if iterator == nil {
		return code
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}
