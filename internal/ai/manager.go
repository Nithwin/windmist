package ai

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
)

// Factory creates an AI provider.
type Factory func(config.ProviderConfig) Provider

// factories maps provider names to their constructors.
var factories = map[string]Factory{
	"gemini": NewGemini,
	// "groq":   NewGroq,
	// "ollama": NewOllama,
}

// New creates the configured AI provider.
func New(cfg *config.Config) (Provider, error) {
	providerCfg, err := cfg.ActiveProvider()
	if err != nil {
		return nil, err
	}

	factory, ok := factories[cfg.AI.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", cfg.AI.Provider)
	}

	return factory(providerCfg), nil
}