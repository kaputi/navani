package fsutils

import (
	"fmt"
	"os"
)

func CreateDir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return nil
}
