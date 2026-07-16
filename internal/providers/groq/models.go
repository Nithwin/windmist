package groq

// ChatRequest represents a request to Groq's /openai/v1/chat/completions endpoint.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`
}

// Message represents a single conversation message.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// Tool represents a tool definition in OpenAI-compatible format for Groq.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function holds the definition details inside a Tool.
type Function struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Parameters  *Schema `json:"parameters,omitempty"`
}

// Schema represents JSON Schema properties for function parameters.
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

// ToolCall represents a tool call requested by the model.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall holds the tool name and JSON string arguments.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatResponse represents the non-streaming JSON response from Groq.
type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single candidate inside a ChatResponse.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage holds token count metrics returned by Groq.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamResponse represents a Server-Sent Event (SSE) chunk from Groq.
type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a single candidate inside a StreamResponse chunk.
type StreamChoice struct {
	Delta        StreamDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

// StreamDelta holds the incremental content updates.
type StreamDelta struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}
