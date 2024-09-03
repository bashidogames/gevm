package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/downloading"
	"github.com/bashidogames/gevm/internal/environment/fetcher"
	"github.com/bashidogames/gevm/internal/environment/fetcher/github/mappings"
	"github.com/bashidogames/gevm/internal/platform"
	"github.com/bashidogames/gevm/internal/repository"
	"github.com/bashidogames/gevm/internal/utils"
	"github.com/bashidogames/gevm/semver"
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

func (g *Github) FetchExportTemplatesAsset(semver semver.Semver) (*repository.Asset, error) {
	return g.fetchAsset(platform.ExportTemplates, semver)
}

func (g *Github) FetchGodotAsset(semver semver.Semver) (*repository.Asset, error) {
	return g.fetchAsset(g.Config.Platform, semver)
}

func (g *Github) FetchRepository(callback func(entry *fetcher.Entry) error) error {
	url := REPOSITORY_URL
	var datas []Data

	for {
		if g.Config.Verbose {
			utils.Printlnf("Fetching data from url: %s", url)
		}

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

		for platform, mapping := range mappings.Mappings {
			for _, asset := range data.Assets {
				parts := AssetRegex.FindStringSubmatch(asset.Name)
				if len(parts) == 0 {
					continue
				}

				mono := len(parts[1]) > 0
				system := parts[2]
				arch := parts[4]

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

func (g *Github) fetchAsset(platform platform.Platform, semver semver.Semver) (*repository.Asset, error) {
	if g.Config.Verbose {
		utils.Printlnf("Fetching '%s' assets for platform: %s", semver.Relver.GodotString(), platform)
	}

	mapping, ok := mappings.Mappings[platform]
	if !ok {
		return nil, fmt.Errorf("invalid platform mapping: %s", platform)
	}

	url := fmt.Sprintf(ASSET_URL, semver.Relver.GodotString())
	var data Data

	if g.Config.Verbose {
		utils.Printlnf("Fetching assets from url: %s", url)
	}

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
			if g.Config.Verbose {
				utils.Printlnf("Asset not matched: %s", asset.Name)
			}

			continue
		}

		system := parts[2]
		arch := parts[4]

		if slices.Index(mapping.System, system) < 0 {
			if g.Config.Verbose {
				utils.Printlnf("Asset matched but invalid system: %s", asset.Name)
			}

			continue
		}

		if slices.Index(mapping.Arch, arch) < 0 {
			if g.Config.Verbose {
				utils.Printlnf("Asset matched but invalid arch: %s", asset.Name)
			}

			continue
		}

		if g.Config.Verbose {
			utils.Printlnf("Asset found: %s", asset.Name)
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
		return nil, fmt.Errorf("multiple assets found: %s", semver.GodotString())
	}

	return &assets[0], nil
}

func New(config *config.Config) *Github {
	return &Github{
		Config: config,
	}
}
