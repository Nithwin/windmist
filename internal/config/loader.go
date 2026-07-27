package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads the configuration from disk and merges it with defaults
// to ensure all required fields are populated even if the config file
// is partial or from an older version.
func Load() (*Config, error) {
	configFile, err := ConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// Start with defaults so any missing fields are populated
	cfg := DefaultConfig()

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// Ensure the Providers map is never nil
	if cfg.Providers == nil {
		cfg.Providers = DefaultConfig().Providers
	}

	return cfg, nil
}
