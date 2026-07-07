package config

import "github.com/spf13/viper"

type ViperRepo interface {
	Load(dataDir string) (*Config, error)
	Save(dataDir string, cfg *Config) error
	Get(dataDir, key string) (interface{}, error)
	Set(dataDir, key, value string) error
}

type viperRepo struct{}

func NewViperRepo() ViperRepo {
	return viperRepo{}
}

func (viperRepo) newViper(dataDir string) *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dataDir)
	return v
}

func (r viperRepo) Load(dataDir string) (*Config, error) {
	v := r.newViper(dataDir)
	cfg := DefaultConfig()
	if err := v.ReadInConfig(); err != nil {
		return cfg, nil
	}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (r viperRepo) Save(dataDir string, cfg *Config) error {
	v := r.newViper(dataDir)
	v.Set("app", cfg.App)
	v.Set("defaults", cfg.Defaults)
	v.Set("feishu", cfg.Feishu)
	v.Set("bitable", cfg.Bitable)
	v.Set("bot", cfg.Bot)
	v.Set("workspace", cfg.Workspace)
	return v.WriteConfigAs(dataDir + "/config.yaml")
}

func (r viperRepo) Get(dataDir, key string) (interface{}, error) {
	v := r.newViper(dataDir)
	if err := v.ReadInConfig(); err != nil {
		cfg := DefaultConfig()
		_ = v.MergeConfigMap(map[string]interface{}{
			"app":       cfg.App,
			"defaults":  cfg.Defaults,
			"feishu":    cfg.Feishu,
			"bitable":   cfg.Bitable,
			"bot":       cfg.Bot,
			"workspace": cfg.Workspace,
		})
	}
	return v.Get(key), nil
}

func (r viperRepo) Set(dataDir, key, value string) error {
	v := r.newViper(dataDir)
	_ = v.ReadInConfig()
	v.Set(key, value)
	return v.WriteConfigAs(dataDir + "/config.yaml")
}
