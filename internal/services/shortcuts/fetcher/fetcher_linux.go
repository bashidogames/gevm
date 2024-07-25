package fetcher

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/bashidogames/gdvm/config"
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
	return f.locateExecutable(filepath.Join(f.Config.GodotRootDirectory, semver.String()))
}

func (f *Fetcher) ShortcutName(semver semver.Semver) string {
	return fmt.Sprintf(SHORTCUT_NAME, semver)
}

func (f *Fetcher) shortcutFilename(semver semver.Semver) string {
	return fmt.Sprintf(SHORTCUT_FILENAME, semver)
}

func (f *Fetcher) locateExecutable(root string) (string, error) {
	var result string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk path: %s", path)
		}

		if d.IsDir() {
			return nil
		}

		if !ExecutableRegex.MatchString(filepath.Base(path)) {
			return nil
		}

		result = path
		return filepath.SkipAll
	})
	if err != nil {
		return "", fmt.Errorf("could not walk directory: %w", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("executable not found")
	}

	return result, nil
}

func New(config *config.Config) *Fetcher {
	return &Fetcher{
		Config: config,
	}
}
