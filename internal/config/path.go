package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDir returns the directory used to store WindMist configuration.
func ConfigDir() (string, error){
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(configDir, "windmist"), nil
}

// ConfigFile returns the path to the WindMist configuration file.
func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}