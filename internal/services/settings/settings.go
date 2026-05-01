package settings

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bashmills/gevm/config"
	"github.com/bashmills/gevm/internal/utils"
)

var ErrNotFound = errors.New("not found")

type Service struct {
	Config *config.Config
}

func (s *Service) Reset() error {
	s.Config.Logger.Debug("Attempting to reset settings...")

	err := s.Config.Reset()
	if err != nil {
		return fmt.Errorf("cannot reset config: %w", err)
	}

	err = s.Config.Save()
	if err != nil {
		return fmt.Errorf("cannot save config: %w", err)
	}

	s.Config.Logger.Info("Settings reset")
	return nil
}

func (s *Service) List() error {
	err := s.iterateFields(func(field reflect.Value, name string) error {
		utils.Printlnf("%s = %s", name, field.String())
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to iterate fields: %w", err)
	}

	return nil
}

func (s *Service) Set(key string, value string) error {
	err := s.findField(key, func(field reflect.Value, name string) error {
		utils.Printlnf("%s = %s => %s", name, field.String(), value)
		field.SetString(value)
		return nil
	})
	if err == ErrNotFound {
		s.Config.Logger.Errorf("Settings key '%s' not found. Use 'gevm settings list' to see available settings.", key)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to find field: %w", err)
	}

	err = s.Config.Save()
	if err != nil {
		return fmt.Errorf("cannot save config: %w", err)
	}

	return nil
}

func (s *Service) Get(key string) error {
	err := s.findField(key, func(field reflect.Value, name string) error {
		utils.Printlnf("%s = %s", name, field.String())
		return nil
	})
	if err == ErrNotFound {
		s.Config.Logger.Errorf("Settings key '%s' not found. Use 'gevm settings list' to see available settings.", key)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to find field: %w", err)
	}

	return nil
}

func (s *Service) Path() error {
	utils.Println(s.Config.ConfigPath)
	return nil
}

func (s *Service) iterateFields(callback func(reflect.Value, string) error) error {
	element := reflect.ValueOf(s.Config).Elem()
	for i := 0; i < element.Type().NumField(); i++ {
		tag := element.Type().Field(i).Tag.Get("json")
		field := element.Field(i)
		if tag == "-" || tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		err := callback(field, name)
		if err != nil {
			return fmt.Errorf("iterate fields callback failed: %w", err)
		}
	}

	return nil
}

func (s *Service) findField(key string, callback func(reflect.Value, string) error) error {
	found := false
	err := s.iterateFields(func(field reflect.Value, name string) error {
		if name != key {
			return nil
		}

		err := callback(field, name)
		if err != nil {
			return fmt.Errorf("find field callback failed: %w", err)
		}

		found = true
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to iterate fields: %w", err)
	}

	if !found {
		return ErrNotFound
	}

	return nil
}

func New(config *config.Config) *Service {
	return &Service{
		Config: config,
	}
}
