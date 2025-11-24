package app

import tea "github.com/charmbracelet/bubbletea"

type initMsg struct{}

func InitMsg() tea.Msg {
	return initMsg{}
}

type contentMsg struct {
	filePath string
}

func ContentMsg(filePath string) func() tea.Msg {
	return func() tea.Msg {
		return contentMsg{filePath: filePath}
	}
}

type snippetMetadataMsg struct {
	metadataStrings []string
}
