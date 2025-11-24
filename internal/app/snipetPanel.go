package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SnippetPanel struct {
	metadtaStrings []string
}

func NewSnippePanel() SnippetPanel {
	return SnippetPanel{
		metadtaStrings: []string{"TEEEST"},
	}
}

func (s SnippetPanel) Init() tea.Cmd {
	return nil
}

func (s SnippetPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case snippetMetadataMsg:
		s.metadtaStrings = msg.metadataStrings
	}
	return s, nil
}

func (s SnippetPanel) View() string {
	return strings.Join(s.metadtaStrings, "\n")
}
