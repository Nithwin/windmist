package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI        AIConfig        `yaml:"ai"`
	Providers ProvidersConfig `yaml:"providers"`
	UI        UIConfig        `yaml:"ui"`
	Cache     CacheConfig     `yaml:"cache"`
}

// AIConfig stores settings for AI providers and models.
type AIConfig struct {
	Provider string `yaml:"provider"`
}

// CloudProviderConfig stores configuration for cloud AI providers.
type CloudProviderConfig struct {
	Model  string `yaml:"model"`
	APIKey string `yaml:"api_key"`
}

// OllamaConfig stores Ollama settings.
type OllamaConfig struct {
	Model   string `yaml:"model"`
	BaseURL string `yaml:"base_url"`
}

// UIConfig stores user interface preferences.
type UIConfig struct {
	Theme string `yaml:"theme"`
}

// CacheConfig stores cache-related settings.
type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ProvidersConfig stores configuration for all supported providers.
type ProvidersConfig struct {
	Gemini CloudProviderConfig `yaml:"gemini"`
	Groq   CloudProviderConfig `yaml:"groq"`
	Ollama OllamaConfig        `yaml:"ollama"`
}
