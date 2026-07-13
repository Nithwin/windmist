package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	configFile, err := ConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}