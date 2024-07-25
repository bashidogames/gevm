package application

import (
	"fmt"

	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/semver"
)

type Remove struct {
	Version string `arg:"" help:"Application shortcut to remove in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Remove) Run(app *gdvm.App) error {
	err := app.Shortcuts.Application.Remove(semver.Maybe(c.Version, c.Release, c.Mono), true)
	if err != nil {
		return fmt.Errorf("cannot remove application shortcut: %w", err)
	}

	return nil
}

type Add struct {
	Version string `arg:"" help:"Application shortcut to add in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Add) Run(app *gdvm.App) error {
	err := app.Shortcuts.Application.Add(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot add application shortcut: %w", err)
	}

	return nil
}

type Application struct {
	Remove Remove `cmd:"" help:"Remove application shortcuts for godot by version"`
	Add    Add    `cmd:"" help:"Add application shortcuts for godot by version"`
}
