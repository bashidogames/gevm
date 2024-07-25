package exporttemplates

import (
	"fmt"

	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/semver"
)

type Download struct {
	Version string `arg:"" help:"Export templates version to download to cache in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Download) Run(app *gdvm.App) error {
	err := app.ExportTemplates.Download(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot download export templates: %w", err)
	}

	return nil
}

type Uninstall struct {
	Version string `arg:"" help:"Export templates version to uninstall in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Uninstall) Run(app *gdvm.App) error {
	err := app.ExportTemplates.Uninstall(semver.Maybe(c.Version, c.Release, c.Mono), true)
	if err != nil {
		return fmt.Errorf("cannot uninstall export templates: %w", err)
	}

	return nil
}

type Install struct {
	Version string `arg:"" help:"Export templates version to download and install in the format x.x.x.x, x.x.x or x.x"`
	Release string `default:"stable" help:"Release to use (dev1, alpha2, beta3, rc4, stable, etc)"`
	Mono    bool   `help:"Use mono version"`
}

func (c *Install) Run(app *gdvm.App) error {
	err := app.ExportTemplates.Install(semver.Maybe(c.Version, c.Release, c.Mono))
	if err != nil {
		return fmt.Errorf("cannot install export templates: %w", err)
	}

	return nil
}

type List struct{}

func (c *List) Run(app *gdvm.App) error {
	err := app.ExportTemplates.List()
	if err != nil {
		return fmt.Errorf("cannot list export templates: %w", err)
	}

	return nil
}

type ExportTemplates struct {
	Download  Download  `cmd:"" help:"Download export templates to the cache by version"`
	Uninstall Uninstall `cmd:"" help:"Uninstall export templates by version"`
	Install   Install   `cmd:"" help:"Install export templates by version"`
	List      List      `cmd:"" help:"List all current export template versions"`
}
