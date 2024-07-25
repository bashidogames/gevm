package shortcuts

import (
	"github.com/bashidogames/gdvm/cmd/gdvm/shortcuts/application"
	"github.com/bashidogames/gdvm/cmd/gdvm/shortcuts/desktop"
)

type Shortcuts struct {
	Application application.Application `cmd:"" help:"Run commands for application shortcuts"`
	Desktop     desktop.Desktop         `cmd:"" help:"Run commands for desktop shortcuts"`
}
