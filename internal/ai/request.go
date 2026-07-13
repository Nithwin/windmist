package ai

// GenerateRequest represents a request sent to an AI provider.
type GenerateRequest struct {
	Prompt string

	SystemPrompt string

	Temperature float32

	MaxTokens int
}