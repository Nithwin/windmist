package ai

// ToolParameter represents a single parameter for a tool.
type ToolParameter struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolDefinition defines the schema of a tool available to the model.
type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  []ToolParameter `json:"parameters,omitempty"`
}

// ToolCall represents a request from the model to invoke a tool.
type ToolCall struct {
	ID   string         `json:"id,omitempty"`
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

// ToolResult represents the output of a tool execution sent back to the model.
type ToolResult struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
	IsError bool   `json:"is_error,omitempty"`
}
