package ui

import tea "github.com/charmbracelet/bubbletea"

type Lang struct {
}

func NewLang() Lang {
	return Lang{}
}

func (l Lang) Init() tea.Cmd {
	return nil
}

func (l Lang) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l Lang) View() string {
	return "Language panel (not implemented yet)"
}
