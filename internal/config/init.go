package config

import (
	"os"
)

// Init initializes the WindMist configuration directory structure (e.g., ~/.config/windmist).
// It ensures that all required configuration directories exist with appropriate file permissions (0755).
func Init() error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile, err := ConfigFile()
	if err != nil {
		return err
	}

	_, err = os.Stat(configFile)

	if os.IsNotExist(err) {
		cfg := DefaultConfig()

		if err := Save(cfg); err != nil {
			return err
		}
	}

	return nil
}
