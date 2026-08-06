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
	cwd, _ := os.Getwd()
	injector := NewReferenceInjector(cwd)
	processedPrompt := injector.Inject(userPrompt)

	messages = appendUser(messages, processedPrompt)
	a.saveMessage(messages[len(messages)-1])

	var totalUsage ai.Usage

	effectiveMode := a.config.Mode
	if effectiveMode == string(ModeAuto) {
		resolvedMode, viaLLM := a.resolveMode(ctx, userPrompt)
		effectiveMode = resolvedMode
		// Only announce when we spent an LLM call or selected an engineering mode.
		if onChunk != nil && (viaLLM || resolvedMode != string(ModeChat)) {
			onChunk(fmt.Sprintf("\n> 🤖 **Auto-Router**: Selected `%s` mode.\n\n", resolvedMode))
		}
	}

	for turn := 0; turn < a.config.MaxTurns; turn++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		prunedHistory := a.config.Memory.Prune(messages, a.config.MaxContextTokens)
		cwd, _ = os.Getwd()
		modeConfig := GetModeConfig(Mode(effectiveMode))
		dynamicSystemPrompt := prompt.Build(cwd, modeConfig.SystemPrompt, prompt.Options{
			IncludeRepoMap: modeConfig.IncludeRepoMap,
			IncludeGuides:  modeConfig.AllowTools,
		})

		req := &ai.GenerateRequest{
			System:   dynamicSystemPrompt,
			Messages: prunedHistory,
			Tools:    a.toolDefinitions(modeConfig),
		}

		if err := a.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit exceeded or context cancelled: %w", err)
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

			if ctx.Err() != nil {
				break
			}

			if attempt == maxRetries {
				break
			}

			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
			case <-timer.C:
			}
			if ctx.Err() != nil {
				break
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

		calls := resp.ToolCalls
		if len(calls) > MaxToolCallsPerTurn {
			if onChunk != nil {
				onChunk(fmt.Sprintf("\n> ⚠️ **Warning**: Truncated tool calls from %d to %d (max per turn).\n\n", len(calls), MaxToolCallsPerTurn))
			}
			calls = calls[:MaxToolCallsPerTurn]
		}

		results := a.execute(ctx, calls, onChunk)
		messages = appendToolResults(messages, results)
		a.saveMessage(messages[len(messages)-1])
	}

	return nil, ErrMaxTurnsExceeded
}

// resolveMode picks chat/plan/build. Prefer local heuristics to save free-tier quota.
// viaLLM is true only when the LLM classifier was actually called.
func (a *Agent) resolveMode(ctx context.Context, userPrompt string) (mode string, viaLLM bool) {
	if m, ok := ClassifyIntent(userPrompt); ok {
		return string(m), false
	}
	return a.orchestrateMode(ctx, userPrompt), true
}

// orchestrateMode sends a fast prompt to the LLM to classify the user's intent.
func (a *Agent) orchestrateMode(ctx context.Context, userPrompt string) string {
	systemPrompt := `Classify the user message into exactly one mode:
- chat: greeting, thanks, small talk, or a short question that needs no codebase tools
- plan: analyze/explain/search/review the project without editing files
- build: write, edit, fix, create, or otherwise change files/code

Reply with ONE word only: chat, plan, or build.`

	req := &ai.GenerateRequest{
		System: systemPrompt,
		Messages: []ai.Message{
			{Role: ai.RoleUser, Content: userPrompt},
		},
	}

	resp, err := a.provider.Generate(ctx, req)
	if err != nil {
		return string(ModeChat) // Prefer cheap fallback on free tier
	}

	res := strings.ToLower(strings.TrimSpace(resp.Text))
	switch {
	case strings.Contains(res, "build"):
		return string(ModeBuild)
	case strings.Contains(res, "plan"):
		return string(ModePlan)
	default:
		return string(ModeChat)
	}
}
