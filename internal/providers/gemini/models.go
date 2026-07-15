package gemini

// GenerateContentRequest represents a Gemini generateContent request.
type GenerateContentRequest struct {
	Contents          []Content          `json:"contents"`
	SystemInstruction *SystemInstruction `json:"systemInstruction,omitempty"`
	GenerationConfig  *GenerationConfig  `json:"generationConfig,omitempty"`
	Tools             []Tool             `json:"tools,omitempty"`
}

// Tool represents Gemini tool definitions.
type Tool struct {
	FunctionDeclarations []FunctionDeclaration `json:"functionDeclarations,omitempty"`
}

// FunctionDeclaration defines a function callable by Gemini.
type FunctionDeclaration struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Parameters  *Schema `json:"parameters,omitempty"`
}

// Schema represents JSON schema for function parameters.
type Schema struct {
	Type        string             `json:"type"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Enum        []string           `json:"enum,omitempty"`
}

// SystemInstruction represents Gemini's system instruction.
type SystemInstruction struct {
	Parts []Part `json:"parts"`
}

// Content represents a conversation message.
type Content struct {
	Role  string `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}

// Part represents a single content part.
type Part struct {
	Text             string            `json:"text,omitempty"`
	FunctionCall     *FunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *FunctionResponse `json:"functionResponse,omitempty"`
}

// FunctionCall represents a request from Gemini to call a function.
type FunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

// FunctionResponse represents tool output returned to Gemini.
type FunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

// GenerationConfig controls generation behavior.
type GenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GenerateContentResponse represents Gemini's response.
type GenerateContentResponse struct {
	Candidates []Candidate `json:"candidates"`
	Usage      Usage       `json:"usageMetadata,omitempty"`
}

// Candidate represents a generated candidate.
type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason,omitempty"`
}

// Usage represents Gemini token usage.
type Usage struct {
	PromptTokenCount     int `json:"promptTokenCount,omitempty"`
	CandidatesTokenCount int `json:"candidatesTokenCount,omitempty"`
	TotalTokenCount      int `json:"totalTokenCount,omitempty"`
}

// ErrorResponse represents a Gemini API error.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents an error returned by Gemini.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
