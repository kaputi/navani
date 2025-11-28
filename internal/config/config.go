package config

import (
	"fmt"
	"path/filepath"

	"github.com/kaputi/navani/internal/utils/fsutils"
)

var (
	localConfig config
)

type config struct {
	DataPath          string
	LogsPath          string
	UserFiletypes     map[string]string
	UserFiletypeIcons map[string]string
	MetaExtension     string
}

func Config() config {
	return localConfig
}

func LoadConfig() error {
	mainDir := ".navani"
	dataDir := "data"
	logDir := "logs"

	// TODO: reads the config from the config directory (not configurable, depends on os)
	// userDataPath = "~" + string(os.PathSeparator) + dataDirName // TODO: make this cross-platform
	mainPath := filepath.Join(".", mainDir)

	localConfig = config{
		DataPath:          filepath.Join(mainPath, dataDir),
		LogsPath:          filepath.Join(mainPath, logDir),
		UserFiletypes:     make(map[string]string),
		UserFiletypeIcons: make(map[string]string),
		MetaExtension:     ".meta.json",
	}

	err := fsutils.CreateDir(localConfig.DataPath)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	err = fsutils.CreateDir(localConfig.LogsPath)
	if err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	return nil
}
