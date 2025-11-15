package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaputi/navani/internal/app/ui"
	"github.com/kaputi/navani/internal/config"
)

type focusState uint

const (
	langPanel focusState = iota
	treePanel
	snippetPanel
	contentPanel
)

type app struct {
	focusPanel focusState

	config *config.Config

	leftColumn  ui.Container
	rightColumn ui.Container
	langUI      ui.Lang
	treeUI      ui.Tree
	snippetUI   ui.SnippetList
	contenUI    ui.Content
}

func NewApp(c *config.Config) app {
	return app{
		focusPanel: 0,

		config: c,

		leftColumn:  ui.NewContainer(),
		rightColumn: ui.NewContainer(),
		langUI:      ui.NewLang(),
		treeUI:      ui.NewTree(),
		snippetUI:   ui.NewSnippetList(),
		contenUI:    ui.NewContent(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.langUI.Init(),
		m.treeUI.Init(),
		m.snippetUI.Init(),
		m.contenUI.Init(),
	)
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "j", "down":
			m.focusPanel++
			if m.focusPanel > snippetPanel {
				m.focusPanel = langPanel
			}
		case "shift+tab", "k", "up":
			m.focusPanel--
			// we use > instead of < because focusPanel is an unsigned int and will wrap around to max value
			if m.focusPanel > snippetPanel {
				m.focusPanel = snippetPanel
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m app) View() string {

	langStyle := m.config.Theme.LangPanelStyle
	treeStyle := m.config.Theme.TreePanelStyle
	snippetStyle := m.config.Theme.SnippetPanelStyle

	switch m.focusPanel {
	case langPanel:
		langStyle = m.config.Theme.FocusPanel(langStyle)
	case treePanel:
		treeStyle = m.config.Theme.FocusPanel(treeStyle)
	case snippetPanel:
		snippetStyle = m.config.Theme.FocusPanel(snippetStyle)
	}

	langString := langStyle.Render(m.langUI.View())
	treeString := treeStyle.Render(m.treeUI.View())
	snippetString := snippetStyle.Render(m.snippetUI.View())

	leftContent := lipgloss.JoinVertical(lipgloss.Top, langString, treeString, snippetString)
	m.leftColumn.SetContent(leftContent)

	rightContent := m.config.Theme.ContentPanelStyle.Render(m.contenUI.View())
	m.rightColumn.SetContent(rightContent)

	s := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	return s
}
