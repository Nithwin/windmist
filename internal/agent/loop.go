package agent

import (
	"context"
	"fmt"
	"os"
	"strings"
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

	effectiveMode := a.config.Mode
	if effectiveMode == string(ModeAuto) {
		resolvedMode := a.orchestrateMode(ctx, userPrompt)
		effectiveMode = resolvedMode
		if onChunk != nil {
			onChunk(fmt.Sprintf("\n> 🤖 **Auto-Router**: Selected `%s` mode.\n\n", resolvedMode))
		}
	}

	for turn := 0; turn < a.config.MaxTurns; turn++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		prunedHistory := pruneMessages(messages, a.config.MaxContextTokens)
		// Build dynamic system prompt based on mode
		cwd, _ := os.Getwd()
		modeConfig := GetModeConfig(Mode(effectiveMode))
		dynamicSystemPrompt := prompt.Build(cwd, modeConfig.SystemPrompt)

		req := &ai.GenerateRequest{
			System:   dynamicSystemPrompt,
			Messages: prunedHistory,
			Tools:    a.toolDefinitions(modeConfig),
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

// orchestrateMode sends a fast prompt to the LLM to classify the user's intent.
func (a *Agent) orchestrateMode(ctx context.Context, userPrompt string) string {
	systemPrompt := `You are an AI orchestrator for a coding assistant. 
Your ONLY job is to classify the user's prompt into one of two modes:
1. "build" - The user wants you to write code, fix a bug, create a file, or modify the codebase.
2. "plan" - The user just wants you to analyze, review, explain, search, or output a step-by-step plan WITHOUT modifying any files.

Reply with EXACTLY ONE WORD: either "build" or "plan". Do not include any punctuation or extra text.`

	req := &ai.GenerateRequest{
		System: systemPrompt,
		Messages: []ai.Message{
			{Role: ai.RoleUser, Content: userPrompt},
		},
		// No tools needed for classification
	}

	// We use the same provider to classify, but ideally a smaller model.
	// For now, we just use the active model.
	resp, err := a.provider.Generate(ctx, req)
	if err != nil {
		return string(ModeBuild) // Default fallback
	}

	res := strings.ToLower(strings.TrimSpace(resp.Text))
	if strings.Contains(res, "plan") {
		return string(ModePlan)
	}
	return string(ModeBuild)
}
