package fetcher

import (
	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/internal/repository"
	"github.com/bashidogames/gdvm/semver"
)

type Fetcher interface {
	FetchBuildTemplatesAsset(semver semver.Semver) (*repository.Asset, error)
	FetchGodotAsset(semver semver.Semver) (*repository.Asset, error)
	FetchRepository(func(entry *Entry) error) error
}

type Entry struct {
	Platform platform.Platform
	Relver   semver.Relver
	Asset    repository.Asset
	Mono     bool
}
