package environment

import (
	"fmt"

	"github.com/bashidogames/gevm/internal/environment/fetcher"
	"github.com/bashidogames/gevm/internal/platform"
	"github.com/bashidogames/gevm/internal/repository"
	"github.com/bashidogames/gevm/semver"
)

type Environment struct {
	Fetcher fetcher.Fetcher
}

func (e *Environment) FetchExportTemplatesAsset(semver semver.Semver) (*repository.Asset, error) {
	asset, err := e.Fetcher.FetchExportTemplatesAsset(semver)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch export templates asset: %w", err)
	}

	return asset, nil
}

func (e *Environment) FetchGodotAsset(semver semver.Semver) (*repository.Asset, error) {
	asset, err := e.Fetcher.FetchGodotAsset(semver)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch godot asset: %w", err)
	}

	return asset, nil
}

func (e *Environment) FetchRepository() (*repository.Repository, error) {
	result := repository.Repository{
		Downloads: map[semver.Relver]repository.Download{},
	}

	err := e.Fetcher.FetchRepository(func(entry *fetcher.Entry) error {
		download, ok := result.Downloads[entry.Relver]
		if !ok {
			download = repository.Download{
				MonoAssets: map[platform.Platform]repository.Asset{},
				Assets:     map[platform.Platform]repository.Asset{},
			}
		}

		if entry.Mono {
			download.MonoAssets[entry.Platform] = entry.Asset
		} else {
			download.Assets[entry.Platform] = entry.Asset
		}

		result.Downloads[entry.Relver] = download
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	return &result, nil
}

func New(fetcher fetcher.Fetcher) (*Environment, error) {
	return &Environment{
		Fetcher: fetcher,
	}, nil
}
