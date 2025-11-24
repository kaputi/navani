package theme

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/kaputi/navani/internal/utils/logger"
)

var colors = map[string]string{
	"foreground": "#FFFFFF",
	"background": "#000000",
	"primary":    "#FFA500",
	"secondary":  "#00FFFF",
	"accent":     "#FF00FF",
	"selected":   "#008000",
}

// theme variables
var (
	navPanelWidthPercent float32 = 0.25
	navPanelPadding      [2]int  = [2]int{0, 1}
)

// calculated variable
var (
	TermWidth  int
	TermHeight int

	NavPanelWidth int

	LangPanelHeight    int
	TreePanelHeight    int
	SnippetPanelHeight int

	ContentPanelWidth  int
	ContentPanelHeight int

	PanelStyle        lipgloss.Style
	LangPanelStyle    lipgloss.Style
	TreePanelStyle    lipgloss.Style
	SnippetPanelStyle lipgloss.Style
	ContentPanelStyle lipgloss.Style
)

func Init() {
	updateStyle()
	go watchWindowResize()
}

func updateStyle() {
	width, height := getTermSize()

	TermWidth = width
	TermHeight = height

	// NOTE:  style widths dont include padding and border
	xOffset := 2 + navPanelPadding[1]*2
	yOffset := 2 + navPanelPadding[0]*2

	navPanelWidth := int(float32(TermWidth)*navPanelWidthPercent) - xOffset
	navPanelHeight := int((float32(TermHeight))*0.3) - yOffset

	LangPanelHeight = navPanelHeight
	TreePanelHeight = navPanelHeight
	SnippetPanelHeight = TermHeight - LangPanelHeight - TreePanelHeight - yOffset*3

	ContentPanelWidth = TermWidth - navPanelWidth - xOffset
	ContentPanelHeight = LangPanelHeight + TreePanelHeight + SnippetPanelHeight + yOffset*2

	PanelStyle = lipgloss.NewStyle().
		Margin(0, 0).
		Padding(navPanelPadding[0], navPanelPadding[1]).
		Border(lipgloss.RoundedBorder())

	LangPanelStyle = PanelStyle.Width(navPanelWidth).Height(LangPanelHeight)
	TreePanelStyle = PanelStyle.Width(navPanelWidth).Height(TreePanelHeight)
	SnippetPanelStyle = PanelStyle.Width(navPanelWidth).Height(SnippetPanelHeight)
	ContentPanelStyle = PanelStyle.Width(ContentPanelWidth).Height(ContentPanelHeight)
}

func getTermSize() (int, int) {
	width, height, err := term.GetSize(0) // Get terminal size (0 is stdin)
	if err != nil {
		logger.Log("Error getting terminal size")
	}
	return width, height
}

func watchWindowResize() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	for range sig {
		updateStyle()
	}
}

func Color(name string) (string, error) {
	if color, exists := colors[name]; exists {
		return color, nil
	}

	return "", fmt.Errorf("color %s not found", name)
}

func FocusPanel(inputStyle lipgloss.Style) lipgloss.Style {
	accentColor, _ := Color("accent")
	return inputStyle.BorderForeground(lipgloss.Color(accentColor))
}
