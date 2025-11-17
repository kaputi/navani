package ui

import tea "github.com/charmbracelet/bubbletea"

type ContentPanel struct {
}

func NewContentPanel() ContentPanel {
	return ContentPanel{}
}

func (c ContentPanel) Init() tea.Cmd {
	return nil
}

func (c ContentPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c ContentPanel) View() string {
	return "content panel (not implemented yet)"
}
