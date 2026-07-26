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
		},

		UI: UIConfig{
			Theme: "dark",
		},

		Cache: CacheConfig{
			Enabled: true,
		},
	}
}
