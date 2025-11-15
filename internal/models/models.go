package models

import (
	"path/filepath"
	"slices"
	"time"

	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/utils"
)

type Metadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
	Copies      int      `json:"copies"`
	Tags        []string `json:"tags"`
}

func NewMetadataFromFileName(fileName string) Metadata {
	name := filepath.Base(fileName)
	ft, _ := utils.FTbyFileName(fileName)
	language := ft

	return Metadata{
		Name:        name,
		Description: "",
		Language:    language,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		Copies:      0,
		Tags:        []string{},
	}
}

type Snippet struct {
	Metadata
	FilePath string
	DirPath  string
	FileName string
}

func NewSnippet(dirPath, fileName string, metadata Metadata) *Snippet {
	return &Snippet{
		Metadata: metadata,
		FilePath: filepath.Join(dirPath, fileName),
		DirPath:  dirPath,
		FileName: fileName,
	}
}

func (s *Snippet) MetadataPath() string {
	return filepath.Join(s.DirPath, s.FileName+config.MetaExtension)
}

type SnippetIndex struct {
	ByFilePath map[string]*Snippet
	ByDirPath  map[string][]*Snippet
	ByFileName map[string][]*Snippet
	ByName     map[string][]*Snippet
	ByLanguage map[string][]*Snippet
	ByTag      map[string][]*Snippet
}

func NewIndex() *SnippetIndex {
	return &SnippetIndex{
		ByFilePath: make(map[string]*Snippet),
		ByDirPath:  make(map[string][]*Snippet),
		ByFileName: make(map[string][]*Snippet),
		ByName:     make(map[string][]*Snippet),
		ByLanguage: make(map[string][]*Snippet),
		ByTag:      make(map[string][]*Snippet),
	}
}

func addToMapList(m map[string][]*Snippet, key string, snippet *Snippet) {
	if _, exists := m[key]; !exists {
		m[key] = []*Snippet{snippet}
	} else {
		m[key] = append(m[key], snippet)
	}
}

func (idx *SnippetIndex) Add(snippet *Snippet) {
	// do not add duplicates
	if _, exists := idx.ByFilePath[snippet.FilePath]; exists {
		return
	}

	addToMapList(idx.ByDirPath, snippet.DirPath, snippet)
	addToMapList(idx.ByFileName, snippet.FileName, snippet)
	addToMapList(idx.ByName, snippet.Name, snippet)
	addToMapList(idx.ByLanguage, snippet.Language, snippet)
	for _, tag := range snippet.Tags {
		addToMapList(idx.ByTag, tag, snippet)
	}

}

func removeFromSlice(slice []*Snippet, snippet *Snippet) []*Snippet {
	for i, s := range slice {
		if s == snippet {
			return slices.Delete(slice, i, i+1)
		}
	}

	return slice
}

func (idx *SnippetIndex) Remove(snippet *Snippet) {
	_, exists := idx.ByFilePath[snippet.FilePath]
	if !exists {
		return
	}

	delete(idx.ByFilePath, snippet.FilePath)
	idx.ByDirPath[snippet.DirPath] = removeFromSlice(idx.ByDirPath[snippet.DirPath], snippet)
	idx.ByFileName[snippet.FileName] = removeFromSlice(idx.ByFileName[snippet.FileName], snippet)
	idx.ByName[snippet.Name] = removeFromSlice(idx.ByName[snippet.Name], snippet)
	idx.ByLanguage[snippet.Language] = removeFromSlice(idx.ByLanguage[snippet.Language], snippet)
	for _, tag := range snippet.Tags {
		idx.ByTag[tag] = removeFromSlice(idx.ByTag[tag], snippet)
	}
}

func (idx *SnippetIndex) UpdateMetadata(snippetFilePath string, metadata Metadata) {
	snippet, exists := idx.ByFilePath[snippetFilePath]
	if !exists {
		return
	}
	snippet.Metadata = metadata
}

func SnippetPathFromMetadataPath(metadataPath string) string {
	return metadataPath[:len(metadataPath)-len(config.MetaExtension)]
}
