package ai

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
)

// Factory creates an AI provider.
type Factory func(config.ProviderConfig) Provider

// factories maps provider names to their constructors.
var factories = make(map[string]Factory)

// Register is called by provider packages in their init() functions to register themselves.
func Register(name string, factory Factory) {
	factories[name] = factory
}

// New creates the configured AI provider.
func New(cfg *config.Config) (Provider, error) {
	providerCfg, err := cfg.ActiveProvider()
	if err != nil {
		return nil, err
	}

	factory, ok := factories[cfg.AI.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s (did you import it?)", cfg.AI.Provider)
	}

	return factory(providerCfg), nil
}
