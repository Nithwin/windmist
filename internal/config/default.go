package config

// DefaultConfig returns the default WindMist configuration.
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider: "gemini",
		},

		Providers: ProvidersConfig{
			Gemini: ProviderConfig{
				Model:  "gemini-2.5-flash",
				APIKey: "",
			},

			Groq: ProviderConfig{
				Model:  "llama-3.3-70b-versatile",
				APIKey: "",
			},

			Ollama: ProviderConfig{
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