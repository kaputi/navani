package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kaputi/navani/internal/app"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/filesystem"
	"github.com/kaputi/navani/internal/models"
	"github.com/kaputi/navani/internal/utils"
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

	// until here loggs need to go to standard output because logger is not initialized
	logger.Log("Application started")

	// use config
	for ext, ft := range c.UserFiletypes {
		utils.RegisterFileType(ext, ft)
	}
	for ft, icon := range c.UserFiletypeIcons {
		utils.RegisterIcons(ft, icon)
	}

	snippetIndex := models.NewIndex()

	logger.Log(fmt.Sprintf("Crawling data directory: %s", c.DataPath))
	treeRoot := filesystem.NewFileTreeNode("root", c.DataPath, true)
	treeRoot.Open()
	// TODO: save a state for the app, for example which directories are open or closed and what snippet is selected, and try to restore if posible
	filesystem.Crawl(c.DataPath, treeRoot, snippetIndex)

	// LOG FOR DEBUGGING PURPOSES
	snippets := snippetIndex.List()
	snippetsMsg := ""
	for _, snippet := range snippets {
		snippetsMsg += fmt.Sprintf("\n %s", snippet.FilePath)
	}

	logger.Debug(fmt.Sprintf("Snippets found: %d %s", len(snippets), snippetsMsg))
	/////////////////////////////

	go filesystem.WatchDirectory(c.DataPath, snippetIndex)

	p := tea.NewProgram(app.NewApp(c, snippetIndex), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		logger.Fatal(err)
	}

	logger.Log("Application exited")
}
