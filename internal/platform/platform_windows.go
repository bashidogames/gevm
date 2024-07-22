package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

const CONFIG_FILENAME = "config.json"

func DefaultBuildTemplatesRootDirectory() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user config directory: %w", err)
	}

	directory := filepath.Join(root, "Godot", "export_templates")
	return directory, nil
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
