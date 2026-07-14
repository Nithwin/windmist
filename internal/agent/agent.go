package agent

import (
	"context"

	"github.com/Nithwin/WindMist/internal/agent/prompt"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/tools"
)

// Config configures the behavior of the agent.
type Config struct {
	// MaxTurns is the maximum number of reasoning iterations the agent
	// may perform before terminating the request.
	MaxTurns int
}

// Result contains the final output produced by the agent.
type Result struct {
	// Content is the final response returned to the user.
	Content string
}

// Agent coordinates the language model and the available tools to solve
// software engineering tasks.
type Agent struct {
	provider ai.Provider
	manager  *tools.Manager

	config Config

	systemPrompt string

	// messages contains the current conversation for this run.
	messages []ai.Message
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

	return &Agent{
		provider:     provider,
		manager:      manager,
		config:       config,
		systemPrompt: prompt.Build(),
	}
}

// Run executes a single user request.
func (a *Agent) Run(ctx context.Context, userPrompt string) (*Result, error) {
	return a.runLoop(ctx, userPrompt)
}