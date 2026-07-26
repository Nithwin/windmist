package mistral

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/providers/openai"
)

func init() {
	ai.Register("mistral", New)
}

// New creates a new Mistral provider instance using the OpenAI-compatible client.
func New(cfg config.ProviderConfig) ai.Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.mistral.ai/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "mistral-large-latest"
	}
	return openai.New(cfg)
}
