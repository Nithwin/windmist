package agent

import (
	"context"
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
)

// Summarizer handles compressing old conversation context into a dense summary.
type Summarizer struct {
	provider  ai.Provider
	threshold int
}

// NewSummarizer creates a new Summarizer.
// The threshold is the minimum number of tokens before auto-summarization triggers.
func NewSummarizer(provider ai.Provider, threshold int) *Summarizer {
	if threshold <= 0 {
		threshold = 8000 // Default threshold
	}
	return &Summarizer{
		provider:  provider,
		threshold: threshold,
	}
}

// ShouldSummarize checks if the current message history exceeds the token threshold.
func (s *Summarizer) ShouldSummarize(messages []ai.Message) bool {
	if len(messages) <= 3 {
		return false // Too short to summarize meaningfully
	}

	totalTokens := 0
	for _, msg := range messages {
		totalTokens += countMessageTokens(msg)
	}

	return totalTokens > s.threshold
}

// Summarize takes a list of messages and returns a single summary string.
func (s *Summarizer) Summarize(ctx context.Context, messages []ai.Message) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	systemPrompt := `You are an AI assistant tasked with summarizing a conversation between a User and an AI Coding Assistant.
Your goal is to create a dense, highly compressed summary that retains all technical context, decisions, file paths, code snippets, and outstanding tasks.
Do not include conversational filler (e.g., "The user asked for..."). Just provide the facts and the context.`

	req := &ai.GenerateRequest{
		System:   systemPrompt,
		Messages: messages,
	}

	resp, err := s.provider.Generate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return resp.Text, nil
}

// Compact merges older messages into a single summary message, keeping the most recent N messages intact.
// It returns the new, compacted message list.
func (s *Summarizer) Compact(ctx context.Context, messages []ai.Message, keepRecent int) ([]ai.Message, error) {
	if len(messages) <= keepRecent+2 {
		return messages, nil // Nothing to compact
	}

	// Always keep the system message if present
	startIndex := 0
	if messages[0].Role == ai.RoleSystem {
		startIndex = 1
	}

	// We want to keep the last `keepRecent` messages.
	// We summarize from startIndex up to len(messages)-keepRecent.
	endIndex := len(messages) - keepRecent

	if endIndex <= startIndex {
		return messages, nil
	}

	toSummarize := messages[startIndex:endIndex]
	summaryText, err := s.Summarize(ctx, toSummarize)
	if err != nil {
		return nil, err
	}

	// Construct the new history
	compacted := make([]ai.Message, 0, keepRecent+2)
	if startIndex == 1 {
		compacted = append(compacted, messages[0])
	}

	summaryMsg := ai.Message{
		Role:    ai.RoleAssistant,
		Content: fmt.Sprintf("*(Context Summary of earlier conversation)*\n\n%s", summaryText),
	}
	compacted = append(compacted, summaryMsg)

	// Append the recent unsummarized messages
	compacted = append(compacted, messages[endIndex:]...)

	return compacted, nil
}
