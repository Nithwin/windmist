package gemini

// GenerateContentRequest represents a Gemini generateContent request.
type GenerateContentRequest struct {
	Contents          []Content          `json:"contents"`
	SystemInstruction *SystemInstruction `json:"systemInstruction,omitempty"`
	GenerationConfig  *GenerationConfig  `json:"generationConfig,omitempty"`
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
	Text string `json:"text,omitempty"`
}

// GenerationConfig controls generation behavior.
type GenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GenerateContentResponse represents Gemini's response.
type GenerateContentResponse struct {
	Candidates []Candidate `json:"candidates"`
	Usage       Usage       `json:"usageMetadata,omitempty"`
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