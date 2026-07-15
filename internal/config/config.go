package config

import (
	"fmt"
	"os"
)

const (
	EnvGeminiAPIKey = "GEMINI_API_KEY"
	EnvGroqAPIKey   = "GROQ_API_KEY"
)

var envKeys = map[string]string{
	"gemini": EnvGeminiAPIKey,
	"groq":   EnvGroqAPIKey,
}

// ActiveProvider returns the active provider configuration.
func (c *Config) ActiveProvider() (ProviderConfig, error) {
	provider, ok := c.Providers[c.AI.Provider]
	if !ok {
		return ProviderConfig{}, fmt.Errorf("unsupported provider: %s", c.AI.Provider)
	}

	return provider, nil
}

// ActiveModel returns the active model.
func (c *Config) ActiveModel() (string, error) {
	provider, err := c.ActiveProvider()
	if err != nil {
		return "", err
	}

	return provider.Model, nil
}

// ActiveAPIKey returns the active provider API key.
// Environment variables take precedence over config values.
func (c *Config) ActiveAPIKey() (string, error) {
	provider, err := c.ActiveProvider()
	if err != nil {
		return "", err
	}

	if envVar, ok := envKeys[c.AI.Provider]; ok {
		if key := os.Getenv(envVar); key != "" {
			return key, nil
		}
	}

	return provider.APIKey, nil
}

// ActiveBaseURL returns the active provider base URL.
func (c *Config) ActiveBaseURL() (string, error) {
	provider, err := c.ActiveProvider()
	if err != nil {
		return "", err
	}

	return provider.BaseURL, nil
}

// SetProvider changes the active provider.
func (c *Config) SetProvider(name string) error {
	if _, ok := c.Providers[name]; !ok {
		return fmt.Errorf("unsupported provider: %s", name)
	}

	c.AI.Provider = name
	return nil
}

// SetModel updates a provider model.
func (c *Config) SetModel(providerName, model string) error {
	provider, ok := c.Providers[providerName]
	if !ok {
		return fmt.Errorf("unsupported provider: %s", providerName)
	}

	provider.Model = model
	c.Providers[providerName] = provider

	return nil
}

// SetAPIKey updates a provider API key.
func (c *Config) SetAPIKey(providerName, apiKey string) error {
	provider, ok := c.Providers[providerName]
	if !ok {
		return fmt.Errorf("unsupported provider: %s", providerName)
	}

	provider.APIKey = apiKey
	c.Providers[providerName] = provider

	return nil
}

// SetBaseURL updates a provider base URL.
func (c *Config) SetBaseURL(providerName, baseURL string) error {
	provider, ok := c.Providers[providerName]
	if !ok {
		return fmt.Errorf("unsupported provider: %s", providerName)
	}

	provider.BaseURL = baseURL
	c.Providers[providerName] = provider

	return nil
}

// SetTheme updates the UI theme.
func (c *Config) SetTheme(theme string) {
	c.UI.Theme = theme
}
