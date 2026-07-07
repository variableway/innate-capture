package config

// Service defines the config domain API.
type Service interface {
	Load(dataDir string) (*Config, error)
	Save(dataDir string, cfg *Config) error
	Get(dataDir, key string) (interface{}, error)
	Set(dataDir, key, value string) error
}
