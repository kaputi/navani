package config

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kaputi/navani/internal/utils/logger"
	"golang.org/x/term"
)

var (
	localTheme theme
)

type treeChars struct {
	TreeOpenChar       string
	TreeCloseChar      string
	TreeIndentChar     string
	TreeDirIndentChar  string
	TreeLastIndentChar string
	TreeIndentSize     int
}

type theme struct {
	Tree treeChars

	Palette              map[string]string
	NavPanelWidthPercent float32
	NavPanelPadding      [2]int

	PanelStyle        lipgloss.Style
	LangPanelStyle    lipgloss.Style
	FilePanelStyle    lipgloss.Style
	SnippetPanelStyle lipgloss.Style
	ContentPanelStyle lipgloss.Style

	LangPanelHeight    int
	FilePanelHeight    int
	SnippetPanelHeight int
	ContentPanelWidth  int
	ContentPanelHeight int
}

func (t theme) Color(name string) (string, error) {
	if color, exists := t.Palette[name]; exists {
		return color, nil
	}
	return "", fmt.Errorf("color %s not found", name)
}

func (t theme) FocusPanel(inputStyle lipgloss.Style) lipgloss.Style {
	accentColor, _ := t.Color("accent")
	return inputStyle.BorderForeground(lipgloss.Color(accentColor))
}

func defaultTheme() theme {
	return theme{
		Tree: treeChars{
			TreeOpenChar:       "▼ ",
			TreeCloseChar:      "▶ ",
			TreeIndentChar:     "│",
			TreeDirIndentChar:  "├",
			TreeLastIndentChar: "└",
			TreeIndentSize:     2,
		},

		Palette: map[string]string{
			"foreground": "#FFFFFF",
			"background": "#000000",
			"primary":    "#FFA500",
			"secondary":  "#00FFFF",
			"accent":     "#FF00FF",
			"selected":   "#008000",
		},
		NavPanelWidthPercent: 0.25,
		NavPanelPadding:      [2]int{0, 1},
	}
}

func init() {
	localTheme = defaultTheme()
	UpdateStyles()
}

func getTermSize() (int, int) {
	width, height, err := term.GetSize(0) // Get terminal size (0 is stdin)
	if err != nil {
		logger.Log("Error getting terminal size")
	}
	return width, height
}

func UpdateStyles() {

	width, height := getTermSize()

	TermWidth := width
	TermHeight := height

	// NOTE:  style widths dont include padding and border
	xOffset := 2 + localTheme.NavPanelPadding[1]*2
	yOffset := 2 + localTheme.NavPanelPadding[0]*2

	navPanelWidth := int(float32(TermWidth)*localTheme.NavPanelWidthPercent) - xOffset
	navPanelHeight := int((float32(TermHeight))*0.3) - yOffset

	localTheme.LangPanelHeight = navPanelHeight
	localTheme.FilePanelHeight = navPanelHeight
	localTheme.SnippetPanelHeight = TermHeight - localTheme.LangPanelHeight - localTheme.FilePanelHeight - yOffset*3

	localTheme.ContentPanelWidth = TermWidth - navPanelWidth - xOffset
	localTheme.ContentPanelHeight = localTheme.LangPanelHeight + localTheme.FilePanelHeight + localTheme.SnippetPanelHeight + yOffset*2

	localTheme.PanelStyle = lipgloss.NewStyle().
		Margin(0, 0).
		Padding(localTheme.NavPanelPadding[0], localTheme.NavPanelPadding[1]).
		Border(lipgloss.RoundedBorder())

	localTheme.LangPanelStyle = localTheme.PanelStyle.Width(navPanelWidth).Height(localTheme.LangPanelHeight)
	localTheme.FilePanelStyle = localTheme.PanelStyle.Width(navPanelWidth).Height(localTheme.FilePanelHeight)
	localTheme.SnippetPanelStyle = localTheme.PanelStyle.Width(navPanelWidth).Height(localTheme.SnippetPanelHeight)
	localTheme.ContentPanelStyle = localTheme.PanelStyle.Width(localTheme.ContentPanelWidth).Height(localTheme.ContentPanelHeight)
}

func Theme() theme {
	return localTheme
}

func LoadTheme() {
	// TODO: reads the theme from the config directory (not configurable, depends on os)

	UpdateStyles()
}
