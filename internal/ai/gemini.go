package ai

import "context"

// Gemini implements the Provider interface.
type Gemini struct {
	apiKey string
	model  string
}

// NewGemini creates a Gemini provider.
func NewGemini(apiKey, model string) *Gemini {
	return &Gemini{
		apiKey: apiKey,
		model:  model,
	}
}

// Generate generates a response using Gemini.
func (g *Gemini) Generate(
	ctx context.Context,
	req *GenerateRequest,
) (*GenerateResponse, error) {

	return &GenerateResponse{
		Text: "Gemini integration coming next...",
	}, nil
}