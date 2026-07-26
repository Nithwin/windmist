package kimi

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/providers/openai"
)

func init() {
	ai.Register("kimi", New)
}

// New creates a new Kimi (Moonshot AI) provider instance using the OpenAI-compatible client.
func New(cfg config.ProviderConfig) ai.Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.moonshot.ai/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "moonshot-v1-8k"
	}
	return openai.New(cfg)
}
