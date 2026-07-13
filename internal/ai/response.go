package ai

// GenerateResponse represents a response from an AI provider.
type GenerateResponse struct {
	Text string

	Model string

	InputTokens int

	OutputTokens int

	FinishReason string
}