package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/models"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

func WatchDirectory(dirPath string, snippetIndex *models.SnippetIndex, c *config.Config) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create file system watcher: %w", err))
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			logger.Err(fmt.Errorf("failed to close file system watcher: %w", err))
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to add directory to watcher: %w", err))
	}

	logger.Log(fmt.Sprintf("Watching directory: %s", dirPath))

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// logger.Log(fmt.Sprintf("Modified file: %s", event.Name))
				// TODO: if file is snippet, there is not much to do because metadata is separate
				// if file is metadata, update index accordingly

				fullPath := event.Name
				fileName := filepath.Base(fullPath)
				if utils.MatchExtension(fileName, c.MetaExtension) {
					metadata, err := ReadMetadata(fullPath)
					if err != nil {
						logger.Err(fmt.Errorf("failed to read metadata file: %w", err))
					} else {
						ft := metadata.Language
						extensionByFt, err := utils.ExtensionByFT(ft)
						if err == nil {
							dirPath := filepath.Dir(fullPath)
							snippetFilePath := filepath.Join(dirPath, fileName[:len(fileName)-len(c.MetaExtension)]+extensionByFt)
							snippetIndex.UpdateMetadata(snippetFilePath, metadata)
						}
					}
				}
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log(fmt.Sprintf("Created file: %s", event.Name))
				// TODO: if file is snippet, create metadata file too
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				logger.Log(fmt.Sprintf("Removed file: %s", event.Name))
				// TODO: if file is snippet, remove metadata file too
				// if file is metadata, recreate default metadata, only snippets may be removed
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// in some OS the remove triggers a rename first so we handle this way
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					logger.Log(fmt.Sprintf("Removed file: %s", event.Name))
					// TODO: if file is snippet, remove metadata file too
					// if file is metadata, recreate default metadata, only snippets may be removed
				} else {
					logger.Log(fmt.Sprintf("Renamed file: %s", event.Name))
					// TODO: if file is snippet, update and rename metadata
					// if file is metadata, keep old name (WATCH RECURSIVE RENAMES!!!!!)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Err(fmt.Errorf("file system watcher error: %w", err))
		}
	}
}
