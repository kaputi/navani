package app

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

type ContentPanel struct {
	content string
}

func NewContentPanel() ContentPanel {
	return ContentPanel{content: ""}
}

func (c ContentPanel) Init() tea.Cmd {
	return nil
}

func (c ContentPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case contentMsg:

		fileContent, err := os.ReadFile(msg.filePath)

		if err != nil {
			// TODO: better error handling
			c.content = "Error reading file: " + err.Error()
			logger.Err(fmt.Errorf("error reading file: %w", err))
			return c, nil
		}

		fileText := string(fileContent)

		fileName := filepath.Base(msg.filePath)
		language, err := utils.FTbyFileName(fileName)
		if err == nil {
			fileText = highlightCode(string(fileContent), language)
		}

		c.content = fileText
	}

	return c, nil
}

func (c ContentPanel) View() string {
	if c.content != "" {
		return c.content
	}

	return "content panel (no snippet selected)"
}
