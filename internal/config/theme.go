package config

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/kaputi/navani/internal/utils/logger"
	"golang.org/x/term"
)

var (
	defaultPalette = map[string]string{
		"foreground": "#FFFFFF",
		"background": "#000000",
		"primary":    "#FFA500",
		"secondary":  "#00FFFF",
		"accent":     "#FF00FF",
		"selected":   "#008000",
	}
	defaultWidthPercent float32 = 0.25
	defaultPanelPadding [2]int  = [2]int{0, 1}
)

type theme struct {
	palette              map[string]string
	navPanelWidthPercent float32
	navPanelPadding      [2]int

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

func watchWindowResize(t *theme) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	for range sig {
		t.updateStyle()
	}
}

func getTermSize() (int, int) {
	width, height, err := term.GetSize(0) // Get terminal size (0 is stdin)
	if err != nil {
		logger.Log("Error getting terminal size")
	}
	return width, height
}

// defaults
func newTheme() *theme {
	return &theme{
		palette: defaultPalette,

		navPanelWidthPercent: defaultWidthPercent,
		navPanelPadding:      defaultPanelPadding,
	}
}

func (t *theme) updateStyle() {

	width, height := getTermSize()

	TermWidth := width
	TermHeight := height

	// NOTE:  style widths dont include padding and border
	xOffset := 2 + t.navPanelPadding[1]*2
	yOffset := 2 + t.navPanelPadding[0]*2

	navPanelWidth := int(float32(TermWidth)*t.navPanelWidthPercent) - xOffset
	navPanelHeight := int((float32(TermHeight))*0.3) - yOffset

	t.LangPanelHeight = navPanelHeight
	t.FilePanelHeight = navPanelHeight
	t.SnippetPanelHeight = TermHeight - t.LangPanelHeight - t.FilePanelHeight - yOffset*3

	t.ContentPanelWidth = TermWidth - navPanelWidth - xOffset
	t.ContentPanelHeight = t.LangPanelHeight + t.FilePanelHeight + t.SnippetPanelHeight + yOffset*2

	t.PanelStyle = lipgloss.NewStyle().
		Margin(0, 0).
		Padding(t.navPanelPadding[0], t.navPanelPadding[1]).
		Border(lipgloss.RoundedBorder())

	t.LangPanelStyle = t.PanelStyle.Width(navPanelWidth).Height(t.LangPanelHeight)
	t.FilePanelStyle = t.PanelStyle.Width(navPanelWidth).Height(t.FilePanelHeight)
	t.SnippetPanelStyle = t.PanelStyle.Width(navPanelWidth).Height(t.SnippetPanelHeight)
	t.ContentPanelStyle = t.PanelStyle.Width(t.ContentPanelWidth).Height(t.ContentPanelHeight)
}

func (t *theme) init() {
	// TODO: load or pass theme from config file
	t.updateStyle()
	go watchWindowResize(t)
}

func (t *theme) Color(name string) (string, error) {
	if color, exists := t.palette[name]; exists {
		return color, nil
	}
	return "", fmt.Errorf("color %s not found", name)
}

func (t *theme) FocusPanel(inputStyle lipgloss.Style) lipgloss.Style {
	accentColor, _ := t.Color("accent")
	return inputStyle.BorderForeground(lipgloss.Color(accentColor))
}

func (t *theme) UpdatePalette(palette map[string]string) {
	for key, value := range palette {
		if _, exists := t.palette[key]; exists {
			t.palette[key] = value
		}
	}
}
