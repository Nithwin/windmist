package config

// DefaultConfig returns the default WindMist configuration.
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider: "gemini",
		},

		Providers: map[string]ProviderConfig{
			"gemini": {
				Model: "gemini-2.5-flash",
			},
			"groq": {
				Model: "llama-3.3-70b-versatile",
			},
			"ollama": {
				Model:   "qwen3:8b",
				BaseURL: "http://localhost:11434",
			},
			"openai": {
				Model: "gpt-4o",
			},
			"anthropic": {
				Model: "claude-3-5-sonnet-latest",
			},
			"deepseek": {
				Model: "deepseek-chat",
			},
			"mistral": {
				Model: "mistral-large-latest",
			},
			"kimi": {
				Model: "kimi-k3",
			},
			"perplexity": {
				Model: "llama-3.1-sonar-large-128k-online",
			},
			"together": {
				Model: "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo",
			},
		},

		UI: UIConfig{
			Theme: "dark",
		},

		Cache: CacheConfig{
			Enabled: true,
		},
	}
}
