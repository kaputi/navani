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
		// file is a snippet
		snippet, exist := snippetIndex.ByFilePath[fullPath]
		if exist {
			snippetIndex.Remove(snippet)
		}
		metadataPath := filepath.Join(dirPath, fileName+config.MetaExtension)
		err := os.Remove(metadataPath)
		if err != nil {
			logger.Err(fmt.Errorf("failed to remove metadata file: %w", err))
		}
	} else {
		metadata := models.NewMetadataFromFileName(fileName)
		snippetPath := models.SnippetPathFromMetadataPath(fullPath)
		snippet, exists := snippetIndex.ByFilePath[snippetPath]
		// TODO: it says the snippet does not exist.. so it does not rerecreate the metadata file
		logger.Log(fmt.Sprintf("\n ====== SNIPPET PATH: %s %v", snippetPath, exists))
		if exists {
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
			dirPath := filepath.Dir(fullPath)
			fileName := filepath.Base(fullPath)
			extension := utils.GetExtension(fileName)
			_, err := utils.FTbyExtension(extension)
			if !utils.MatchExtension(fileName, config.MetaExtension) && err != nil {
				// only care about snippet files and metadata files
				continue
			}
			isMeta := extension == config.MetaExtension

			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Log(fmt.Sprintf("Modified file: %s", fullPath))

				if isMeta {
					metadata, err := ReadMetadata(fullPath)
					if err != nil {
						logger.Err(fmt.Errorf("failed to read metadata file: %w", err))
					} else {
						ft := metadata.Language
						extensionByFt, err := utils.ExtensionByFT(ft)
						if err == nil {
							snippetFilePath := filepath.Join(dirPath, fileName[:len(fileName)-len(config.MetaExtension)]+extensionByFt)
							snippetIndex.UpdateMetadata(snippetFilePath, metadata)
							logger.Log(fmt.Sprintf("Updated metadata for snippet: %s", snippetFilePath))
						}
					}
				}
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log(fmt.Sprintf("Created file: %s", fullPath))
				if !isMeta {
					metadata := models.NewMetadataFromFileName(fileName)
					snippet := models.NewSnippet(filepath.Dir(fullPath), fileName, metadata)
					snippetIndex.Add(snippet)
					logger.Log(fmt.Sprintf("Added new snippet to index: %s", fullPath))

					err := WriteMetadata(snippet)
					if err != nil {
						logger.Err(fmt.Errorf("failed to write metadata file: %w", err))
					} else {
						logger.Log(fmt.Sprintf("Created metadata file for snippet: %s", snippet.MetadataPath()))
					}
				}
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				logger.Log(fmt.Sprintf("Removed file: %s", fullPath))
				handleRemoved(fullPath, isMeta, snippetIndex)
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// in some OS the remove triggers a rename first so we handle this way
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					logger.Log(fmt.Sprintf("Removed file: %s", fullPath))
					handleRemoved(fullPath, isMeta, snippetIndex)
				} else {
					logger.Log(fmt.Sprintf("Renamed file: %s", fullPath))
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
