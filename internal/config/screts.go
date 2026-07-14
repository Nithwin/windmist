package config

import "os"

// APIKey returns the configured API key.
func APIKey(cfg *Config) string {

	switch cfg.AI.Provider {

	case "gemini":

		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			return key
		}

		return cfg.Providers.Gemini.APIKey

	case "groq":

		if key := os.Getenv("GROQ_API_KEY"); key != "" {
			return key
		}

		return cfg.Providers.Groq.APIKey

	default:
		return ""
	}
}