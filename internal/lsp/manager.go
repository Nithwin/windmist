package lsp

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config represents the configuration for an LSP server.
type Config struct {
	Command string
	Args    []string
}

// Manager manages language servers for different file types.
type Manager struct {
	servers map[string]*Client
	mu      sync.Mutex
	configs map[string]Config // extension -> Config mapping
}

// NewManager creates a new LSP manager.
func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]*Client),
		configs: map[string]Config{
			".go": {Command: "gopls", Args: []string{"serve"}},
			".py": {Command: "pyright-langserver", Args: []string{"--stdio"}},
			".ts": {Command: "typescript-language-server", Args: []string{"--stdio"}},
			".js": {Command: "typescript-language-server", Args: []string{"--stdio"}},
			".rs": {Command: "rust-analyzer", Args: []string{}},
		},
	}
}

// GetClient returns a running client for the file extension, or starts one if not running.
func (m *Manager) GetClient(ctx context.Context, projectPath string, filePath string) (*Client, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	m.mu.Lock()
	defer m.mu.Unlock()

	// If already running, return it and reset its idle timer
	if client, ok := m.servers[ext]; ok {
		client.ResetIdleTimer()
		return client, nil
	}

	// Lookup config
	cfg, ok := m.configs[ext]
	if !ok {
		return nil, fmt.Errorf("no LSP configured for extension %s", ext)
	}

	// Start new client
	client := NewClient(cfg.Command, cfg.Args, projectPath)
	
	// Add an idle callback to automatically shut down the LSP to save RAM
	client.OnIdle(30*time.Second, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if c, exists := m.servers[ext]; exists && c == client {
			c.Close()
			delete(m.servers, ext)
		}
	})

	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start LSP for %s: %w", ext, err)
	}

	m.servers[ext] = client
	return client, nil
}

// CloseAll shuts down all running LSP servers.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for ext, client := range m.servers {
		client.Close()
		delete(m.servers, ext)
	}
}
