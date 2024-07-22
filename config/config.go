package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bashidogames/gdvm/internal/platform"
)

type Config struct {
	GodotRootDirectory string
	CacheDirectory     string

	Platform platform.Platform
	Verbose  bool
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

	platform, err := Platform()
	if err != nil {
		return nil, fmt.Errorf("could not determine platform: %w", err)
	}

	return &Config{
		GodotRootDirectory: defaultGodotRootDirectory,
		CacheDirectory:     defaultCacheDirectory,
		Platform:           platform,
		Verbose:            false,
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

	return config, nil
}

type Option func(*Config)

func OptionSetVerbose(verbose bool) Option {
	return func(config *Config) {
		config.Verbose = verbose
	}
}
