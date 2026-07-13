package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Save writes the configuration to disk.
func Save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	configFile, err := ConfigFile()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFile, data, 0644)

	if err != nil {
		return err
	}

	return nil
}