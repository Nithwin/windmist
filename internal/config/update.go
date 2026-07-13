package config

import "fmt"

// Update updates a configuration value.
func Update(key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	switch key {
	case "provider":
		cfg.AI.Provider = value
	
	case "model":
		cfg.AI.Model = value

	case "theme":
		cfg.UI.Theme = value

	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return Save(cfg)
}