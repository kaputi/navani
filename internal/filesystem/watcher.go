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

func handleRemoved(fullPath string, isMeta bool, snippetIndex *models.SnippetIndex) {
	dirPath := filepath.Dir(fullPath)
	fileName := filepath.Base(fullPath)

	if !isMeta {
		snippet, exist := snippetIndex.ByFilePath[fullPath]
		if exist {
			snippetIndex.Remove(snippet)
		}
		metadataPath := filepath.Join(dirPath, fileName+config.MetaExtension)
		err := os.Remove(metadataPath)
		if err != nil {
			logger.Err(fmt.Errorf("failed to remove metadata file: %w", err))
		} else {
			logger.Log(fmt.Sprintf("Removed metadata file for snippet: %s", fullPath))
		}
	} else {
		snippetPath := models.SnippetPathFromMetadataPath(fullPath)
		snippet, exists := snippetIndex.ByFilePath[snippetPath]
		if exists {
			metadata := models.NewMetadataFromFileName(fileName)
			snippet.Metadata = metadata
			err := WriteMetadata(snippet)
			if err != nil {
				logger.Err(fmt.Errorf("failed to recreate metadata file: %w", err))
			}
		}
	}
}

func WatchDirectory(wathchPath string, snippetIndex *models.SnippetIndex) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create file system watcher: %w", err))
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			logger.Err(fmt.Errorf("failed to close file system watcher: %w", err))
		}
	}()

	err = watcher.Add(wathchPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to add directory to watcher: %w", err))
	}

	logger.Log(fmt.Sprintf("Watching directory: %s", wathchPath))

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			fullPath := event.Name
			fileName := filepath.Base(fullPath)
			extension := utils.GetExtension(fileName)

			_, err := utils.FTbyExtension(extension)
			if extension != config.MetaExtension && err != nil {
				// only care about snippet files and metadata files
				continue
			}

			isMeta := extension == config.MetaExtension

			// WRITE ==================================================
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Log(fmt.Sprintf("Modified file: %s", fullPath))

				if isMeta {
					metadata, err := ReadMetadata(fullPath)
					if err != nil {
						logger.Err(fmt.Errorf("failed to read metadata file: %w", err))
					} else {
						snippetFilePath := models.SnippetPathFromMetadataPath(fullPath)
						snippetIndex.UpdateMetadata(snippetFilePath, metadata)
					}
				}
			}
			// CREATE ==================================================
			if event.Op&fsnotify.Create == fsnotify.Create {

				logger.Log(fmt.Sprintf("Created file: %s", fullPath))
				if !isMeta {
					metadata := models.NewMetadataFromFileName(fileName)
					snippet := models.NewSnippet(filepath.Dir(fullPath), fileName, metadata)
					snippetIndex.Add(snippet)
					err := WriteMetadata(snippet)
					if err != nil {
						logger.Err(fmt.Errorf("failed to write metadata file: %w", err))
					}
				}
			}
			// REMOVE ==================================================
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				logger.Log(fmt.Sprintf("Removed file: %s", fullPath))
				handleRemoved(fullPath, isMeta, snippetIndex)
			}
			// RENAME ==================================================
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// in some OS the rename does a remove and then a create
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					logger.Log(fmt.Sprintf("Removed file from rename: %s", fullPath))
					handleRemoved(fullPath, isMeta, snippetIndex)
				} else {

					logger.Log(fmt.Sprintf("Renamed file: %s", fullPath))
					// TODO: if file is snippet, update and rename metadata
					// if file is metadata, keep old name (WATCH RECURSIVE RENAMES!!!!!)
					// event.renameName is not exported by fsnotify because of a bug,
					// at the moment there is no handling for rename,if this is not fixed soon
					// i should track the rename and implement the handler myself.
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
