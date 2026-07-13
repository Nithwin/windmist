package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI    AIConfig
	UI    UIConfig
	Cache CacheConfig
}

// AIConfig stores settings for AI providers and models.
type AIConfig struct {
	Provider string
	Model    string
	APIKey   string
}

// UIConfig stores user interface preferences.
type UIConfig struct {
	Theme string
}

// CacheConfig stores cache-related settings.
type CacheConfig struct {
	Enabled bool
}