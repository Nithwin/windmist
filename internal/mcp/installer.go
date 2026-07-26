package mcp

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
)

// InstallerCatalog holds the top 5 essential MCP servers that WindMist supports out of the box.
var InstallerCatalog = []CatalogEntry{
	{
		ID:          "github",
		Name:        "GitHub",
		Icon:        "🐙",
		Description: "Read private repos, create PRs, and manage issues",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-github"},
		RequiredEnv: []string{"GITHUB_PERSONAL_ACCESS_TOKEN"},
	},
	{
		ID:          "postgres",
		Name:        "PostgreSQL",
		Icon:        "🐘",
		Description: "Query live databases and analyze schemas",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-postgres"},
		RequiredEnv: []string{"POSTGRES_CONNECTION_STRING"},
		EnvPrompt:   map[string]string{"POSTGRES_CONNECTION_STRING": "Enter Postgres DB URL (postgres://user:pass@localhost/db)"},
	},
	{
		ID:          "sqlite",
		Name:        "SQLite",
		Icon:        "🗄️",
		Description: "Query local SQLite database files",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-sqlite"},
		RequiredEnv: []string{"SQLITE_DB_PATH"},
		EnvPrompt:   map[string]string{"SQLITE_DB_PATH": "Enter absolute path to SQLite file (e.g. /tmp/db.sqlite)"},
	},
	{
		ID:          "puppeteer",
		Name:        "Web Browser (Puppeteer)",
		Icon:        "🌐",
		Description: "Allows the AI to open a web browser and navigate visually",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-puppeteer"},
	},
	{
		ID:          "slack",
		Name:        "Slack",
		Icon:        "💬",
		Description: "Read and send messages in your team workspace",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-slack"},
		RequiredEnv: []string{"SLACK_BOT_TOKEN"},
	},
}

type CatalogEntry struct {
	ID          string
	Name        string
	Icon        string
	Description string
	Command     string
	Args        []string
	RequiredEnv []string
	EnvPrompt   map[string]string // Maps env var to a custom prompt
}

// GetCatalogList returns a formatted string list of available servers for the UI
func GetCatalogList() []string {
	var list []string
	for _, entry := range InstallerCatalog {
		list = append(list, fmt.Sprintf("%s %s - %s", entry.Icon, entry.Name, entry.Description))
	}
	return list
}

// GetCatalogEntry returns a CatalogEntry by its index in the catalog list.
func GetCatalogEntry(index int) (*CatalogEntry, bool) {
	if index >= 0 && index < len(InstallerCatalog) {
		return &InstallerCatalog[index], true
	}
	return nil, false
}

// Install adds the server to the global configuration and saves it.
func Install(entry *CatalogEntry, envValues map[string]string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.MCPServers == nil {
		cfg.MCPServers = make(map[string]config.MCPServerConfig)
	}

	// Create the configuration for this server
	srvConfig := config.MCPServerConfig{
		Command: entry.Command,
		Args:    entry.Args,
		Env:     envValues,
	}

	// Append the dynamic DB path to the args for some servers like SQLite or Postgres
	// Some MCP servers take the DB path as an argument rather than an env var
	if entry.ID == "sqlite" && envValues["SQLITE_DB_PATH"] != "" {
		srvConfig.Args = append(srvConfig.Args, envValues["SQLITE_DB_PATH"])
		delete(srvConfig.Env, "SQLITE_DB_PATH") // Remove from env if passed as arg
	}

	if entry.ID == "postgres" && envValues["POSTGRES_CONNECTION_STRING"] != "" {
		srvConfig.Args = append(srvConfig.Args, envValues["POSTGRES_CONNECTION_STRING"])
		delete(srvConfig.Env, "POSTGRES_CONNECTION_STRING")
	}

	cfg.MCPServers[entry.ID] = srvConfig

	// Save the config back to disk
	return config.Save(cfg)
}
