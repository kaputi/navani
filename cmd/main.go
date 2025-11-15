package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kaputi/navani/internal/app"
	"github.com/kaputi/navani/internal/config/theme"
	"github.com/kaputi/navani/internal/utils/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatal(err)
	}

	logger.Log("Application started")

	theme.Init()

	p := tea.NewProgram(app.NewApp())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	err = logger.Close()

	if err != nil {
		log.Fatal(err)
	}
}
