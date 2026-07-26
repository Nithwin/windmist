package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI           AIConfig                  `yaml:"ai"`
	Providers    map[string]ProviderConfig `yaml:"providers"`
	UI           UIConfig                  `yaml:"ui"`
	Cache        CacheConfig               `yaml:"cache"`
	SubAgent     SubAgentConfig            `yaml:"subagent,omitempty"`
	CustomModels map[string][]string       `yaml:"custom_models,omitempty"`
}

// SubAgentConfig stores the provider and model for sub-agents.
type SubAgentConfig struct {
	Provider string `yaml:"provider,omitempty"`
	Model    string `yaml:"model,omitempty"`
}

// AIConfig stores the active AI provider.
type AIConfig struct {
	Provider string `yaml:"provider"`
}

// ProviderConfig stores configuration for a single AI provider.
type ProviderConfig struct {
	Model   string `yaml:"model"`
	APIKey  string `yaml:"api_key,omitempty"`
	BaseURL string `yaml:"base_url,omitempty"`
}

// UIConfig stores UI-related settings.
type UIConfig struct {
	Theme string `yaml:"theme"`
}

// CacheConfig stores cache-related settings.
type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
}
