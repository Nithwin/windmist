package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI        AIConfig        `yaml:"ai"`
	Providers ProvidersConfig `yaml:"providers"`
	UI        UIConfig        `yaml:"ui"`
	Cache     CacheConfig     `yaml:"cache"`
}

type ProviderConfig struct {
	Model   string `yaml:"model"`
	APIKey  string `yaml:"api_key,omitempty"`
	BaseURL string `yaml:"base_url,omitempty"`
}

// AIConfig stores settings for AI providers and models.
type AIConfig struct {
	Provider string `yaml:"provider"`
}

type ProvidersConfig struct {
	Gemini ProviderConfig `yaml:"gemini"`
	Groq   ProviderConfig `yaml:"groq"`
	Ollama ProviderConfig `yaml:"ollama"`
}

// UIConfig stores user interface preferences.
type UIConfig struct {
	Theme string `yaml:"theme"`
}

// CacheConfig stores cache-related settings.
type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
}
