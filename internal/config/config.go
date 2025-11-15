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
	dataDirName = ".navani"
	// defaultDataDirPath = "~" + string(os.PathSeparator) + dataDirName // TODO: make this cross-platform
	defaultDataDirPath = "." + string(os.PathSeparator)
)

type Config struct {
	Theme       *theme
	DataDirPath string
}

func New() *Config {
	// TODO: read config file and set these values accordingly
	return &Config{
		Theme:       newTheme(),
		DataDirPath: filepath.Join(defaultDataDirPath, dataDirName),
	}
}

func (c *Config) Init() {
	err := fsutils.CreateDir(c.DataDirPath)
	if err != nil {
		logger.Critical(fmt.Errorf("failed to create data directory: %w", err))
	}

	c.Theme.init()
}
