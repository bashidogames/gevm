package gdvm

import (
	"fmt"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/environment"
	"github.com/bashidogames/gdvm/internal/environment/fetcher/github"
	"github.com/bashidogames/gdvm/internal/services/buildtemplates"
	"github.com/bashidogames/gdvm/internal/services/cache"
	"github.com/bashidogames/gdvm/internal/services/godot"
	"github.com/bashidogames/gdvm/internal/services/settings"
	"github.com/bashidogames/gdvm/internal/services/shortcuts"
	"github.com/bashidogames/gdvm/internal/services/versions"
)

type App struct {
	BuildTemplates *buildtemplates.Service
	Godot          *godot.Service
	Shortcuts      *shortcuts.Service
	Versions       *versions.Service
	Settings       *settings.Service
	Cache          *cache.Service
}

func New(config *config.Config) (*App, error) {
	environment, err := environment.New(github.New(config))
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return &App{
		BuildTemplates: buildtemplates.New(environment, config),
		Godot:          godot.New(environment, config),
		Shortcuts:      shortcuts.New(config),
		Versions:       versions.New(environment, config),
		Settings:       settings.New(config),
		Cache:          cache.New(config),
	}, nil
}
