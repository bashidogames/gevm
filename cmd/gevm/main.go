package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/bashidogames/gevm"
	"github.com/bashidogames/gevm/cmd/gevm/cache"
	"github.com/bashidogames/gevm/cmd/gevm/exporttemplates"
	"github.com/bashidogames/gevm/cmd/gevm/godot"
	"github.com/bashidogames/gevm/cmd/gevm/settings"
	"github.com/bashidogames/gevm/cmd/gevm/shortcuts"
	"github.com/bashidogames/gevm/cmd/gevm/version"
	"github.com/bashidogames/gevm/cmd/gevm/versions"
	"github.com/bashidogames/gevm/config"
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

	app, err := gevm.New(config)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	err = ctx.Run(app)
	if err != nil {
		log.Fatalf("failed to run: %s", err)
	}
}
