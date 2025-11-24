package app

import tea "github.com/charmbracelet/bubbletea"

type ContentPanel struct {
	filePath string
}

func NewContentPanel() ContentPanel {
	return ContentPanel{
		filePath: "",
	}
}

func (c ContentPanel) Init() tea.Cmd {
	return nil
}

func (c ContentPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case contentMsg:
		fp := msg.filePath
		c.filePath = fp
	}

	return c, nil
}

func (c ContentPanel) View() string {
	if c.filePath != "" {
		return "content panel for snippet: " + c.filePath
	}

	return "content panel (no snippet selected)"
}
