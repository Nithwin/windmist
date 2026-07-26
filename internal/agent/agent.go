package agent

import (
	"context"
	"encoding/json"

	"time"

	"github.com/Nithwin/WindMist/internal/ai"
	appconfig "github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/mcp"
	"github.com/Nithwin/WindMist/internal/store"
	"github.com/Nithwin/WindMist/internal/tools"
)

// Config configures the behavior of the agent.
type Config struct {
	// MaxTurns is the maximum number of reasoning iterations the agent
	// may perform before terminating the request.
	MaxTurns int
	// MaxContextTokens is the maximum number of tokens retained in the
	// sliding window context memory.
	MaxContextTokens int
	// Store is the optional database connection for session persistence.
	Store *store.Store
	// SessionID is the unique identifier for the current session, if persistence is enabled.
	SessionID string
	// Mode is the operating mode of the agent (e.g., build, plan, auto).
	Mode string
	// Memory defines the token pruning strategy.
	Memory MemoryStrategy
}

// Result contains the final output produced by the agent.
type Result struct {
	// Content is the final response returned to the user.
	Content string
	// Usage accumulates token usage across all reasoning turns.
	Usage ai.Usage
	// Turns is the number of turns executed during the request.
	Turns int
}

// Agent coordinates the language model and the available tools to solve
// software engineering tasks.
type Agent struct {
	provider   ai.Provider
	manager    *tools.Manager
	config     Config
	mcpManager *mcp.Manager
}

// New creates a new Agent.
func New(
	provider ai.Provider,
	manager *tools.Manager,
	config Config,
) *Agent {
	if config.MaxTurns <= 0 {
		config.MaxTurns = DefaultMaxTurns
	}
	if config.MaxContextTokens <= 0 {
		config.MaxContextTokens = DefaultMaxContextTokens
	}
	if config.Mode == "" {
		config.Mode = string(ModeBuild)
	}
	if config.Memory == nil {
		config.Memory = SlidingWindowMemory{}
	}

	a := &Agent{
		provider:   provider,
		manager:    manager,
		config:     config,
		mcpManager: mcp.NewManager(),
	}

	// Start MCP servers asynchronously so it doesn't block UI load
	go func() {
		// Create a temporary context for startup
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Load global config to get MCPServers
		globalCfg, err := appconfig.Load()
		if err == nil {
			_ = a.mcpManager.StartAll(ctx, globalCfg)
		}
	}()

	return a
}

// Close gracefully shuts down any resources held by the agent.
func (a *Agent) Close() {
	if a.mcpManager != nil {
		a.mcpManager.CloseAll()
	}
}

// Manager returns the tools manager associated with the agent.
func (a *Agent) Manager() *tools.Manager {
	return a.manager
}

// Run executes a single user request.
func (a *Agent) Run(ctx context.Context, initialMessages []ai.Message, userPrompt string, onChunk func(string)) (*Result, error) {
	if initialMessages == nil {
		initialMessages = make([]ai.Message, 0, 8)
	}
	return a.runLoop(ctx, initialMessages, userPrompt, onChunk)
}

func (a *Agent) saveMessage(msg ai.Message) {
	if a.config.Store == nil || a.config.SessionID == "" {
		return
	}

	sMsg := &store.Message{
		SessionID: a.config.SessionID,
		Role:      string(msg.Role),
		Content:   msg.Content,
	}

	if len(msg.ToolCalls) > 0 {
		b, _ := json.Marshal(msg.ToolCalls)
		sMsg.ToolCalls = string(b)
	}

	if len(msg.ToolResults) > 0 {
		b, _ := json.Marshal(msg.ToolResults)
		sMsg.ToolResults = string(b)
	}

	_ = a.config.Store.SaveMessage(sMsg)
}
