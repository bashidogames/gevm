package godot

import (
	"fmt"

	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/semver"
)

type Download struct {
	Version               string `arg:"" help:"Godot version to download to cache in the format x.x.x.x, x.x.x or x.x"`
	IncludeBuildTemplates bool   `help:"Include build templates in download"`
	Release               string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                  bool   `help:"Use mono version"`
}

func (c *Download) Run(app *gdvm.App) error {
	if c.IncludeBuildTemplates {
		err := app.BuildTemplates.Download(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot download build templates: %w", err)
		}
	}

	err := app.Godot.Download(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot download godot: %w", err)
	}

	return nil
}

type Uninstall struct {
	Version               string `arg:"" help:"Godot version to uninstall in the format x.x.x.x, x.x.x or x.x"`
	ExcludeBuildTemplates bool   `help:"Exclude build templates in uninstall"`
	Release               string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                  bool   `help:"Use mono version"`
}

func (c *Uninstall) Run(app *gdvm.App) error {
	if !c.ExcludeBuildTemplates {
		err := app.BuildTemplates.Uninstall(semver.Maybe(c.Version, c.Release, c.Mono), true)
		if err != nil {
			return fmt.Errorf("cannot uninstall build templates: %w", err)
		}
	}

	err := app.Godot.Uninstall(semver.Maybe(c.Version, c.Release, c.Mono), true)
	if err != nil {
		return fmt.Errorf("cannot uninstall godot: %w", err)
	}

	return nil
}

type Install struct {
	Version               string `arg:"" help:"Godot version to download and install in the format x.x.x.x, x.x.x or x.x"`
	IncludeBuildTemplates bool   `help:"Include build templates in install"`
	Release               string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                  bool   `help:"Use mono version"`
}

func (c *Install) Run(app *gdvm.App) error {
	if c.IncludeBuildTemplates {
		err := app.BuildTemplates.Install(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot install build templates: %w", err)
		}
	}

	err := app.Godot.Install(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot install godot: %w", err)
	}

	return nil
}

type List struct{}

func (c *List) Run(app *gdvm.App) error {
	err := app.Godot.List()
	if err != nil {
		return fmt.Errorf("cannot list godot: %w", err)
	}

	return nil
}

type Godot struct {
	Download  Download  `cmd:"" help:"Download godot engine to the cache by version"`
	Uninstall Uninstall `cmd:"" help:"Uninstall godot engine by version"`
	Install   Install   `cmd:"" help:"Install godot engine by version"`
	List      List      `cmd:"" help:"List all current godot versions"`
}
