package fetcher

import (
	"github.com/bashidogames/gevm/internal/platform"
	"github.com/bashidogames/gevm/internal/repository"
	"github.com/bashidogames/gevm/semver"
)

type Fetcher interface {
	FetchExportTemplatesAsset(semver semver.Semver) (*repository.Asset, error)
	FetchGodotAsset(semver semver.Semver) (*repository.Asset, error)
	FetchRepository(func(entry *Entry) error) error
}

type Entry struct {
	Platform platform.Platform
	Relver   semver.Relver
	Asset    repository.Asset
	Mono     bool
}
