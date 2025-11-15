package ui

import tea "github.com/charmbracelet/bubbletea"

type SnippetList struct {
}

func NewSnippetList() SnippetList {
	return SnippetList{}
}

func (s SnippetList) Init() tea.Cmd {
	return nil
}

func (s SnippetList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s SnippetList) View() string {
	return "Snippet list panel (not implemented yet)"
}
