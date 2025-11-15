package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kaputi/navani/internal/app"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/filesystem"
	"github.com/kaputi/navani/internal/models"
	"github.com/kaputi/navani/internal/utils/logger"
)

func main() {

	c := config.New()
	c.Init()

	err := logger.Init(c.LogsPath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := logger.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// NOTE: until here loggs need to go to standard output because logger is not initialized

	logger.Log("Application started")

	snippetIndex := models.NewIndex()

	filesystem.Crawl(c.DataPath, snippetIndex)

	go filesystem.WatchDirectory(c.DataPath, snippetIndex)

	p := tea.NewProgram(app.NewApp(c), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		logger.Fatal(err)
	}

	logger.Log("Application exited")
}
