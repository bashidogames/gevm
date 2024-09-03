package application

import (
	"fmt"
	"os"

	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/services/shortcuts/fetcher"
	"github.com/bashidogames/gevm/internal/shortcut"
	"github.com/bashidogames/gevm/internal/utils"
	"github.com/bashidogames/gevm/semver"
)

type Service struct {
	Fetcher *fetcher.Fetcher
	Config  *config.Config
}

func (s *Service) Remove(semver semver.Semver, logMissing bool) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to remove '%s' application shortcut...", semver.GodotString())
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
			utils.Printlnf("Application shortcut '%s' not found", semver.GodotString())
		}

		return nil
	}

	err = os.RemoveAll(shortcutPath)
	if err != nil {
		return fmt.Errorf("cannot remove shortcut: %w", err)
	}

	utils.Printlnf("Application '%s' shortcut removed", semver.GodotString())
	return nil
}

func (s *Service) Add(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to add '%s' application shortcut...", semver.GodotString())
	}

	targetPath, err := s.Fetcher.TargetPath(semver)
	if err != nil {
		return fmt.Errorf("cannot determine target path: %w", err)
	}

	shortcutPath := s.Fetcher.ApplicationShortcutPath(semver)
	shortcutName := s.Fetcher.ShortcutName(semver)

	if s.Config.Verbose {
		utils.Printlnf("Adding '%s' application shortcut: %s => %s", shortcutName, shortcutPath, targetPath)
	}

	exists, err := utils.DoesExist(shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		utils.Printlnf("Application shortcut '%s' already added", semver.GodotString())
		return nil
	}

	err = shortcut.Create(shortcutPath, targetPath, shortcutName)
	if err != nil {
		return fmt.Errorf("cannot create shortcut: %w", err)
	}

	utils.Printlnf("Application '%s' shortcut added", semver.GodotString())
	return nil
}

func New(fetcher *fetcher.Fetcher, config *config.Config) *Service {
	return &Service{
		Fetcher: fetcher,
		Config:  config,
	}
}
