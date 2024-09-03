package shortcuts

import (
	"github.com/bashidogames/gevm/cmd/gevm/shortcuts/application"
	"github.com/bashidogames/gevm/cmd/gevm/shortcuts/desktop"
)

type Shortcuts struct {
	Application application.Application `cmd:"" help:"Remove and add application shortcuts"`
	Desktop     desktop.Desktop         `cmd:"" help:"Remove and add desktop shortcuts"`
}
