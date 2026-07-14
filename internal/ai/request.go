package ai

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateRequest contains everything required to generate a response.
type GenerateRequest struct {
	System      string
	Messages    []Message
	Temperature float32
	MaxTokens   int
	Stream      bool
}