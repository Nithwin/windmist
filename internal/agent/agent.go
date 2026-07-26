package agent

import (
	"context"
	"encoding/json"

	"github.com/Nithwin/WindMist/internal/ai"
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
	provider ai.Provider
	manager  *tools.Manager

	config Config
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

	return &Agent{
		provider: provider,
		manager:  manager,
		config:   config,
	}
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
