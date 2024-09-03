package desktop

import (
	"fmt"

	"github.com/bashidogames/gevm"
	"github.com/bashidogames/gevm/semver"
)

type Remove struct {
	Version string `arg:"" help:"Desktop shortcut to remove in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Remove) Run(app *gevm.App) error {
	err := app.Shortcuts.Desktop.Remove(semver.Maybe(c.Version, c.Release, c.Mono), true)
	if err != nil {
		return fmt.Errorf("cannot remove desktop shortcut: %w", err)
	}

	return nil
}

type Add struct {
	Version string `arg:"" help:"Desktop shortcut to add in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Add) Run(app *gevm.App) error {
	err := app.Shortcuts.Desktop.Add(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot add desktop shortcut: %w", err)
	}

	return nil
}

type Desktop struct {
	Remove Remove `cmd:"" help:"Remove desktop shortcuts for godot by version"`
	Add    Add    `cmd:"" help:"Add desktop shortcuts for godot by version"`
}
