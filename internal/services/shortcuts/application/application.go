package application

import (
	"fmt"
	"os"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/services/shortcuts/fetcher"
	"github.com/bashidogames/gdvm/internal/shortcut"
	"github.com/bashidogames/gdvm/internal/utils"
	"github.com/bashidogames/gdvm/semver"
)

type Service struct {
	Fetcher *fetcher.Fetcher
	Config  *config.Config
}

func (s *Service) Remove(semver semver.Semver, logMissing bool) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to remove '%s' application shortcut...", semver)
	}

	shortcutPath := s.Fetcher.ApplicationShortcutPath(semver)

	if s.Config.Verbose {
		utils.Printlnf("Removing application shortcut: %s", shortcutPath)
	}

	exists, err := utils.DoesExist(shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if !exists {
		if logMissing {
			utils.Printlnf("Application shortcut '%s' not found", semver)
		}

		return nil
	}

	err = os.RemoveAll(shortcutPath)
	if err != nil {
		return fmt.Errorf("cannot remove shortcut: %w", err)
	}

	utils.Printlnf("Application '%s' shortcut removed", semver)
	return nil
}

func (s *Service) Add(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to add '%s' application shortcut...", semver)
	}

	targetPath, err := s.Fetcher.TargetPath(semver)
	if err != nil {
		return fmt.Errorf("cannot determine target path: %w", err)
	}

	shortcutPath := s.Fetcher.ApplicationShortcutPath(semver)
	shortcutName := s.Fetcher.ShortcutName(semver)

	if s.Config.Verbose {
		utils.Printlnf("Adding application shortcut: %s", shortcutPath)
		utils.Printlnf("Target path: %s", targetPath)
	}

	exists, err := utils.DoesExist(shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		utils.Printlnf("Application shortcut '%s' already added", semver)
		return nil
	}

	err = shortcut.Create(shortcutPath, targetPath, shortcutName)
	if err != nil {
		return fmt.Errorf("cannot create shortcut: %w", err)
	}

	utils.Printlnf("Application '%s' shortcut added", semver)
	return nil
}

func New(fetcher *fetcher.Fetcher, config *config.Config) *Service {
	return &Service{
		Fetcher: fetcher,
		Config:  config,
	}
}
