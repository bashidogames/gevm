package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/internal/utils"
)

type Config struct {
	ExportTemplatesRootDirectory string `json:"export-templates-root-directory,omitempty"`
	GodotRootDirectory           string `json:"godot-root-directory,omitempty"`
	ApplicationShortcutDirectory string `json:"application-shortcut-directory,omitempty"`
	DesktopShortcutDirectory     string `json:"desktop-shortcut-directory,omitempty"`
	CacheDirectory               string `json:"cache-directory,omitempty"`

	ConfigPath string            `json:"-"`
	Platform   platform.Platform `json:"-"`
	Verbose    bool              `json:"-"`
}

func (c *Config) Reset() error {
	if c.Verbose {
		utils.Printlnf("Attempting to reset config...")
	}

	config, err := DefaultConfig()
	if err != nil {
		return fmt.Errorf("could not create default config: %w", err)
	}

	*c = *config

	if c.Verbose {
		utils.Printlnf("Config reset")
	}

	return nil
}

func (c *Config) Save() error {
	if c.Verbose {
		utils.Printlnf("Attempting to save config: %s", c.ConfigPath)
	}

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

	if c.Verbose {
		utils.Printlnf("Config saved")
	}

	return nil
}

func (c *Config) load() error {
	if c.Verbose {
		utils.Printlnf("Attempting to load config: %s", c.ConfigPath)
	}

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

	if c.Verbose {
		utils.Printlnf("Config loaded")
	}

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

	defaultApplicationShortcutDirectory, err := platform.DefaultApplicationShortcutDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default application shortcut directory: %w", err)
	}

	defaultDesktopShortcutDirectory, err := platform.DefaultDesktopShortcutDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default desktop shortcut directory: %w", err)
	}

	defaultCacheDirectory, err := platform.DefaultCacheDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default cache directory: %w", err)
	}

	configPath, err := platform.ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("cannot get config path: %w", err)
	}

	platform, err := platform.Get()
	if err != nil {
		return nil, fmt.Errorf("could not determine platform: %w", err)
	}

	return &Config{
		ExportTemplatesRootDirectory: defaultExportTemplatesRootDirectory,
		GodotRootDirectory:           defaultGodotRootDirectory,
		ApplicationShortcutDirectory: defaultApplicationShortcutDirectory,
		DesktopShortcutDirectory:     defaultDesktopShortcutDirectory,
		CacheDirectory:               defaultCacheDirectory,

		ConfigPath: configPath,
		Platform:   platform,
		Verbose:    false,
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

func OptionSetVerbose(verbose bool) Option {
	return func(config *Config) {
		config.Verbose = verbose
	}
}
