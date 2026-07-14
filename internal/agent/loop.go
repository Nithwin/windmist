package agent

import (
	"context"

	"github.com/Nithwin/WindMist/internal/ai"
)

// runLoop executes the iterative reasoning and tool execution loop for the agent.
func (a *Agent) runLoop(ctx context.Context, userPrompt string) (*Result, error) {
	// Initialize the conversation with the user prompt if empty.
	if len(a.messages) == 0 {
		a.messages = append(a.messages, ai.Message{
			Role:    ai.RoleUser,
			Content: userPrompt,
		})
	}

	for turn := 0; turn < a.config.MaxTurns; turn++ {
		// Check for context cancellation before each turn.
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		req := &ai.GenerateRequest{
			System:   a.systemPrompt,
			Messages: a.messages,
		}

		resp, err := a.provider.Generate(ctx, req)
		if err != nil {
			return nil, err
		}

		// Append the assistant's response to the conversation history.
		a.messages = append(a.messages, ai.Message{
			Role:    ai.RoleAssistant,
			Content: resp.Text,
		})

		// TODO: Implement tool call parsing and execution inside executor.go once tool calling schema is wired up.
		// For now, return the model's text response directly when no further tool actions are required.
		return &Result{
			Content: resp.Text,
		}, nil
	}

	return &Result{
		Content: "Agent reached the maximum turn limit without completing the task.",
	}, nil
}
