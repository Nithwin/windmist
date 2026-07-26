package perplexity

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/providers/openai"
)

func init() {
	ai.Register("perplexity", New)
}

// New creates a new Perplexity provider instance using the OpenAI-compatible client.
func New(cfg config.ProviderConfig) ai.Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.perplexity.ai"
	}
	if cfg.Model == "" {
		cfg.Model = "llama-3.1-sonar-large-128k-online"
	}
	return openai.New(cfg)
}
