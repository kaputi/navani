package filesystem

import (
	"encoding/json"
	"os"

	"github.com/kaputi/navani/internal/models"
)

func CreateSnippet(snippet *models.Snippet, content string) error {
	err := WriteSnippet(snippet, content)
	if err != nil {
		return err
	}
	err = WriteMetadata(snippet)
	if err != nil {
		return err
	}
	return nil
}

func WriteSnippet(snippet *models.Snippet, content string) error {
	err := os.WriteFile(snippet.FilePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func WriteMetadata(snippet *models.Snippet) error {
	metadata := snippet.Metadata
	metadataFilePath := snippet.MetadataPath()

	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(metadataFilePath, metadataBytes, 0644)
	if err != nil {
		return err
	}

	return nil

}

func ReadSnippetContent(snippet *models.Snippet) (string, error) {
	content, err := os.ReadFile(snippet.FilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ReadMetadata(metadataFilePath string) (models.Metadata, error) {
	var metadata models.Metadata

	bytes, err := os.ReadFile(metadataFilePath)

	if err != nil {
		return metadata, err
	}

	err = json.Unmarshal(bytes, &metadata)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}
