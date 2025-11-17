package ui

import tea "github.com/charmbracelet/bubbletea"

type SnippetPanel struct {
}

func NewSnippePanel() SnippetPanel {
	return SnippetPanel{}
}

func (s SnippetPanel) Init() tea.Cmd {
	return nil
}

func (s SnippetPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s SnippetPanel) View() string {
	return "Snippet list panel (not implemented yet)"
}
