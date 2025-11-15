package models

import (
	"path/filepath"
	"time"

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
	bareFileName := filepath.Base(s.FileName)
	return filepath.Join(s.DirPath, bareFileName+".meta.json")
}

type SnippetIndex struct {
	uniqueMap  map[string]*Snippet
	ByDirPath  map[string][]*Snippet
	ByFileName map[string][]*Snippet
	ByName     map[string][]*Snippet
	ByLanguage map[string][]*Snippet
	ByTag      map[string][]*Snippet
}

func NewIndex() *SnippetIndex {
	return &SnippetIndex{
		uniqueMap:  make(map[string]*Snippet),
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
	if _, exists := idx.uniqueMap[snippet.FilePath]; exists {
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

func (idx *SnippetIndex) UpdateMetadata(snippetFilePath string, metadata Metadata) {
	snippet, exists := idx.uniqueMap[snippetFilePath]
	if !exists {
		return
	}
	snippet.Metadata = metadata
}
