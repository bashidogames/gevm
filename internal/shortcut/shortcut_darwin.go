package shortcut

import (
	"os"
)

func Create(shortcutPath string, targetPath string, shortcutName string) error {
	return os.Symlink(targetPath, shortcutPath)
}
