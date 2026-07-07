package config

import (
	"os"
	"path/filepath"
)

func DefaultDataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".capture")
}

func Load(dataDir string) (*Config, error) {
	return DefaultService().Load(dataDir)
}

func Save(dataDir string, cfg *Config) error {
	return DefaultService().Save(dataDir, cfg)
}

func Get(dataDir, key string) (interface{}, error) {
	return DefaultService().Get(dataDir, key)
}

func Set(dataDir, key, value string) error {
	return DefaultService().Set(dataDir, key, value)
}
