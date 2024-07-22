package environment

import (
	"fmt"

	"github.com/bashidogames/gdvm/internal/github"
	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/internal/repository"
	"github.com/bashidogames/gdvm/semver"
)

type Environment struct {
	Github *github.Github
}

func (e *Environment) FetchRepository() (*repository.Repository, error) {
	repo := repository.Repository{
		Downloads: map[semver.Relver]repository.Download{},
	}

	err := e.Github.FetchRepository(func(entry *repository.Entry) error {
		download, ok := repo.Downloads[entry.Relver]
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

		repo.Downloads[entry.Relver] = download
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	return &repo, nil
}

func (e *Environment) FetchAsset(semver semver.Semver) (*repository.Asset, error) {
	asset, err := e.Github.FetchAsset(semver)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset: %w", err)
	}

	return asset, nil
}

func New(github *github.Github) (*Environment, error) {
	return &Environment{
		Github: github,
	}, nil
}
