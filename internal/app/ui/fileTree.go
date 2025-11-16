package ui

import tea "github.com/charmbracelet/bubbletea"

type Tree struct {
}

func NewTree() Tree {
	return Tree{}
}

func (t Tree) Init() tea.Cmd {
	return nil
}

func (t Tree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t Tree) View() string {
	return "Tree panel (not implemented yet)"
}
