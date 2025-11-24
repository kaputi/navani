package ui

import tea "github.com/charmbracelet/bubbletea"

type Content struct {
}

func NewContent() Content {
	return Content{}
}

func (c Content) Init() tea.Cmd {
	return nil
}

func (c Content) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c Content) View() string {
	return "content panel (not implemented yet)"
}
