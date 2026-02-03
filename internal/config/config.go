package config

import (
	"os"

	"github.com/adminlove520/vulnDb-Notifier/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
	PushMode string   `yaml:"push_mode,omitempty"` // daily 或 keyword，默认为 daily
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, &errors.ConfigError{Message: "Failed to read config file: " + err.Error()}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &errors.ConfigError{Message: "Failed to unmarshal config data: " + err.Error()}
	}

	return &cfg, nil
}

func GetConfigPath() string {
	configPath := "config.yaml"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		configPath = envConfigPath
	}
	return configPath
}
