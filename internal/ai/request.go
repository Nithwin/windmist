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
	Parts       []Part       `json:"parts,omitempty"`
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

// PartType defines the type of content in a multipart message.
type PartType string

const (
	PartText  PartType = "text"
	PartImage PartType = "image" // base64 encoded image
)

// Part represents a segment of a multi-modal message.
type Part struct {
	Type     PartType `json:"type"`
	Text     string   `json:"text,omitempty"`
	MIMEType string   `json:"mime_type,omitempty"`
	Data     string   `json:"data,omitempty"`
}
