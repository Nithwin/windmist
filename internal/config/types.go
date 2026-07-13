package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI    AIConfig    `yaml:"ai"`
	UI    UIConfig    `yaml:"ui"`
	Cache CacheConfig `yaml:"cache"`
}

// AIConfig stores settings for AI providers and models.
type AIConfig struct {
	Provider string `yaml:"provider"`
	Model    string	`yaml:"model"`
	APIKey   string	`yaml:"api_key"`
}

// UIConfig stores user interface preferences.
type UIConfig struct {
	Theme string `yaml:"theme"`
}

// CacheConfig stores cache-related settings.
type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
}
