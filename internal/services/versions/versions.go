package versions

import (
	"fmt"
	"os"

	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/environment"
	"github.com/bashidogames/gevm/internal/platform"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Service struct {
	Environment *environment.Environment
	Config      *config.Config
}

func (s *Service) Detailed(all bool, mono bool) error {
	repository, err := s.Environment.FetchRepository()
	if err != nil {
		return fmt.Errorf("cannot fetch repository: %w", err)
	}

	header := table.Row{"Version", "Release"}
	for _, platform := range platform.Platforms {
		header = append(header, platform)
	}

	t := table.NewWriter()
	t.AppendHeader(header)

	for _, relver := range repository.SortedDownloadKeys() {
		if relver.Release.Original != "stable" && !all {
			continue
		}

		download := repository.Downloads[relver]

		availabilities := map[platform.Platform]bool{}
		for _, platform := range platform.Platforms {
			if mono {
				availabilities[platform] = download.HasMonoAsset(platform)
			} else {
				availabilities[platform] = download.HasAsset(platform)
			}
		}

		row := table.Row{
			relver.Version,
			relver.Release,
		}

		for _, platform := range platform.Platforms {
			row = append(row, availabilities[platform])
		}

		t.AppendRow(row)
	}

	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}

func (s *Service) List(all bool, mono bool) error {
	repository, err := s.Environment.FetchRepository()
	if err != nil {
		return fmt.Errorf("cannot fetch repository: %w", err)
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Version", "Release", s.Config.Platform})

	for _, relver := range repository.SortedDownloadKeys() {
		if relver.Release.Original != "stable" && !all {
			continue
		}

		download := repository.Downloads[relver]

		var available bool
		if mono {
			available = download.HasMonoAsset(s.Config.Platform)
		} else {
			available = download.HasAsset(s.Config.Platform)
		}

		row := table.Row{
			relver.Version,
			relver.Release,
			available,
		}

		t.AppendRow(row)
	}

	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}

func New(environment *environment.Environment, config *config.Config) *Service {
	return &Service{
		Environment: environment,
		Config:      config,
	}
}
