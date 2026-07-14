package ai

// Role represents the role of a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents a single conversation message.
type Message struct {
	Role    Role   `json:"role"`
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