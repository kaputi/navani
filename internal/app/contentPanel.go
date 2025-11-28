package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

type ContentPanel struct {
	content  string
	filepath string
}

func NewContentPanel() ContentPanel {
	return ContentPanel{}
}

func (c ContentPanel) Init() tea.Cmd {
	return nil
}

func (c ContentPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case contentMsg:

		fileContent, err := os.ReadFile(msg.filePath)

		if err != nil {
			// TODO: better error handling
			c.content = "Error reading file: " + err.Error()
			logger.Err(fmt.Errorf("error reading file: %w", err))
			return c, nil
		}

		fileText := string(fileContent)

		fileName := filepath.Base(msg.filePath)
		language, err := utils.FTbyFileName(fileName)

		c.filepath = msg.filePath

		if err == nil {
			fileText = highlightCode(string(fileContent), language)
		}

		c.content = fileText
	case editModeMsg:
		cmd = c.openInEditor()
	}

	return c, cmd
}

func (c ContentPanel) View() string {

	if c.content != "" {
		return c.content
	}

	return "(no snippet selected)"
}

func (c ContentPanel) openInEditor() tea.Cmd {
	if c.filepath == "" {
		return func() tea.Msg {
			return errorMsg{err: fmt.Errorf("no file selected")}
		}
	}

	editor, args := getEditor()

	cmdArgs := append(args, c.filepath)

	return tea.ExecProcess(exec.Command(editor, cmdArgs...), func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: fmt.Errorf("failed to open editor: %w", err)}
		}

		return contentMsg{filePath: c.filepath}
	})
}

// getEditor returns the appropriate editor for the platform
// TODO: read from config if available
func getEditor() (string, []string) {
	// Try $EDITOR first (Unix/Linux/Mac)
	if editor := os.Getenv("EDITOR"); editor != "" {
		parts := strings.Fields(editor)
		if len(parts) > 1 {
			editorCmd := parts[0]
			existingArgs := parts[1:]

			waitEditor, waitArgs := addWaitFlag(editorCmd)

			allArgs := append(waitArgs, existingArgs...)
			return waitEditor, allArgs
		}
		return addWaitFlag(editor)
	}

	// Try $VISUAL (Unix alternative)
	if visual := os.Getenv("VISUAL"); visual != "" {
		parts := strings.Fields(visual)
		if len(parts) > 1 {
			editorCmd := parts[0]
			existingArgs := parts[1:]
			waitEditor, waitArgs := addWaitFlag(editorCmd)
			allArgs := append(existingArgs, waitArgs...)
			return waitEditor, allArgs
		}
		return addWaitFlag(visual)
	}

	// Platform-specific defaults
	switch runtime.GOOS {
	case "windows":
		// Try common Windows editors
		if _, err := exec.LookPath("code"); err == nil {
			return "code", []string{"--wait"} // VS Code
		}

		return "notepad", []string{} // Fallback to notepad

	case "darwin": // macOS
		// Prefer vim/nano over vi
		if _, err := exec.LookPath("vim"); err == nil {
			return "vim", []string{}
		}
		return "nano", []string{}

	default: // Linux and other Unix-like
		// Try to find best available editor
		if _, err := exec.LookPath("vim"); err == nil {
			return "vim", []string{}
		}
		if _, err := exec.LookPath("nano"); err == nil {
			return "nano", []string{}
		}
		return "vi", []string{} // vi should always exist on Unix
	}
}

func addWaitFlag(editor string) (string, []string) {
	editorLower := strings.ToLower(filepath.Base(editor))

	// Map of editors that need wait flags
	waitFlags := map[string][]string{
		"code":          {"--wait"},
		"code-insiders": {"--wait"},
		"subl":          {"--wait"},
		"sublime":       {"--wait"},
		"sublime_text":  {"--wait"},
		"atom":          {"--wait"},
		"gedit":         {"--wait"},
		"kate":          {"--block"},
		"notepad++.exe": {"-multiInst", "-notabbar", "-nosession"},
		"notepad++":     {"-multiInst", "-notabbar", "-nosession"},
	}

	if flags, ok := waitFlags[editorLower]; ok {
		return editor, flags
	}

	// Default: assume it waits (most terminal editors do)
	return editor, []string{}
}
