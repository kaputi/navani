package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/models"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

func Crawl(dirPath string, snippetIndex *models.SnippetIndex, c *config.Config) {
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
	// files that are not .json or snippet files are ignored
	// all .json files are considered metadata files // TODO: this may need to be more strict
	for _, fileEntry := range filesInDir {
		if fileEntry.IsDir() {
			directories = append(directories, fileEntry)
			continue
		}

		fileName := fileEntry.Name()
		extension := filepath.Ext(fileName)
		if utils.MatchExtension(fileName, c.MetaExtension) {
			extension = c.MetaExtension
		}

		switch extension {
		case c.MetaExtension:
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
		extension := filepath.Ext(snippetFileName)
		bareName := snippetFileName[:len(snippetFileName)-len(extension)]
		metaFileName := bareName + c.MetaExtension

		metadata := models.NewMetadataFromFileName(snippetFileName)

		if metaFile, exists := allMetaFiles[metaFileName]; exists {
			remainingMetaFiles[metaFileName] = false
			fileMetadata, err := ReadMetadata(filepath.Join(dirPath, metaFile.Name()))
			if err == nil {
				metadata = fileMetadata
			}
		}

		newSnippet := models.NewSnippet(dirPath, snippetFileName, metadata)

		snippetIndex.Add(newSnippet)
	}

	for _, dirEntry := range directories {
		Crawl(filepath.Join(dirPath, dirEntry.Name()), snippetIndex, c)
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
