package shortcut

import (
	"runtime"

	"github.com/jxeng/shortcut"
)

func Create(shortcutPath string, targetPath string, shortcutName string) error {
	// This is needed to fix "CoInitialize has not been called" bug (https://github.com/go-ole/go-ole/issues/124#issuecomment-339117036)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	return shortcut.Create(shortcut.Shortcut{
		ShortcutPath: shortcutPath,
		IconLocation: targetPath,
		Target:       targetPath,
	})
}
