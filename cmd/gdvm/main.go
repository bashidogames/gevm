package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/cmd/gdvm/cache"
	"github.com/bashidogames/gdvm/cmd/gdvm/exporttemplates"
	"github.com/bashidogames/gdvm/cmd/gdvm/godot"
	"github.com/bashidogames/gdvm/cmd/gdvm/settings"
	"github.com/bashidogames/gdvm/cmd/gdvm/shortcuts"
	"github.com/bashidogames/gdvm/cmd/gdvm/version"
	"github.com/bashidogames/gdvm/cmd/gdvm/versions"
	"github.com/bashidogames/gdvm/config"
)

var CLI struct {
	Versions        versions.Versions               `cmd:"" help:"View available versions for download"`
	ExportTemplates exporttemplates.ExportTemplates `cmd:"" help:"Run commands related to export templates"`
	Godot           godot.Godot                     `cmd:"" help:"Run commands related to godot engines" default:"withargs"`
	Shortcuts       shortcuts.Shortcuts             `cmd:"" help:"Remove and add godot shortcuts"`
	Settings        settings.Settings               `cmd:"" help:"View and adjust config settings"`
	Cache           cache.Cache                     `cmd:"" help:"Run commands on the cache"`
	Version         version.Version                 `cmd:"" help:"Print current version"`

	ConfigPath string `help:"Override which config path to use"`
	Verbose    bool   `help:"Use verbose debug logging"`
}

func main() {
	ctx := kong.Parse(&CLI)

	config, err := config.New(
		config.OptionSetConfigPath(CLI.ConfigPath),
		config.OptionSetVerbose(CLI.Verbose),
	)
	if err != nil {
		log.Fatalf("failed to create config: %s", err)
	}

	app, err := gdvm.New(config)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	err = ctx.Run(app)
	if err != nil {
		log.Fatalf("failed to run: %s", err)
	}
}
