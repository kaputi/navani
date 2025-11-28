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

func Crawl(dirPath string, parentNode *TreeNode, snippetIndex *models.SnippetIndex) {
	filesInDir, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Err(err)
		return
	}

	var (
		dirMap             = make(map[fs.DirEntry]*TreeNode)
		snippetFiles       []fs.DirEntry
		allMetaFiles       = make(map[string]fs.DirEntry)
		remainingMetaFiles = make(map[string]bool)
	)

	c := config.Config()

	// Categorize files
	// files that are not c.MetaExtension or snippet files are ignored
	for _, fileEntry := range filesInDir {
		if fileEntry.IsDir() {
			dirNode := NewFileTreeNode(fileEntry.Name(), filepath.Join(dirPath, fileEntry.Name()), true)
			parentNode.AddChild(dirNode)
			dirMap[fileEntry] = dirNode
			continue
		}

		fileName := fileEntry.Name()
		extension := utils.GetExtension(fileName)

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
		metaFileName := snippetFileName + c.MetaExtension

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
		parentNode.AddChild(NewFileTreeNode(newSnippet.FileName, newSnippet.FilePath, false))
	}

	for dirFileEntry, dirNode := range dirMap {
		innerDirPath := filepath.Join(dirPath, dirFileEntry.Name())
		Crawl(innerDirPath, dirNode, snippetIndex)
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
