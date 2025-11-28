package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/filesystem"
	"github.com/kaputi/navani/internal/utils"
)

type FilePanel struct {
	fileTree    *filesystem.FileTree
	cursor      int
	start       int
	end         int
	fileStrings []string
}

func NewFilePanel(fileTree *filesystem.FileTree) FilePanel {
	t := config.Theme()
	return FilePanel{
		fileTree: fileTree,
		cursor:   0,
		start:    0,
		end:      t.FilePanelHeight,
	}
}

func (f FilePanel) getNode() *filesystem.TreeNode {
	nodeList := f.fileTree.OpenNodeList()
	if len(nodeList)-1 >= f.cursor {
		return nodeList[f.cursor]
	}
	return nil
}

func (f FilePanel) Init() tea.Cmd {
	return nil
}

func (f FilePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case WindowResizeMsg:
		currSize := f.end - f.start
		nextSize := config.Theme().FilePanelHeight
		diff := currSize - nextSize
		f.end -= diff
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			f.cursor++
		case "k", "up":
			if f.cursor > 0 {
				f.cursor--
			}
		case "e":
			node := f.getNode()
			if node != nil && !node.IsDir() {
				cmds = append(cmds, EditModeMsg)
			}
		case "enter", " ":
			node := f.getNode()
			if node != nil && node.IsDir() {
				node.Toggle()
				f.fileTree.UpdateOpenNodeList()
			}
		}
	}

	fileStrings := f.fileTree.Strings()
	totalStrs := len(fileStrings)

	if f.cursor > totalStrs-1 {
		f.cursor = totalStrs - 1
	}

	if f.cursor > f.end-1 {
		diff := f.cursor - f.end + 1
		f.end += diff
		f.start += diff
	}

	if f.cursor < f.start {
		diff := f.start - f.cursor
		f.start -= diff
		f.end -= diff
	}

	// TODO: highlight selected line with the theme
	fileStrings[f.cursor] = fileStrings[f.cursor] + " <=="

	f.fileStrings = fileStrings[f.start:utils.Min(totalStrs, f.end)]

	openList := f.fileTree.OpenNodeList()
	if len(openList) >= f.cursor {
		node := openList[f.cursor]
		if !node.IsDir() {
			cmds = append(cmds, ContentMsg(node.Path()))
		} else {
			cmds = append(cmds, ContentMsg(""))
		}
	}

	return f, tea.Batch(cmds...)
}

func (f FilePanel) View() string {
	return strings.Join(f.fileStrings, "\n")
}
