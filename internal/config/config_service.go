package config

import (
	"fmt"
	"os"
)

type configService struct {
	repo ViperRepo
}

var defaultService Service = NewService(NewViperRepo())

func NewService(r ViperRepo) Service {
	return configService{repo: r}
}

func DefaultService() Service {
	return defaultService
}

func (s configService) Load(dataDir string) (*Config, error) {
	return s.repo.Load(dataDir)
}

func (s configService) Save(dataDir string, cfg *Config) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	return s.repo.Save(dataDir, cfg)
}

func (s configService) Get(dataDir, key string) (interface{}, error) {
	val, err := s.repo.Get(dataDir, key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return val, nil
}

func (s configService) Set(dataDir, key, value string) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	return s.repo.Set(dataDir, key, value)
}
