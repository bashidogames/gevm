package shortcut

import (
	"fmt"
	"os"

	"github.com/bashidogames/gdvm/internal/utils"
)

const CONTENTS = `[Desktop Entry]
Name=%s
Comment=The game engine you've been waiting for.
GenericName=Game Engine
Exec=%s
Icon=godot
Type=Application
Categories=Development;IDE;
MimeType=text/plain;inode/directory;application/x-godot-project;
Keywords=godot;`

func Create(shortcutPath string, targetPath string, shortcutName string) error {
	file, err := os.OpenFile(shortcutPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, utils.OS_FILE)
	if err != nil {
		return fmt.Errorf("could not create shortcut file: %w", err)
	}
	defer file.Close()

	contents := fmt.Sprintf(CONTENTS, shortcutName, targetPath)

	_, err = file.WriteString(contents)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}
