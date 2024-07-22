package repository

import (
	"slices"

	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/semver"
)

type Repository struct {
	Downloads map[semver.Relver]Download
}

func (r Repository) SortedDownloadKeys() []semver.Relver {
	var keys []semver.Relver
	for key := range r.Downloads {
		keys = append(keys, key)
	}

	slices.SortStableFunc(keys, func(a semver.Relver, b semver.Relver) int { return a.Compare(b) })

	return keys
}

type Download struct {
	MonoAssets map[platform.Platform]Asset
	Assets     map[platform.Platform]Asset
}

func (d Download) HasMonoAsset(platform platform.Platform) bool {
	return d.MonoAssets[platform].IsValid()
}

func (d Download) HasAsset(platform platform.Platform) bool {
	return d.Assets[platform].IsValid()
}

type Asset struct {
	DownloadURL string
	Name        string
}

func (a Asset) IsValid() bool {
	return len(a.DownloadURL) > 0 && len(a.Name) > 0
}

type Entry struct {
	Platform platform.Platform
	Relver   semver.Relver
	Asset    Asset
	Mono     bool
}
