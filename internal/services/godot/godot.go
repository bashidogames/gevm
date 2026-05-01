package godot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bashmills/gevm/config"
	"github.com/bashmills/gevm/internal/archiving"
	"github.com/bashmills/gevm/internal/downloading"
	"github.com/bashmills/gevm/internal/environment"
	"github.com/bashmills/gevm/internal/utils"
	"github.com/bashmills/gevm/semver"
	"github.com/jedib0t/go-pretty/v6/table"
)

const CACHE_FOLDER = "godot"

type ExportTemplatesChecker interface {
	Exists(semver semver.Semver) (bool, error)
}

type ExecutableLocator interface {
	Find(semver semver.Semver) (string, error)
}

type Service struct {
	Environment            *environment.Environment
	ExportTemplatesChecker ExportTemplatesChecker
	ExecutableLocator      ExecutableLocator
	Config                 *config.Config
}

func (s *Service) Download(semver semver.Semver) error {
	s.Config.Logger.Debugf("Attempting to download '%s' godot...", semver.GodotString())

	asset, err := s.Environment.FetchGodotAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		s.Config.Logger.Errorf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("fetch asset failed: %w", err)
	}

	archivePath := s.archivePath(asset.Name)

	exists, err := utils.DoesExist(archivePath)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		s.Config.Logger.Infof("Godot '%s' already downloaded", semver.GodotString())
		return nil
	}

	s.Config.Logger.Debugf("Downloading from: %s", asset.DownloadURL)
	s.Config.Logger.Debugf("Downloading to: %s", archivePath)

	err = downloading.Download(s.Config.Logger, asset.DownloadURL, archivePath, s.Config.Silent)
	if errors.Is(err, downloading.ErrNotFound) {
		s.Config.Logger.Errorf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	s.Config.Logger.Infof("Godot '%s' downloaded", semver.GodotString())
	return nil
}

func (s *Service) Uninstall(semver semver.Semver, logMissing bool) error {
	s.Config.Logger.Debugf("Attempting to uninstall '%s' godot...", semver.GodotString())

	targetDirectory := s.targetDirectory(semver)

	exists, err := utils.DoesExist(targetDirectory)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if !exists {
		if logMissing {
			s.Config.Logger.Errorf("Godot '%s' not found", semver.GodotString())
		}

		return nil
	}

	s.Config.Logger.Debugf("Removing directory: %s", targetDirectory)

	err = os.RemoveAll(targetDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove target directory: %w", err)
	}

	s.Config.Logger.Infof("Godot '%s' uninstalled", semver.GodotString())
	return nil
}

func (s *Service) Install(semver semver.Semver) error {
	s.Config.Logger.Debugf("Attempting to install '%s' godot...", semver.GodotString())

	asset, err := s.Environment.FetchGodotAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		s.Config.Logger.Errorf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("fetch asset failed: %w", err)
	}

	targetDirectory := s.targetDirectory(semver)
	archivePath := s.archivePath(asset.Name)

	exists, err := utils.DoesExist(targetDirectory)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		s.Config.Logger.Infof("Godot '%s' already installed", semver.GodotString())
		return nil
	}

	err = os.MkdirAll(s.Config.GodotRootDirectory, utils.OS_DIRECTORY)
	if err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	err = os.RemoveAll(targetDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove target directory: %w", err)
	}

	s.Config.Logger.Debugf("Downloading from: %s", asset.DownloadURL)
	s.Config.Logger.Debugf("Downloading to: %s", archivePath)

	err = downloading.Download(s.Config.Logger, asset.DownloadURL, archivePath, s.Config.Silent)
	if errors.Is(err, downloading.ErrNotFound) {
		s.Config.Logger.Errorf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	s.Config.Logger.Debugf("Unzipping from: %s", archivePath)
	s.Config.Logger.Debugf("Unzipping to: %s", targetDirectory)

	err = archiving.Unzip(s.Config.Logger, archivePath, targetDirectory)
	if err != nil {
		return fmt.Errorf("unzip failed: %w", err)
	}

	s.Config.Logger.Infof("Godot '%s' installed", semver.GodotString())
	return nil
}

func (s *Service) Path(semver semver.Semver) error {
	targetPath, err := s.ExecutableLocator.Find(semver)
	if errors.Is(err, os.ErrNotExist) {
		s.Config.Logger.Errorf("Godot '%s' not found. Use `gevm godot list` to see installed versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot determine target path: %w", err)
	}

	utils.Println(targetPath)
	return nil
}

func (s *Service) List() error {
	entries, err := os.ReadDir(s.Config.GodotRootDirectory)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return fmt.Errorf("cannot read godot root directory: %w", err)
	}

	if len(entries) == 0 {
		s.Config.Logger.Info("No godot engine versions installed")
		return nil
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Version", "Release", "Export Templates?", "Mono?"})

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		semver, err := semver.Parse(entry.Name())
		if err != nil {
			s.Config.Logger.Warningf("Failed to recognize version: %s", err)
			continue
		}

		exportTemplates, err := s.ExportTemplatesChecker.Exists(semver)
		if err != nil {
			s.Config.Logger.Warningf("Failed to check export templates existence: %s", err)
		}

		version := semver.Relver.Version.String()
		release := semver.Relver.Release.String()
		mono := semver.Mono

		t.AppendRow(table.Row{version, release, exportTemplates, mono})
	}

	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}

func (s *Service) Clear() error {
	entries, err := os.ReadDir(s.Config.GodotRootDirectory)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return fmt.Errorf("cannot read godot root directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		semver, err := semver.Parse(entry.Name())
		if err != nil {
			s.Config.Logger.Warningf("Failed to recognize version: %s", err)
			continue
		}

		err = s.Uninstall(semver, true)
		if err != nil {
			return fmt.Errorf("cannot clear godot: %w", err)
		}
	}

	empty, err := utils.IsDirectoryEmpty(s.Config.GodotRootDirectory)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return fmt.Errorf("failed to check emptiness: %w", err)
	}

	if empty {
		err = os.Remove(s.Config.GodotRootDirectory)
		if err != nil {
			return fmt.Errorf("cannot remove godot root directory: %w", err)
		}
	}

	s.Config.Logger.Info("Godot cleared")
	return nil
}

func (s *Service) targetDirectory(semver semver.Semver) string {
	return filepath.Join(s.Config.GodotRootDirectory, semver.GodotString())
}

func (s *Service) archivePath(name string) string {
	return filepath.Join(s.Config.CacheDirectory, CACHE_FOLDER, name)
}

func New(environment *environment.Environment, exportTemplatesChecker ExportTemplatesChecker, executableLocator ExecutableLocator, config *config.Config) *Service {
	return &Service{
		Environment:            environment,
		ExportTemplatesChecker: exportTemplatesChecker,
		ExecutableLocator:      executableLocator,
		Config:                 config,
	}
}
