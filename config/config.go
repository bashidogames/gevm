package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/internal/utils"
)

const CONFIG_FILENAME = "config.json"

type Config struct {
	GodotRootDirectory string `json:"godot-root-directory,omitempty"`
	CacheDirectory     string `json:"cache-directory,omitempty"`

	ConfigPath string            `json:"-"`
	Platform   platform.Platform `json:"-"`
	Verbose    bool              `json:"-"`
}

func (c *Config) Reset() error {
	config, err := DefaultConfig()
	if err != nil {
		return fmt.Errorf("could not create default config: %w", err)
	}

	*c = *config
	return nil
}

func (c *Config) Save() error {
	if c.Verbose {
		utils.Printlnf("Attempting to save config: %s", c.ConfigPath)
	}

	err := os.MkdirAll(filepath.Dir(c.ConfigPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	bytes, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return fmt.Errorf("cannot parse config: %w", err)
	}

	err = os.WriteFile(c.ConfigPath, bytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write config: %w", err)
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

	return nil
}

func DefaultGodotRootDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	directory := filepath.Join(userHomeDir, "Godot")
	return directory, nil
}

func DefaultCacheDirectory() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user cache directory: %w", err)
	}

	directory := filepath.Join(userCacheDir, "bashidogames", "gdvm")
	return directory, nil
}

func ConfigPath() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user config directory: %w", err)
	}

	directory := filepath.Join(root, "bashidogames", "gdvm")
	path := filepath.Join(directory, CONFIG_FILENAME)
	return path, nil
}

func Platform() (platform.Platform, error) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return platform.WindowsAmd64, nil
		case "386":
			return platform.Windows386, nil
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			return platform.LinuxArm64, nil
		case "amd64":
			return platform.LinuxAmd64, nil
		case "arm":
			return platform.LinuxArm, nil
		case "386":
			return platform.Linux386, nil
		}
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return platform.DarwinArm64, nil
		case "amd64":
			return platform.DarwinAmd64, nil
		case "386":
			return platform.Darwin386, nil
		}
	}

	return "", fmt.Errorf("invalid platform")
}

func DefaultConfig() (*Config, error) {
	defaultGodotRootDirectory, err := DefaultGodotRootDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default godot root directory: %w", err)
	}

	defaultCacheDirectory, err := DefaultCacheDirectory()
	if err != nil {
		return nil, fmt.Errorf("cannot get default cache directory: %w", err)
	}

	configPath, err := ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("cannot get config path: %w", err)
	}

	platform, err := Platform()
	if err != nil {
		return nil, fmt.Errorf("could not determine platform: %w", err)
	}

	return &Config{
		GodotRootDirectory: defaultGodotRootDirectory,
		CacheDirectory:     defaultCacheDirectory,

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
		config.ConfigPath = configPath
	}
}

func OptionSetVerbose(verbose bool) Option {
	return func(config *Config) {
		config.Verbose = verbose
	}
}
