package anthropic

// MessagesRequest represents a request to Anthropic's /v1/messages endpoint.
type MessagesRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	System      string    `json:"system,omitempty"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float32   `json:"temperature,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
	Stream      bool      `json:"stream"`
}

// Message represents a conversation turn (`role: "user"` or `role: "assistant"`).
type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents a text or tool use/result block inside a Message.
type ContentBlock struct {
	Type      string         `json:"type"`                  // "text", "tool_use", or "tool_result"
	Text      string         `json:"text,omitempty"`        // for type "text"
	ID        string         `json:"id,omitempty"`          // for type "tool_use"
	Name      string         `json:"name,omitempty"`        // for type "tool_use"
	Input     map[string]any `json:"input,omitempty"`       // for type "tool_use"
	ToolUseID string         `json:"tool_use_id,omitempty"` // for type "tool_result"
	Content   string         `json:"content,omitempty"`     // for type "tool_result"
}

// Tool represents a tool definition in Anthropic format.
type Tool struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	InputSchema *Schema `json:"input_schema"`
}

// Schema represents JSON Schema properties for Anthropic tools (`type: "object"`).
type Schema struct {
	Type       string               `json:"type"`
	Properties map[string]*Property `json:"properties,omitempty"`
	Required   []string             `json:"required,omitempty"`
}

// Property represents a single parameter definition inside a Schema.
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// MessagesResponse represents the non-streaming JSON response from Anthropic.
type MessagesResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []ContentBlock `json:"content"`
	Model      string         `json:"model"`
	StopReason string         `json:"stop_reason"`
	Usage      Usage          `json:"usage"`
}

// Usage holds token count metrics returned by Anthropic.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// StreamEvent represents an SSE event chunk from Anthropic.
type StreamEvent struct {
	Type         string       `json:"type"`
	Index        int          `json:"index,omitempty"`
	ContentBlock ContentBlock `json:"content_block,omitempty"`
	Delta        StreamDelta  `json:"delta,omitempty"`
	Usage        Usage        `json:"usage,omitempty"`
}

// StreamDelta holds incremental updates during streaming.
type StreamDelta struct {
	Type       string `json:"type"` // e.g. "text_delta"
	Text       string `json:"text,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
}
