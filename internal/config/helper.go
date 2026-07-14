package config

import (
	"fmt"
	"os"
)

// ActiveProvider returns the configuration of the currently selected AI provider.
func (c *Config) ActiveProvider() (*ProviderConfig, error) {
	switch c.AI.Provider {

	case "gemini":
		return &c.Providers.Gemini, nil

	case "groq":
		return &c.Providers.Groq, nil

	case "ollama":
		return &c.Providers.Ollama, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", c.AI.Provider)
	}
}

// ActiveModel returns the model of the currently selected provider.
func (c *Config) ActiveModel() (string, error) {
	provider, err := c.ActiveProvider()
	if err != nil {
		return "", err
	}

	return provider.Model, nil
}

// ActiveBaseURL returns the BaseURL of the currently selected provider.
func (c *Config) ActiveBaseURL() (string, error) {
	provider, err := c.ActiveProvider()
	if err != nil {
		return "", err
	}

	return provider.BaseURL, nil
}

// ActiveAPIKey returns the API key of the active provider.
//
// Priority:
//   1. Environment variable
//   2. config.yaml
func (c *Config) ActiveAPIKey() (string, error) {

	switch c.AI.Provider {

	case "gemini":

		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			return key, nil
		}

		return c.Providers.Gemini.APIKey, nil

	case "groq":

		if key := os.Getenv("GROQ_API_KEY"); key != "" {
			return key, nil
		}

		return c.Providers.Groq.APIKey, nil

	case "ollama":
		// Ollama doesn't require an API key.
		return "", nil

	default:
		return "", fmt.Errorf("unsupported provider: %s", c.AI.Provider)
	}
}