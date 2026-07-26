package together

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/providers/openai"
)

func init() {
	ai.Register("together", New)
}

// New creates a new Together AI provider instance using the OpenAI-compatible client.
func New(cfg config.ProviderConfig) ai.Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.together.xyz/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo"
	}
	return openai.New(cfg)
}
