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
	filePanel
	snippetPanel
	contentPanel
)

type app struct {
	config *config.Config

	snippetIndex *models.SnippetIndex
	treeRoot     *filesystem.TreeNode

	leftColumn  ui.Container
	rightColumn ui.Container

	focusPanel focusState
	panels     map[focusState]tea.Model
}

func NewApp(c *config.Config, snippetIndex *models.SnippetIndex, treeRoot *filesystem.TreeNode) app {
	return app{
		config: c,

		snippetIndex: snippetIndex,
		treeRoot:     treeRoot,

		focusPanel:  0,
		leftColumn:  ui.NewContainer(),
		rightColumn: ui.NewContainer(),

		panels: map[focusState]tea.Model{
			langPanel:    ui.NewLangPanel(),
			filePanel:    ui.NewFilePanel(treeRoot, c),
			snippetPanel: ui.NewSnippePanel(),
			contentPanel: ui.NewContentPanel(),
		},
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.panels[langPanel].Init(),
		m.panels[filePanel].Init(),
		m.panels[snippetPanel].Init(),
		m.panels[contentPanel].Init(),
		func() tea.Msg {
			return initMsg{}
		},
	)
}

type initMsg struct{}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		for key, panel := range m.panels {
			var cmd tea.Cmd
			m.panels[key], cmd = panel.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "l", "right":
			m.focusPanel++
			if m.focusPanel > snippetPanel {
				m.focusPanel = langPanel
			}
		case "shift+tab", "h", "left":
			m.focusPanel--
			// we use > instead of < because focusPanel is an unsigned int and will wrap around to max value
			if m.focusPanel > snippetPanel {
				m.focusPanel = snippetPanel
			}
		// NOTE: this are all propagated to focused pannel
		case "j", "down", "k", "up":
			if panel, ok := m.panels[m.focusPanel]; ok {
				var cmd tea.Cmd
				m.panels[m.focusPanel], cmd = panel.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m app) View() string {

	langStyle := m.config.Theme.LangPanelStyle
	filesStyle := m.config.Theme.FilePanelStyle
	snippetStyle := m.config.Theme.SnippetPanelStyle

	switch m.focusPanel {
	case langPanel:
		langStyle = m.config.Theme.FocusPanel(langStyle)
	case filePanel:
		filesStyle = m.config.Theme.FocusPanel(filesStyle)
	case snippetPanel:
		snippetStyle = m.config.Theme.FocusPanel(snippetStyle)
	}

	langString := langStyle.Render(m.panels[langPanel].View())
	fileString := filesStyle.Render(m.panels[filePanel].View())
	snippetString := snippetStyle.Render(m.panels[snippetPanel].View())

	leftContent := lipgloss.JoinVertical(lipgloss.Top, langString, fileString, snippetString)
	m.leftColumn.SetContent(leftContent)

	rightContent := m.config.Theme.ContentPanelStyle.Render(m.panels[contentPanel].View())
	m.rightColumn.SetContent(rightContent)

	s := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	return s
}
