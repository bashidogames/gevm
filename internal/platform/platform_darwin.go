package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

const CONFIG_FILENAME = "config.json"

func DefaultExportTemplatesRootDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	directory := filepath.Join(userHomeDir, "Library", "Application Support", "Godot", "templates")
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

func DefaultApplicationShortcutDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	directory := filepath.Join(userHomeDir, "Applications")
	return directory, nil
}

func DefaultDesktopShortcutDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	directory := filepath.Join(userHomeDir, "Desktop")
	return directory, nil
}

func DefaultCacheDirectory() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user cache directory: %w", err)
	}

	directory := filepath.Join(userCacheDir, "bashidogames", "gevm")
	return directory, nil
}

func DefaultBinDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	directory := filepath.Join(userHomeDir, "bin")
	return directory, nil
}

func ConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user config directory: %w", err)
	}

	directory := filepath.Join(userConfigDir, "bashidogames", "gevm")
	path := filepath.Join(directory, CONFIG_FILENAME)
	return path, nil
}
