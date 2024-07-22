package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/bashidogames/gdvm"
	"github.com/bashidogames/gdvm/cmd/gdvm/cache"
	"github.com/bashidogames/gdvm/cmd/gdvm/godot"
	"github.com/bashidogames/gdvm/cmd/gdvm/versions"
	"github.com/bashidogames/gdvm/config"
)

var CLI struct {
	Versions versions.Versions `cmd:"" help:"View available versions for download"`
	Godot    godot.Godot       `cmd:"" help:"Run commands related to the godot engine"`
	Cache    cache.Cache       `cmd:"" help:"Run commands on the cache"`

	Verbose bool `help:"Use verbose debug logging"`
}

func main() {
	ctx := kong.Parse(&CLI)

	config, err := config.New(
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
