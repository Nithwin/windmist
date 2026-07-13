package ai

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
)

// New returns an AI provider based on the configured provider.
func New(cfg *config.Config) (Provider, error) {

	switch cfg.AI.Provider {

	case "gemini":
		return NewGemini(cfg.AI.APIKey, cfg.AI.Model), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.AI.Provider)
	}
}