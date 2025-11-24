package app

import tea "github.com/charmbracelet/bubbletea"

type LangPanel struct {
}

func NewLangPanel() LangPanel {
	return LangPanel{}
}

func (l LangPanel) Init() tea.Cmd {
	return nil
}

func (l LangPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l LangPanel) View() string {
	return "Language panel (not implemented yet)"
}
