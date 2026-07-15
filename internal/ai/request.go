package ai

// Role represents the role of a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a single conversation message.
type Message struct {
	Role        Role         `json:"role"`
	Content     string       `json:"content"`
	ToolCalls   []ToolCall   `json:"tool_calls,omitempty"`
	ToolResults []ToolResult `json:"tool_results,omitempty"`
}

// GenerateRequest contains everything required to generate a response.
type GenerateRequest struct {
	System      string
	Messages    []Message
	Tools       []ToolDefinition
	Temperature float32
	MaxTokens   int
	Stream      bool
}
