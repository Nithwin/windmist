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

	case "gemini.model":
		cfg.Providers.Gemini.Model = value

	case "gemini.api_key":
		cfg.Providers.Gemini.APIKey = value

	case "groq.model":
		cfg.Providers.Groq.Model = value

	case "groq.api_key":
		cfg.Providers.Groq.APIKey = value

	case "ollama.model":
		cfg.Providers.Ollama.Model = value

	case "ollama.base_url":
		cfg.Providers.Ollama.BaseURL = value

	case "theme":
		cfg.UI.Theme = value

	default:
		return fmt.Errorf("unknown configuration key: %s", key)

	}

	return Save(cfg)
}