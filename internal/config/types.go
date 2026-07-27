package config

// Config represents the complete WindMist configuration.
type Config struct {
	AI           AIConfig                   `yaml:"ai"`
	Providers    map[string]ProviderConfig  `yaml:"providers"`
	UI           UIConfig                   `yaml:"ui"`
	Cache        CacheConfig                `yaml:"cache"`
	SubAgent     SubAgentConfig             `yaml:"subagent,omitempty"`
	CustomModels map[string][]string        `yaml:"custom_models,omitempty"`
	MCPServers   map[string]MCPServerConfig `yaml:"mcp_servers,omitempty"`
	Remote       RemoteConfig               `yaml:"remote,omitempty"`
}

// MCPServerConfig stores the configuration for an MCP server.
type MCPServerConfig struct {
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
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

// RemoteConfig stores configurations for remote control integrations.
type RemoteConfig struct {
	Telegram TelegramConfig `yaml:"telegram,omitempty"`
}

// TelegramConfig stores configuration for the Telegram bot.
type TelegramConfig struct {
	Enabled   bool   `yaml:"enabled"`
	BotToken  string `yaml:"bot_token,omitempty"`
	AllowedID string `yaml:"allowed_id,omitempty"` // String type to support large ints safely or usernames
}
