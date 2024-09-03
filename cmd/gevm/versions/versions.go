package versions

import (
	"fmt"

	"github.com/bashidogames/gevm"
)

type Detailed struct {
	All  bool `help:"View all versions (otherwise only view stable versions)"`
	Mono bool `help:"View mono versions"`
}

func (c *Detailed) Run(app *gevm.App) error {
	err := app.Versions.Detailed(c.All, c.Mono)
	if err != nil {
		return fmt.Errorf("cannot view detailed versions: %w", err)
	}

	return nil
}

type List struct {
	All  bool `help:"List all versions (otherwise only list stable versions)"`
	Mono bool `help:"List mono versions"`
}

func (c *List) Run(app *gevm.App) error {
	err := app.Versions.List(c.All, c.Mono)
	if err != nil {
		return fmt.Errorf("cannot list versions: %w", err)
	}

	return nil
}

type Versions struct {
	Detailed Detailed `cmd:"" help:"View detailed available versions"`
	List     List     `cmd:"" help:"List available versions"`
}
