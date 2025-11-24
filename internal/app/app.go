package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	fileTree     *filesystem.FileTree

	leftColumn  Container
	rightColumn Container

	focusPanel focusState
	panels     map[focusState]tea.Model
}

func NewApp(c *config.Config, snippetIndex *models.SnippetIndex, fileTree *filesystem.FileTree) app {
	return app{
		config: c,

		snippetIndex: snippetIndex,
		fileTree:     fileTree,

		focusPanel:  0,
		leftColumn:  NewContainer(),
		rightColumn: NewContainer(),

		panels: map[focusState]tea.Model{
			langPanel:    NewLangPanel(),
			filePanel:    NewFilePanel(fileTree, c),
			snippetPanel: NewSnippePanel(),
			contentPanel: NewContentPanel(),
		},
	}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(
		a.panels[langPanel].Init(),
		a.panels[filePanel].Init(),
		a.panels[snippetPanel].Init(),
		a.panels[contentPanel].Init(),
		InitMsg,
	)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case initMsg:
		for key, panel := range a.panels {
			var cmd tea.Cmd
			a.panels[key], cmd = panel.Update(msg)
			cmds = append(cmds, cmd)
		}

	case contentMsg:
		var cmd tea.Cmd
		a.panels[contentPanel], cmd = a.panels[contentPanel].Update(msg)
		cmds = append(cmds, cmd)
		snp, ok := a.snippetIndex.ByFilePath[msg.filePath]
		if ok {
			metadata := snp.Metadata
			metadataStrs := metadata.Strings()
			a.panels[snippetPanel], cmd = a.panels[snippetPanel].Update(snippetMetadataMsg{metadataStrings: metadataStrs})
			cmds = append(cmds, cmd)
		} else {
			a.panels[snippetPanel], cmd = a.panels[snippetPanel].Update(snippetMetadataMsg{metadataStrings: []string{"No snippet selected"}})
			cmds = append(cmds, cmd)

		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "tab", "l", "right":
			a.focusPanel++
			if a.focusPanel > snippetPanel {
				a.focusPanel = langPanel
			}
		case "shift+tab", "h", "left":
			a.focusPanel--
			// we use > instead of < because focusPanel is an unsigned int and will wrap around to max value
			if a.focusPanel > snippetPanel {
				a.focusPanel = snippetPanel
			}
		// NOTE: this are all propagated to focused pannel
		case "j", "down", "k", "up", "enter", "backspace", " ":
			if panel, ok := a.panels[a.focusPanel]; ok {
				var cmd tea.Cmd
				a.panels[a.focusPanel], cmd = panel.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	langStyle := a.config.Theme.LangPanelStyle
	filesStyle := a.config.Theme.FilePanelStyle
	snippetStyle := a.config.Theme.SnippetPanelStyle

	switch a.focusPanel {
	case langPanel:
		langStyle = a.config.Theme.FocusPanel(langStyle)
	case filePanel:
		filesStyle = a.config.Theme.FocusPanel(filesStyle)
	case snippetPanel:
		snippetStyle = a.config.Theme.FocusPanel(snippetStyle)
	}

	langString := langStyle.Render(a.panels[langPanel].View())
	fileString := filesStyle.Render(a.panels[filePanel].View())
	snippetString := snippetStyle.Render(a.panels[snippetPanel].View())

	leftContent := lipgloss.JoinVertical(lipgloss.Top, langString, fileString, snippetString)
	a.leftColumn.SetContent(leftContent)

	rightContent := a.config.Theme.ContentPanelStyle.Render(a.panels[contentPanel].View())
	a.rightColumn.SetContent(rightContent)

	s := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	return s
}
