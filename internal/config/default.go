package config

// DefaultConfig returns the default WindMist configuration.
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider: "gemini",
		},

		Providers: ProvidersConfig{
			Gemini: CloudProviderConfig{
				Model:  "gemini-2.5-flash",
				APIKey: "",
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
