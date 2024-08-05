package fetcher

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/utils"
	"github.com/bashidogames/gdvm/semver"
)

const EXECUTABLE_REGEX_PATTERN = "Godot(.*?)([-_.]mono)?[-_.](linux|x11)([-_.]?x86)?[-_.]?(arm64|arm32|64|32)"
const SHORTCUT_FILENAME = "godot-%s.desktop"
const SHORTCUT_NAME = "Godot %s"

var ExecutableRegex = regexp.MustCompile(EXECUTABLE_REGEX_PATTERN)

type Fetcher struct {
	Config *config.Config
}

func (f *Fetcher) ApplicationShortcutPath(semver semver.Semver) string {
	return filepath.Join(f.Config.ApplicationShortcutDirectory, f.shortcutFilename(semver))
}

func (f *Fetcher) DesktopShortcutPath(semver semver.Semver) string {
	return filepath.Join(f.Config.DesktopShortcutDirectory, f.shortcutFilename(semver))
}

func (f *Fetcher) TargetPath(semver semver.Semver) (string, error) {
	return f.locateExecutable(filepath.Join(f.Config.GodotRootDirectory, semver.GodotString()))
}

func (f *Fetcher) ShortcutName(semver semver.Semver) string {
	return fmt.Sprintf(SHORTCUT_NAME, semver.GodotString())
}

func (f *Fetcher) shortcutFilename(semver semver.Semver) string {
	return fmt.Sprintf(SHORTCUT_FILENAME, semver.GodotString())
}

func (f *Fetcher) locateExecutable(root string) (string, error) {
	return utils.LocateExecutable(ExecutableRegex, root, false)
}

func New(config *config.Config) *Fetcher {
	return &Fetcher{
		Config: config,
	}
}
