package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaputi/navani/internal/filesystem"
)

type TreePanel struct {
	root *filesystem.TreeNode
}

func NewTreePanel(root *filesystem.TreeNode) TreePanel {
	return TreePanel{
		root,
	}
}

func (t TreePanel) Init() tea.Cmd {
	return nil
}

func (t TreePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t TreePanel) View() string {
	treeStrings := t.root.Strings()
	str := strings.Join(treeStrings, "\n")
	return str
}
