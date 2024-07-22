package cache

import (
	"fmt"
	"os"

	"github.com/bashidogames/gdvm/config"
	"github.com/bashidogames/gdvm/internal/utils"
)

type Service struct {
	Config *config.Config
}

func (s *Service) Clear() error {
	if s.Config.Verbose {
		utils.Printlnf("Attempting to clear cache directory: %s", s.Config.CacheDirectory)
	}

	err := os.RemoveAll(s.Config.CacheDirectory)
	if err != nil {
		return fmt.Errorf("cannot remove cache directory: %w", err)
	}

	utils.Printlnf("Cache cleared")
	return nil
}

func New(config *config.Config) *Service {
	return &Service{
		Config: config,
	}
}
