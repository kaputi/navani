package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kaputi/navani/internal/app"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/utils/logger"
)

func main() {

	c := config.New()
	c.Init()

	err := logger.Init(c.DataDirPath)
	if err != nil {
		log.Fatal(err)
	}

	logger.Log("Application started")

	p := tea.NewProgram(app.NewApp(c))

	if _, err := p.Run(); err != nil {
		logger.Critical(err)
	}

	logger.Log("Application exited")

	err = logger.Close()

	if err != nil {
		log.Fatal(err)
	}
}
