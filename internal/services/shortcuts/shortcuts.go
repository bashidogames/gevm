package shortcuts

import (
	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/services/shortcuts/application"
	"github.com/bashidogames/gdvm/internal/services/shortcuts/desktop"
	"github.com/bashidogames/gdvm/internal/services/shortcuts/fetcher"
)

type Service struct {
	Application *application.Service
	Desktop     *desktop.Service
}

func New(config *config.Config) *Service {
	fetcher := fetcher.New(config)
	return &Service{
		Application: application.New(fetcher, config),
		Desktop:     desktop.New(fetcher, config),
	}
}
