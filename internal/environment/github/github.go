package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"github.com/bashmills/gevm/config"
	"github.com/bashmills/gevm/internal/downloading"
	"github.com/bashmills/gevm/internal/environment/github/mappings"
	"github.com/bashmills/gevm/internal/platform"
	"github.com/bashmills/gevm/internal/repository"
	"github.com/bashmills/gevm/semver"
)

const REPOSITORY_URL = "https://api.github.com/repos/godotengine/godot-builds/releases?per_page=100"
const ASSET_URL = "https://api.github.com/repos/godotengine/godot-builds/releases/tags/%s"
const ASSET_REGEX_PATTERN = "([-_.]mono)?[-_.](export|linux|x11|windows|win|macos|osx)([-_.]?x86)?[-_.]?(templates|universal|fat|arm64|64)([-_.]?exe)?.(tpz|zip)"
const NEXT_REGEX_PATTERN = "<([^>]*)>[^<]*(next)"
const OLD_REGEX_PATTERN = "^(OLD)[-_.]"

var AssetRegex = regexp.MustCompile(ASSET_REGEX_PATTERN)
var NextRegex = regexp.MustCompile(NEXT_REGEX_PATTERN)
var OldRegex = regexp.MustCompile(OLD_REGEX_PATTERN)

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

func (g *Github) FetchAsset(platform platform.Platform, semver semver.Semver) (*repository.Asset, error) {
	g.Config.Logger.Tracef("Fetching '%s' assets for platform: %s", semver.Relver.GodotString(), platform)

	mapping, ok := mappings.Mappings[platform]
	if !ok {
		return nil, fmt.Errorf("invalid platform mapping: %s", platform)
	}

	url := fmt.Sprintf(ASSET_URL, semver.Relver.GodotString())
	var data Data

	g.Config.Logger.Tracef("Fetching data from url: %s", url)

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

		if OldRegex.MatchString(asset.Name) {
			continue
		}

		isMono := len(parts[1]) > 0
		system := parts[2]
		arch := parts[4]

		if slices.Index(mapping.System, system) < 0 {
			g.Config.Logger.Tracef("Invalid system for asset: %s", asset.Name)
			continue
		}

		if slices.Index(mapping.Arch, arch) < 0 {
			g.Config.Logger.Tracef("Invalid arch for asset: %s", asset.Name)
			continue
		}

		if semver.Mono != isMono {
			g.Config.Logger.Tracef("Invalid mono for asset: %s", asset.Name)
			continue
		}

		g.Config.Logger.Tracef("Asset found: %s", asset.Name)

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

func (g *Github) FetchDownloads(mono bool) ([]repository.Download, error) {
	url := REPOSITORY_URL
	var datas []Data

	for {
		g.Config.Logger.Tracef("Fetching data from url: %s", url)

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
			return nil, fmt.Errorf("fetch failed: %w", err)
		}
	}

	var downloads []repository.Download

	for _, data := range datas {
		relver, err := semver.ParseRelver(data.Name)
		if err != nil {
			return nil, fmt.Errorf("could not parse version release: %w", err)
		}

		download := repository.Download{
			Assets: map[platform.Platform]repository.Asset{},
			Relver: relver,
		}

		for _, asset := range data.Assets {
			parts := AssetRegex.FindStringSubmatch(asset.Name)
			if len(parts) == 0 {
				continue
			}

			if OldRegex.MatchString(asset.Name) {
				continue
			}

			isMono := len(parts[1]) > 0
			if isMono != mono {
				continue
			}

			system := parts[2]
			arch := parts[4]
			found := false

			for platform, mapping := range mappings.Mappings {
				if slices.Index(mapping.System, system) < 0 {
					continue
				}

				if slices.Index(mapping.Arch, arch) < 0 {
					continue
				}

				existing, exists := download.Assets[platform]
				if exists {
					override := mappings.Overrides[platform]
					if len(override) <= 0 {
						g.Config.Logger.Warningf("Asset already exists for '%s' platform: %s == %s", platform, existing.Name, asset.Name)
						continue
					}

					if slices.Index(override, arch) < 0 {
						continue
					}
				}

				download.Assets[platform] = repository.Asset{
					DownloadURL: asset.DownloadURL,
					Name:        asset.Name,
				}

				found = true
			}

			if !found {
				g.Config.Logger.Warningf("No mapping found for asset: %s", asset.Name)
			}
		}

		downloads = append(downloads, download)
	}

	return downloads, nil
}

func New(config *config.Config) *Github {
	return &Github{
		Config: config,
	}
}
