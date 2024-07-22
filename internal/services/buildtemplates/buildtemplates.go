package buildtemplates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/archiving"
	"github.com/bashidogames/gdvm/internal/downloading"
	"github.com/bashidogames/gdvm/internal/environment"
	"github.com/bashidogames/gdvm/internal/utils"
	"github.com/bashidogames/gdvm/semver"
	"github.com/jedib0t/go-pretty/v6/table"
)

const CACHE_FOLDER = "build-templates"
const TEMP_FOLDER = "templates"

type Service struct {
	Environment *environment.Environment
	Config      *config.Config
}

func (s *Service) Download(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to download '%s' build templates...", semver)
	}

	asset, err := s.Environment.FetchBuildTemplatesAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Build templates '%s' not found. Use 'gdvm versions list' to see available versions.", semver)
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
		utils.Printlnf("Build templates '%s' already downloaded", semver)
		return nil
	}

	if s.Config.Verbose {
		utils.Printlnf("Downloading from: %s", asset.DownloadURL)
		utils.Printlnf("Downloading to: %s", archivePath)
	}

	err = downloading.Download(asset.DownloadURL, archivePath)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Build templates '%s' not found. Use 'gdvm versions list' to see available versions.", semver)
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	utils.Printlnf("Build templates '%s' downloaded", semver)
	return nil
}

func (s *Service) Uninstall(semver semver.Semver, logMissing bool) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to uninstall '%s' build templates...", semver)
	}

	targetDirectory := s.targetDirectory(semver)

	exists, err := utils.DoesExist(targetDirectory)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if !exists {
		if logMissing {
			utils.Printlnf("Build templates '%s' not found", semver)
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

	utils.Printlnf("Build templates '%s' uninstalled", semver)
	return nil
}

func (s *Service) Install(semver semver.Semver) error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to install '%s' build templates...", semver)
	}

	asset, err := s.Environment.FetchBuildTemplatesAsset(semver)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Build templates '%s' not found. Use 'gdvm versions list' to see available versions.", semver)
		return nil
	}
	if err != nil {
		return fmt.Errorf("fetch asset failed: %w", err)
	}

	targetDirectory := s.targetDirectory(semver)
	archivePath := s.archivePath(asset.Name)
	tempDirectory := s.tempDirectory()
	rootDirectory := s.rootDirectory()

	exists, err := utils.DoesExist(targetDirectory)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}

	if exists {
		utils.Printlnf("Build templates '%s' already installed", semver)
		return nil
	}

	err = os.MkdirAll(s.Config.BuildTemplatesRootDirectory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	err = os.RemoveAll(targetDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove target directory: %w", err)
	}

	err = os.RemoveAll(tempDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove temp directory: %w", err)
	}

	if s.Config.Verbose {
		utils.Printlnf("Downloading from: %s", asset.DownloadURL)
		utils.Printlnf("Downloading to: %s", archivePath)
	}

	err = downloading.Download(asset.DownloadURL, archivePath)
	if errors.Is(err, downloading.ErrNotFound) {
		utils.Printlnf("Build templates '%s' not found. Use 'gdvm versions list' to see available versions.", semver)
		return nil
	}
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if s.Config.Verbose {
		utils.Printlnf("Unzipping from: %s", archivePath)
		utils.Printlnf("Unzipping to: %s", rootDirectory)
	}

	err = archiving.Unzip(archivePath, rootDirectory)
	if err != nil {
		return fmt.Errorf("unzip failed: %w", err)
	}

	if s.Config.Verbose {
		utils.Printlnf("Moving from: %s", tempDirectory)
		utils.Printlnf("Moving to: %s", targetDirectory)
	}

	err = os.Rename(tempDirectory, targetDirectory)
	if err != nil {
		return fmt.Errorf("move failed: %w", err)
	}

	utils.Printlnf("Build templates '%s' installed", semver)
	return nil
}

func (s *Service) List() error {
	entries, err := os.ReadDir(s.Config.BuildTemplatesRootDirectory)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return fmt.Errorf("cannot read build templates root directory: %w", err)
	}

	if len(entries) == 0 {
		utils.Printlnf("No build templates installed")
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
	return filepath.Join(s.Config.BuildTemplatesRootDirectory, semver.String())
}

func (s *Service) archivePath(name string) string {
	return filepath.Join(s.Config.CacheDirectory, CACHE_FOLDER, name)
}

func (s *Service) tempDirectory() string {
	return filepath.Join(s.Config.BuildTemplatesRootDirectory, TEMP_FOLDER)
}

func (s *Service) rootDirectory() string {
	return s.Config.BuildTemplatesRootDirectory
}

func New(environment *environment.Environment, config *config.Config) *Service {
	return &Service{
		Environment: environment,
		Config:      config,
	}
}
