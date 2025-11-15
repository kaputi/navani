package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Container struct {
	content string
}

func (c Container) Init() tea.Cmd {
	return nil
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *Container) SetContent(content string) {
	c.content = content
}

func (c Container) View() string {
	return c.content
}

func NewContainer() Container {
	return Container{content: ""}
}
