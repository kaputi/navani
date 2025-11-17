package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaputi/navani/internal/app/ui"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/filesystem"
	"github.com/kaputi/navani/internal/models"
)

type focusState uint

const (
	langPanel focusState = iota
	treePanel
	snippetPanel
	contentPanel
)

type app struct {
	config *config.Config

	snippetIndex *models.SnippetIndex
	treeRoot     *filesystem.TreeNode

	focusPanel   focusState
	leftColumn   ui.Container
	rightColumn  ui.Container
	langPanel    ui.LangPanel
	treePanel    ui.TreePanel
	snippePanel  ui.SnippetPanel
	contentPanel ui.ContentPanel
}

func NewApp(c *config.Config, snippetIndex *models.SnippetIndex, treeRoot *filesystem.TreeNode) app {
	return app{
		config: c,

		snippetIndex: snippetIndex,
		treeRoot:     treeRoot,

		focusPanel:   0,
		leftColumn:   ui.NewContainer(),
		rightColumn:  ui.NewContainer(),
		langPanel:    ui.NewLangPanel(),
		treePanel:    ui.NewTreePanel(treeRoot),
		snippePanel:  ui.NewSnippePanel(),
		contentPanel: ui.NewContentPanel(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.langPanel.Init(),
		m.treePanel.Init(),
		m.snippePanel.Init(),
		m.contentPanel.Init(),
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

	langString := langStyle.Render(m.langPanel.View())
	treeString := treeStyle.Render(m.treePanel.View())
	snippetString := snippetStyle.Render(m.snippePanel.View())

	leftContent := lipgloss.JoinVertical(lipgloss.Top, langString, treeString, snippetString)
	m.leftColumn.SetContent(leftContent)

	rightContent := m.config.Theme.ContentPanelStyle.Render(m.contentPanel.View())
	m.rightColumn.SetContent(rightContent)

	s := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	return s
}
