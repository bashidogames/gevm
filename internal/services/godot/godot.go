package godot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/archiving"
	"github.com/bashidogames/gevm/internal/downloading"
	"github.com/bashidogames/gevm/internal/environment"
	"github.com/bashidogames/gevm/internal/services/godot/fetcher"
	"github.com/bashidogames/gevm/internal/utils"
	"github.com/bashidogames/gevm/semver"
	"github.com/jedib0t/go-pretty/v6/table"
)

const CACHE_FOLDER = "godot"

type Service struct {
	Environment *environment.Environment
	Fetcher     *fetcher.Fetcher
	Config      *config.Config
}

func (s *Service) Download(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to download '%s' godot...", semver.GodotString())
	}

	asset, err := s.Environment.FetchGodotAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
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
		utils.Printlnf("Godot '%s' already downloaded", semver.GodotString())
		return nil
	}

	if s.Config.Verbose {
		utils.Printlnf("Downloading from: %s", asset.DownloadURL)
		utils.Printlnf("Downloading to: %s", archivePath)
	}

	err = downloading.Download(asset.DownloadURL, archivePath)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	utils.Printlnf("Godot '%s' downloaded", semver.GodotString())
	return nil
}

func (s *Service) Uninstall(semver semver.Semver, logMissing bool) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to uninstall '%s' godot...", semver.GodotString())
	}

	targetDirectory := s.targetDirectory(semver)

	exists, err := utils.DoesExist(targetDirectory)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if !exists {
		if logMissing {
			utils.Printlnf("Godot '%s' not found", semver.GodotString())
		}

		return nil
	}

	if s.Config.Verbose {
		utils.Printlnf("Removing directory: %s", targetDirectory)
	}

	err = os.RemoveAll(targetDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove target directory: %w", err)
	}

	utils.Printlnf("Godot '%s' uninstalled", semver.GodotString())
	return nil
}

func (s *Service) Install(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to install '%s' godot...", semver.GodotString())
	}

	asset, err := s.Environment.FetchGodotAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
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
		utils.Printlnf("Godot '%s' already installed", semver.GodotString())
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

	if s.Config.Verbose {
		utils.Printlnf("Downloading from: %s", asset.DownloadURL)
		utils.Printlnf("Downloading to: %s", archivePath)
	}

	err = downloading.Download(asset.DownloadURL, archivePath)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Godot '%s' not found. Use 'gevm versions list' to see available versions.", semver.GodotString())
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if s.Config.Verbose {
		utils.Printlnf("Unzipping from: %s", archivePath)
		utils.Printlnf("Unzipping to: %s", targetDirectory)
	}

	err = archiving.Unzip(archivePath, targetDirectory)
	if err != nil {
		return fmt.Errorf("unzip failed: %w", err)
	}

	utils.Printlnf("Godot '%s' installed", semver.GodotString())
	return nil
}

func (s *Service) Use(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to use '%s' godot...", semver.GodotString())
	}

	targetPath, err := s.Fetcher.TargetPath(semver)
	if err != nil {
		return fmt.Errorf("cannot determine target path: %w", err)
	}

	linkPath := s.Fetcher.LinkPath(semver)

	if s.Config.Verbose {
		utils.Printlnf("Creating godot symlink: %s => %s", linkPath, targetPath)
	}

	err = os.MkdirAll(s.Config.BinDirectory, utils.OS_DIRECTORY)
	if err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	err = os.RemoveAll(linkPath)
	if err != nil {
		return fmt.Errorf("cannot remove link path: %w", err)
	}

	err = os.Symlink(targetPath, linkPath)
	if err != nil {
		return fmt.Errorf("cannot create link: %w", err)
	}

	utils.Printlnf("Using '%s' godot", semver.GodotString())
	return nil
}

func (s *Service) List() error {
	entries, err := os.ReadDir(s.Config.GodotRootDirectory)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return fmt.Errorf("cannot read godot root directory: %w", err)
	}

	if len(entries) == 0 {
		utils.Printlnf("No godot engine versions installed")
		return nil
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Version", "Release", "Mono?"})

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		semver, err := semver.Parse(entry.Name())
		if err != nil {
			if s.Config.Verbose {
				utils.Printlnf("Failed to recognize version: %s", err)
			}

			continue
		}

		version := semver.Relver.Version.String()
		release := semver.Relver.Release.String()
		mono := semver.Mono

		t.AppendRow(table.Row{version, release, mono})
	}

	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}

func (s *Service) targetDirectory(semver semver.Semver) string {
	return filepath.Join(s.Config.GodotRootDirectory, semver.GodotString())
}

func (s *Service) archivePath(name string) string {
	return filepath.Join(s.Config.CacheDirectory, CACHE_FOLDER, name)
}

func New(environment *environment.Environment, config *config.Config) *Service {
	fetcher := fetcher.New(config)
	return &Service{
		Environment: environment,
		Fetcher:     fetcher,
		Config:      config,
	}
}
