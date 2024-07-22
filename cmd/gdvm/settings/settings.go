package settings

import (
	"fmt"

	"github.com/bashidogames/gdvm"
)

type Reset struct{}

func (c *Reset) Run(app *gdvm.App) error {
	err := app.Settings.Reset()
	if err != nil {
		return fmt.Errorf("cannot reset settings: %w", err)
	}

	return nil
}

type List struct{}

func (c *List) Run(app *gdvm.App) error {
	err := app.Settings.List()
	if err != nil {
		return fmt.Errorf("cannot list settings: %w", err)
	}

	return nil
}

type Set struct {
	Key   string `arg:"" help:"Key to set config value for"`
	Value string `arg:"" help:"Value to set"`
}

func (c *Set) Run(app *gdvm.App) error {
	err := app.Settings.Set(c.Key, c.Value)
	if err != nil {
		return fmt.Errorf("cannot set setting: %w", err)
	}

	return nil
}

type Get struct {
	Key string `arg:"" help:"Key to get config value for"`
}

func (c *Get) Run(app *gdvm.App) error {
	err := app.Settings.Get(c.Key)
	if err != nil {
		return fmt.Errorf("cannot get setting: %w", err)
	}

	return nil
}

type Path struct{}

func (c *Path) Run(app *gdvm.App) error {
	err := app.Settings.Path()
	if err != nil {
		return fmt.Errorf("cannot print path: %w", err)
	}

	return nil
}

type Settings struct {
	Reset Reset `cmd:"" help:"Reset config back to defaults"`
	List  List  `cmd:"" help:"List all current config values"`
	Set   Set   `cmd:"" help:"Set config value"`
	Get   Get   `cmd:"" help:"Get config value"`
	Path  Path  `cmd:"" help:"Print path to config file"`
}
