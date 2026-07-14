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
		},

		UI: UIConfig{
			Theme: "dark",
		},

		Cache: CacheConfig{
			Enabled: true,
		},
	}
}