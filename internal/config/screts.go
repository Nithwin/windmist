package config

import (
	"os"
)

// APIKey returns the API key for the configured provider.
func APIKey(cfg *Config) string {
	
	switch cfg.AI.Provider {

	case "gemini":

		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			return key
		}

	case "groq":

		if key := os.Getenv("GROQ_API_KEY"); key != "" {
			return key
		}

	}

	return cfg.AI.APIKey
}
