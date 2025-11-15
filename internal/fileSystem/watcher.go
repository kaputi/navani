package filesystem

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kaputi/navani/internal/utils/logger"
)

func WatchDirectory(dirPath string) {

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
				logger.Log(fmt.Sprintf("Modified file: %s", event.Name))
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log(fmt.Sprintf("Created file: %s", event.Name))
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				logger.Log(fmt.Sprintf("Removed file: %s", event.Name))
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// in some OS the remove triggers a rename first so we handle this way
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					logger.Log(fmt.Sprintf("Removed file: %s", event.Name))
				} else {
					logger.Log(fmt.Sprintf("Renamed file: %s", event.Name))
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
