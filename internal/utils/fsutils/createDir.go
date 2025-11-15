package fsutils

import (
	"fmt"
	"os"
)

func CreateDir(dirPath string) error {
	if !PathExists(dirPath) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}
	return nil
}
