package godot

import (
	"fmt"

	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/semver"
)

type Download struct {
	Version                string `arg:"" help:"Godot version to download to cache in the format x.x.x.x, x.x.x or x.x"`
	IncludeExportTemplates bool   `help:"Include export templates in download"`
	Release                string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                   bool   `help:"Use mono version"`
}

func (c *Download) Run(app *gdvm.App) error {
	if c.IncludeExportTemplates {
		err := app.ExportTemplates.Download(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot download export templates: %w", err)
		}
	}

	err := app.Godot.Download(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot download godot: %w", err)
	}

	return nil
}

type Uninstall struct {
	Version                string `arg:"" help:"Godot version to uninstall in the format x.x.x.x, x.x.x or x.x"`
	ExcludeExportTemplates bool   `help:"Exclude export templates in uninstall"`
	ExcludeShortcuts       bool   `help:"Exclude shortcuts in uninstall"`
	Release                string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                   bool   `help:"Use mono version"`
}

func (c *Uninstall) Run(app *gdvm.App) error {
	if !c.ExcludeExportTemplates {
		err := app.ExportTemplates.Uninstall(semver.Maybe(c.Version, c.Release, c.Mono), true)
		if err != nil {
			return fmt.Errorf("cannot uninstall export templates: %w", err)
		}
	}

	err := app.Godot.Uninstall(semver.Maybe(c.Version, c.Release, c.Mono), true)
	if err != nil {
		return fmt.Errorf("cannot uninstall godot: %w", err)
	}

	if !c.ExcludeShortcuts {
		err := app.Shortcuts.Application.Remove(semver.Maybe(c.Version, c.Release, c.Mono), true)
		if err != nil {
			return fmt.Errorf("cannot remove application shortcut: %w", err)
		}
	}

	if !c.ExcludeShortcuts {
		err := app.Shortcuts.Desktop.Remove(semver.Maybe(c.Version, c.Release, c.Mono), true)
		if err != nil {
			return fmt.Errorf("cannot remove desktop shortcut: %w", err)
		}
	}

	return nil
}

type Install struct {
	Version                string `arg:"" help:"Godot version to download and install in the format x.x.x.x, x.x.x or x.x"`
	IncludeExportTemplates bool   `help:"Include export templates in install"`
	Application            bool   `help:"Add application shortcut"`
	Desktop                bool   `help:"Add desktop shortcut"`
	Release                string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono                   bool   `help:"Use mono version"`
}

func (c *Install) Run(app *gdvm.App) error {
	if c.IncludeExportTemplates {
		err := app.ExportTemplates.Install(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot install export templates: %w", err)
		}
	}

	err := app.Godot.Install(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot install godot: %w", err)
	}

	if c.Application {
		err := app.Shortcuts.Application.Add(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot add application shortcut: %w", err)
		}
	}

	if c.Desktop {
		err := app.Shortcuts.Desktop.Add(semver.Maybe(c.Version, c.Release, c.Mono))
		if err != nil {
			return fmt.Errorf("cannot add desktop shortcut: %w", err)
		}
	}
	return nil
}

type Use struct {
	Version string `arg:"" help:"Godot version to use in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Use) Run(app *gdvm.App) error {
	err := app.Godot.Use(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot use godot: %w", err)
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
	Use       Use       `cmd:"" help:"Use godot engine by version"`
	List      List      `cmd:"" help:"List all current godot versions"`
}
