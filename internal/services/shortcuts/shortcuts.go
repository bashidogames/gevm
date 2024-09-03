package shortcuts

import (
	"github.com/bashidogames/gevm/config"
	"github.com/bashidogames/gevm/internal/services/shortcuts/application"
	"github.com/bashidogames/gevm/internal/services/shortcuts/desktop"
	"github.com/bashidogames/gevm/internal/services/shortcuts/fetcher"
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
