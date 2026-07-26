package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
)

type Manager struct {
	servers map[string]*Client
	tools   map[string]ai.ToolDefinition
	mu      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]*Client),
		tools:   make(map[string]ai.ToolDefinition),
	}
}

// StartAll starts all MCP servers defined in the configuration and registers their tools.
func (m *Manager) StartAll(ctx context.Context, cfg *config.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, srvCfg := range cfg.MCPServers {
		if srvCfg.Command == "" {
			continue
		}

		client := NewClient(name, srvCfg.Command, srvCfg.Args, srvCfg.Env)
		if err := client.Start(ctx); err != nil {
			return fmt.Errorf("failed to start MCP server %s: %w", name, err)
		}

		m.servers[name] = client

		// Fetch tools from the server
		res, err := client.Call(ctx, "tools/list", map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("failed to fetch tools from %s: %w", name, err)
		}

		var toolList struct {
			Tools []struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				InputSchema map[string]interface{} `json:"inputSchema"`
			} `json:"tools"`
		}

		if err := json.Unmarshal(res.Result, &toolList); err != nil {
			return fmt.Errorf("failed to parse tools from %s: %w", name, err)
		}

		// Register tools
		for _, t := range toolList.Tools {
			// Prefix the tool name to avoid collisions
			mcpToolName := fmt.Sprintf("mcp_%s_%s", name, t.Name)

			// Extract parameters
			var params []ai.ToolParameter
			if props, ok := t.InputSchema["properties"].(map[string]interface{}); ok {
				for propName, propVal := range props {
					propMap := propVal.(map[string]interface{})
					desc, _ := propMap["description"].(string)
					typ, _ := propMap["type"].(string)

					required := false
					if reqArr, ok := t.InputSchema["required"].([]interface{}); ok {
						for _, req := range reqArr {
							if req.(string) == propName {
								required = true
								break
							}
						}
					}

					params = append(params, ai.ToolParameter{
						Name:        propName,
						Type:        typ,
						Description: desc,
						Required:    required,
					})
				}
			}

			m.tools[mcpToolName] = ai.ToolDefinition{
				Name:        mcpToolName,
				Description: fmt.Sprintf("[%s] %s", name, t.Description),
				Parameters:  params,
			}
		}
	}
	return nil
}

// GetTools returns all tools registered from MCP servers.
func (m *Manager) GetTools() []ai.ToolDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	var list []ai.ToolDefinition
	for _, t := range m.tools {
		list = append(list, t)
	}
	return list
}

// ExecuteTool calls a tool on the appropriate MCP server.
func (m *Manager) ExecuteTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Parse out the server name and original tool name
	// Format: mcp_{serverName}_{toolName}
	parts := len("mcp_")
	if len(toolName) <= parts {
		return nil, fmt.Errorf("invalid MCP tool name: %s", toolName)
	}

	rest := toolName[parts:]
	serverName := ""
	originalToolName := ""

	// Find the server name by checking prefixes
	for name := range m.servers {
		if len(rest) > len(name) && rest[:len(name)] == name && rest[len(name)] == '_' {
			serverName = name
			originalToolName = rest[len(name)+1:]
			break
		}
	}

	if serverName == "" {
		return nil, fmt.Errorf("could not determine MCP server for tool: %s", toolName)
	}

	client, ok := m.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("MCP server %s not found", serverName)
	}

	// Make the call
	res, err := client.Call(ctx, "tools/call", map[string]interface{}{
		"name":      originalToolName,
		"arguments": args,
	})
	if err != nil {
		return nil, err
	}

	var callResult struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}

	if err := json.Unmarshal(res.Result, &callResult); err != nil {
		return nil, fmt.Errorf("failed to parse MCP tool result: %w", err)
	}

	if callResult.IsError {
		if len(callResult.Content) > 0 {
			return nil, fmt.Errorf("MCP tool error: %s", callResult.Content[0].Text)
		}
		return nil, fmt.Errorf("MCP tool execution failed")
	}

	if len(callResult.Content) > 0 {
		return callResult.Content[0].Text, nil
	}

	return "Success", nil
}

// CloseAll shuts down all MCP servers.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, client := range m.servers {
		client.Close()
	}
	m.servers = make(map[string]*Client)
}
