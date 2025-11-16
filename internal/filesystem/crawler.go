package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/models"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

func Crawl(dirPath string, snippetIndex *models.SnippetIndex) {
	filesInDir, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Err(err)
		return
	}

	var (
		directories        []fs.DirEntry
		snippetFiles       []fs.DirEntry
		allMetaFiles       = make(map[string]fs.DirEntry)
		remainingMetaFiles = make(map[string]bool)
	)

	// Categorize files
	// files that are not c.MetaExtension or snippet files are ignored
	for _, fileEntry := range filesInDir {
		if fileEntry.IsDir() {
			directories = append(directories, fileEntry)
			continue
		}

		fileName := fileEntry.Name()
		extension := utils.GetExtension(fileName)

		switch extension {
		case config.MetaExtension:
			allMetaFiles[fileName] = fileEntry
			remainingMetaFiles[fileName] = true
		default:
			if _, err := utils.FTbyExtension(extension); err == nil {
				snippetFiles = append(snippetFiles, fileEntry)
			}
		}
	}

	for _, snippetFile := range snippetFiles {
		snippetFileName := snippetFile.Name()
		metaFileName := snippetFileName + config.MetaExtension

		metadata := models.NewMetadataFromFileName(snippetFileName)

		needsToWriteMeta := false
		if metaFile, exists := allMetaFiles[metaFileName]; exists {
			fileMetadata, err := ReadMetadata(filepath.Join(dirPath, metaFile.Name()))
			if err == nil {
				remainingMetaFiles[metaFileName] = false
				metadata = fileMetadata
			}
		} else {
			needsToWriteMeta = true
		}

		newSnippet := models.NewSnippet(dirPath, snippetFileName, metadata)

		if needsToWriteMeta {
			err := WriteMetadata(newSnippet)
			if err != nil {
				logger.Err(fmt.Errorf("failed to write metadata file: %w", err))
			}
		}

		snippetIndex.Add(newSnippet)
	}

	for _, dirEntry := range directories {
		Crawl(filepath.Join(dirPath, dirEntry.Name()), snippetIndex)
	}

	for metaFileName, isRemaining := range remainingMetaFiles {
		if isRemaining {
			err := os.Remove(filepath.Join(dirPath, metaFileName))
			if err != nil {
				logger.Err(err)
			}
		}
	}
}
