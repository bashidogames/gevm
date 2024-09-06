package fetcher

import (
	"path/filepath"
	"regexp"

	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/utils"
	"github.com/bashidogames/gevm/semver"
)

const EXECUTABLE_REGEX_PATTERN = "Godot(.*?)([-_.]mono)?[-_.](linux|x11)([-_.]?x86)?[-_.]?(arm64|arm32|64|32)"
const LINK_FILENAME = "godot"

var ExecutableRegex = regexp.MustCompile(EXECUTABLE_REGEX_PATTERN)

type Fetcher struct {
	Config *config.Config
}

func (f *Fetcher) TargetPath(semver semver.Semver) (string, error) {
	return f.locateExecutable(filepath.Join(f.Config.GodotRootDirectory, semver.GodotString()))
}

func (f *Fetcher) LinkPath(semver semver.Semver) string {
	return filepath.Join(f.Config.BinDirectory, LINK_FILENAME)
}

func (f *Fetcher) locateExecutable(root string) (string, error) {
	return utils.LocateExecutable(func(filename string) bool {
		return ExecutableRegex.MatchString(filename)
	}, root, false)
}

func New(config *config.Config) *Fetcher {
	return &Fetcher{
		Config: config,
	}
}
