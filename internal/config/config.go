package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaputi/navani/internal/utils/fsutils"
	"github.com/kaputi/navani/internal/utils/logger"
)

// TODO: use filepath.Join to make this cross-platform instead of os.PathSeparator
var (
	mainDir       = ".navani"
	dataDir       = "data"
	logDir        = "logs"
	MetaExtension = ".meta.json"
	// USER DEFINED Configs
	// this is the path where the main directory will be created, this should be read from a config file or environment variable
	// userDataPath = "~" + string(os.PathSeparator) + dataDirName // TODO: make this cross-platform
	userDataPath       = "." + string(os.PathSeparator)
	TreeOpenChar       = "▼ "
	TreeCloseChar      = "▶ "
	TreeIndentChar     = "│"
	TreeDirIndentChar  = "├"
	TreeLastIndentChar = "└"
	TreeIndentSize     = 2
)

type Config struct {
	Theme             *theme
	DataPath          string
	LogsPath          string
	UserFiletypes     map[string]string
	UserFiletypeIcons map[string]string
}

func New() *Config {
	mainPath := filepath.Join(userDataPath, mainDir)
	return &Config{
		Theme:         newTheme(),
		DataPath:      filepath.Join(mainPath, dataDir),
		LogsPath:      filepath.Join(mainPath, logDir),
		UserFiletypes: make(map[string]string),
	}
}

func (c *Config) Init() {
	err := fsutils.CreateDir(c.DataPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create data directory: %w", err))
	}

	err = fsutils.CreateDir(c.LogsPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create logs directory: %w", err))
	}

	c.Theme.init()
}
