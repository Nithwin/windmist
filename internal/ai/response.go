package ai

// Usage contains token usage information returned by the provider.
type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// GenerateResponse contains the provider's response.
type GenerateResponse struct {
	Text      string
	ToolCalls []ToolCall
	Model     string
	Finish    string
	Usage     Usage
}
