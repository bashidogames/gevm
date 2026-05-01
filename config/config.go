package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bashmills/gevm/internal/logging"
	"github.com/bashmills/gevm/internal/platform"
	"github.com/bashmills/gevm/internal/utils"
	"github.com/bashmills/gevm/logger"
)

type Config struct {
	ExportTemplatesRootDirectory string `json:"export-templates-root-directory"`
	GodotRootDirectory           string `json:"godot-root-directory"`
	CacheDirectory               string `json:"cache-directory"`
	BinDirectory                 string `json:"bin-directory"`

	ConfigPath string            `json:"-"`
	Platform   platform.Platform `json:"-"`
	Logger     logger.Logger     `json:"-"`
	Silent     bool              `json:"-"`
}

func (c *Config) Reset() error {
	c.Logger.Trace("Attempting to reset config...")

	config, err := DefaultConfig()
	if err != nil {
		return fmt.Errorf("could not create default config: %w", err)
	}

	*c = *config

	c.Logger.Trace("Config reset")

	return nil
}

func (c *Config) Save() error {
	c.Logger.Tracef("Attempting to save config: %s", c.ConfigPath)

	err := os.MkdirAll(filepath.Dir(c.ConfigPath), utils.OS_DIRECTORY)
	if err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	bytes, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return fmt.Errorf("cannot parse config: %w", err)
	}

	file, err := os.OpenFile(c.ConfigPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, utils.OS_FILE)
	if err != nil {
		return fmt.Errorf("cannot create config: %w", err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}

	c.Logger.Trace("Config saved")

	return nil
}

func (c *Config) load() error {
	c.Logger.Tracef("Attempting to load config: %s", c.ConfigPath)

	file, err := os.Open(c.ConfigPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot open config: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return fmt.Errorf("cannot parse config: %w", err)
	}

	c.Logger.Trace("Config loaded")

	return nil
}

func DefaultConfig() (*Config, error) {
	defaultExportTemplatesRootDirectory, err := platform.DefaultExportTemplatesRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default export templates root directory: %w", err)
	}

	defaultGodotRootDirectory, err := platform.DefaultGodotRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default godot root directory: %w", err)
	}

	defaultCacheDirectory, err := platform.DefaultCacheDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default cache directory: %w", err)
	}

	defaultBinDirectory, err := platform.DefaultBinDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default bin directory: %w", err)
	}

	configPath, err := platform.ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("cannot get config path: %w", err)
	}

	platform, err := platform.Get()
	if err != nil {
		return nil, fmt.Errorf("could not determine platform: %w", err)
	}

	logger, err := logging.New(logging.DEBUG)
	if err != nil {
		return nil, fmt.Errorf("could not create default logger: %w", err)
	}

	return &Config{
		ExportTemplatesRootDirectory: defaultExportTemplatesRootDirectory,
		GodotRootDirectory:           defaultGodotRootDirectory,
		CacheDirectory:               defaultCacheDirectory,
		BinDirectory:                 defaultBinDirectory,

		ConfigPath: configPath,
		Platform:   platform,
		Logger:     logger,
		Silent:     false,
	}, nil
}

func New(options ...Option) (*Config, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("could not create default config: %w", err)
	}

	for _, option := range options {
		option(config)
	}

	err = config.load()
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	err = config.Save()
	if err != nil {
		return nil, fmt.Errorf("cannot save config: %w", err)
	}

	return config, nil
}

type Option func(*Config)

func OptionSetConfigPath(configPath string) Option {
	return func(config *Config) {
		if len(configPath) > 0 {
			config.ConfigPath = configPath
		}
	}
}

func OptionSetLogger(logger logger.Logger) Option {
	return func(config *Config) {
		if logger != nil {
			config.Logger = logger
		}
	}
}

func OptionSetSilent(silent bool) Option {
	return func(config *Config) {
		config.Silent = silent
	}
}
