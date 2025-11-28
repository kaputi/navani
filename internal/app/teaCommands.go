package app

import tea "github.com/charmbracelet/bubbletea"

type initMsg struct{}
type snippetMetadataMsg struct{ metadataStrings []string }
type editModeMsg struct{}
type contentMsg struct{ filePath string }
type errorMsg struct{ err error }
type WindowResizeMsg struct{}

func InitMsg() tea.Msg {
	return initMsg{}
}

func ContentMsg(filePath string) func() tea.Msg {
	return func() tea.Msg {
		return contentMsg{filePath: filePath}
	}
}

func EditModeMsg() tea.Msg {
	return editModeMsg{}
}
