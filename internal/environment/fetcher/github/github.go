package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/downloading"
	"github.com/bashidogames/gdvm/internal/environment/fetcher"
	"github.com/bashidogames/gdvm/internal/environment/fetcher/github/mappings"
	"github.com/bashidogames/gdvm/internal/platform"
	"github.com/bashidogames/gdvm/internal/repository"
	"github.com/bashidogames/gdvm/semver"
)

const REPOSITORY_URL = "https://api.github.com/repos/godotengine/godot-builds/releases?per_page=100"
const ASSET_URL = "https://api.github.com/repos/godotengine/godot-builds/releases/tags/%s"
const ASSET_REGEX_PATTERN = "([-_.]mono)?[-_.](export|linux|x11|win|macos|osx)([-_.]?x86)?[-_.]?(templates|universal|fat|arm64|arm32|64|32)([-_.]?exe)?.(tpz|zip)"
const NEXT_REGEX_PATTERN = "<([^>]*)>[^<]*(next)"

var AssetRegex = regexp.MustCompile(ASSET_REGEX_PATTERN)
var NextRegex = regexp.MustCompile(NEXT_REGEX_PATTERN)

type Github struct {
	Config *config.Config
}

type Data struct {
	Name   string `json:"tag_name"`
	Assets []struct {
		DownloadURL string `json:"browser_download_url"`
		Name        string `json:"name"`
	} `json:"assets"`
}

func (g *Github) FetchBuildTemplatesAsset(semver semver.Semver) (*repository.Asset, error) {
	return fetchAsset(platform.ExportTemplates, semver)
}

func (g *Github) FetchGodotAsset(semver semver.Semver) (*repository.Asset, error) {
	return fetchAsset(g.Config.Platform, semver)
}

func (g *Github) FetchRepository(callback func(entry *fetcher.Entry) error) error {
	url := REPOSITORY_URL
	var datas []Data

	for {
		err := downloading.Fetch(url, func(header http.Header, bytes []byte) error {
			var data []Data
			err := json.Unmarshal(bytes, &data)
			if err != nil {
				return fmt.Errorf("cannot parse bytes: %w", err)
			}

			datas = append(datas, data...)
			link := header.Get("link")

			parts := NextRegex.FindStringSubmatch(link)
			if len(parts) == 0 {
				return io.EOF
			}

			url = parts[1]
			return nil
		})
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("fetch failed: %w", err)
		}
	}

	for _, data := range datas {
		relver, err := semver.ParseRelver(data.Name)
		if err != nil {
			return fmt.Errorf("could not parse version release: %w", err)
		}

		for _, asset := range data.Assets {
			parts := AssetRegex.FindStringSubmatch(asset.Name)
			if len(parts) == 0 {
				continue
			}

			mono := len(parts[1]) > 0
			system := parts[2]
			arch := parts[4]

			for platform, mapping := range mappings.Mappings {
				if slices.Index(mapping.System, system) < 0 {
					continue
				}

				if slices.Index(mapping.Arch, arch) < 0 {
					continue
				}

				err := callback(&fetcher.Entry{
					Platform: platform,
					Relver:   relver,
					Asset: repository.Asset{
						DownloadURL: asset.DownloadURL,
						Name:        asset.Name,
					},
					Mono: mono,
				})
				if err != nil {
					return fmt.Errorf("repository callback failed: %w", err)
				}
			}
		}
	}

	return nil
}

func fetchAsset(platform platform.Platform, semver semver.Semver) (*repository.Asset, error) {
	mapping, ok := mappings.Mappings[platform]
	if !ok {
		return nil, fmt.Errorf("invalid platform mapping: %s", platform)
	}

	url := fmt.Sprintf(ASSET_URL, semver.Relver)
	var data Data

	err := downloading.Fetch(url, func(header http.Header, bytes []byte) error {
		err := json.Unmarshal(bytes, &data)
		if err != nil {
			return fmt.Errorf("cannot parse bytes: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	var assets []repository.Asset

	for _, asset := range data.Assets {
		parts := AssetRegex.FindStringSubmatch(asset.Name)
		if len(parts) == 0 {
			continue
		}

		mono := len(parts[1]) > 0
		if semver.Mono != mono {
			continue
		}

		system := parts[2]
		arch := parts[4]

		if slices.Index(mapping.System, system) < 0 {
			continue
		}

		if slices.Index(mapping.Arch, arch) < 0 {
			continue
		}

		assets = append(assets, repository.Asset{
			DownloadURL: asset.DownloadURL,
			Name:        asset.Name,
		})
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("fetch asset failure: %w", downloading.ErrNotFound)
	}

	if len(assets) > 1 {
		return nil, fmt.Errorf("multiple assets found: %s", semver)
	}

	return &assets[0], nil
}

func New(config *config.Config) *Github {
	return &Github{
		Config: config,
	}
}
