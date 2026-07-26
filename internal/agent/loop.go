package agent

import (
	"context"

	"github.com/Nithwin/WindMist/internal/ai"
)

// runLoop executes the iterative reasoning and tool execution loop for the agent.
func (a *Agent) runLoop(ctx context.Context, messages []ai.Message, userPrompt string, onChunk func(string)) (*Result, error) {
	if len(messages) == 0 {
		messages = appendUser(messages, userPrompt)
	}

	var totalUsage ai.Usage

	for turn := 0; turn < a.config.MaxTurns; turn++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		prunedHistory := pruneMessages(messages, a.config.MaxContextTokens)

		req := &ai.GenerateRequest{
			System:   a.systemPrompt,
			Messages: prunedHistory,
			Tools:    a.toolDefinitions(),
		}

		var resp *ai.GenerateResponse
		var err error

		// Use stream only for the first turn (to show the user something is happening)
		// Or always stream. Since we patched providers to return GenerateResponse, we can always Stream!
		resp, err = a.provider.Stream(ctx, req, onChunk)
		if err != nil {
			return nil, err
		}

		totalUsage.InputTokens += resp.Usage.InputTokens
		totalUsage.OutputTokens += resp.Usage.OutputTokens
		totalUsage.TotalTokens += resp.Usage.TotalTokens

		messages = appendAssistant(messages, resp.Text, resp.ToolCalls)

		if len(resp.ToolCalls) == 0 {
			return &Result{
				Content: resp.Text,
				Usage:   totalUsage,
				Turns:   turn + 1,
			}, nil
		}

		results := a.execute(ctx, resp.ToolCalls)
		messages = appendToolResults(messages, results)
	}

	return nil, ErrMaxTurnsExceeded
}
