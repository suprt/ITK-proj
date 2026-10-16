package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

func Load[T any](yamlPath string) (*T, error) {
	var cfg T

	if err := env.Parse(&cfg); err != nil {
		return &cfg, err
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &cfg, nil
		}
		return &cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
