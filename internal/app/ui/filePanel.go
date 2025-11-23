package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/filesystem"
	"github.com/kaputi/navani/internal/utils"
)

type FilePanel struct {
	config      *config.Config
	root        *filesystem.TreeNode
	cursor      int
	start       int
	end         int
	fileStrings []string
}

func NewFilePanel(root *filesystem.TreeNode, c *config.Config) FilePanel {
	return FilePanel{
		config: c,
		root:   root,
		cursor: 0,
		start:  0,
		end:    c.Theme.FilePanelHeight,
	}
}

func (f FilePanel) Init() tea.Cmd {
	return nil
}

func (f FilePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			f.cursor++
		case "k", "up":
			if f.cursor > 0 {
				f.cursor--
			}
		}
	}

	fileStrings := f.root.Strings()
	totalStrs := len(fileStrings)

	if f.cursor > totalStrs {
		f.cursor = totalStrs
	}

	if f.cursor > f.end {
		diff := f.cursor - f.end
		f.end += diff
		f.start += diff
	}

	if f.cursor < f.start {
		diff := f.start - f.cursor
		f.start -= diff
		f.end -= diff
	}

	fileStrings[f.cursor] = "=>>" + fileStrings[f.cursor]

	f.fileStrings = fileStrings[f.start:utils.Min(totalStrs, f.end)]

	return f, nil
}

func (f FilePanel) View() string {
	return strings.Join(f.fileStrings, "\n")
}
