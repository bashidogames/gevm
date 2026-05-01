package cache

import (
	"fmt"
	"os"

	"github.com/bashmills/gevm/config"
)

type Service struct {
	Config *config.Config
}

func (s *Service) Clear() error {
	s.Config.Logger.Debugf("Attempting to clear cache directory: %s", s.Config.CacheDirectory)

	err := os.RemoveAll(s.Config.CacheDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove cache directory: %w", err)
	}

	s.Config.Logger.Info("Cache cleared")
	return nil
}

func New(config *config.Config) *Service {
	return &Service{
		Config: config,
	}
}
