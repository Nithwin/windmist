package deepseek

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/providers/openai"
)

func init() {
	ai.Register("deepseek", New)
}

// New creates a new DeepSeek provider instance using the OpenAI-compatible client.
func New(cfg config.ProviderConfig) ai.Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-coder"
	}
	return openai.New(cfg)
}
