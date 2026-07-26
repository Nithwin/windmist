package agent

import (
	"context"
	"os"
	"time"

	"github.com/Nithwin/WindMist/internal/agent/prompt"
	"github.com/Nithwin/WindMist/internal/ai"
)

// runLoop executes the iterative reasoning and tool execution loop for the agent.
func (a *Agent) runLoop(ctx context.Context, messages []ai.Message, userPrompt string, onChunk func(string)) (*Result, error) {
	if len(messages) == 0 {
		messages = appendUser(messages, userPrompt)
		a.saveMessage(messages[len(messages)-1])
	}

	var totalUsage ai.Usage

	for turn := 0; turn < a.config.MaxTurns; turn++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		prunedHistory := pruneMessages(messages, a.config.MaxContextTokens)
		// Build dynamic system prompt
		cwd, _ := os.Getwd()
		dynamicSystemPrompt := prompt.Build(cwd)

		req := &ai.GenerateRequest{
			System:   dynamicSystemPrompt,
			Messages: prunedHistory,
			Tools:    a.toolDefinitions(),
		}

		var resp *ai.GenerateResponse
		var err error

		maxRetries := 3
		backoff := 1 * time.Second

		for attempt := 0; attempt <= maxRetries; attempt++ {
			resp, err = a.provider.Stream(ctx, req, onChunk)
			if err == nil {
				break
			}

			// Don't retry if context is cancelled by user
			if ctx.Err() != nil {
				break
			}

			if attempt == maxRetries {
				break
			}

			// Wait before retrying
			select {
			case <-ctx.Done():
				break
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		if err != nil {
			return nil, err
		}

		totalUsage.InputTokens += resp.Usage.InputTokens
		totalUsage.OutputTokens += resp.Usage.OutputTokens
		totalUsage.TotalTokens += resp.Usage.TotalTokens

		messages = appendAssistant(messages, resp.Text, resp.ToolCalls)
		a.saveMessage(messages[len(messages)-1])

		if len(resp.ToolCalls) == 0 {
			return &Result{
				Content: resp.Text,
				Usage:   totalUsage,
				Turns:   turn + 1,
			}, nil
		}

		results := a.execute(ctx, resp.ToolCalls, onChunk)
		messages = appendToolResults(messages, results)
		a.saveMessage(messages[len(messages)-1])
	}

	return nil, ErrMaxTurnsExceeded
}
