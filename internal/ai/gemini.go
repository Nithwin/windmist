package ai

import (
	"context"

	"github.com/Nithwin/WindMist/internal/config"
)

// Gemini implements the Provider interface.
type Gemini struct {
	apiKey string
	model  string
}

// NewGemini creates a new Gemini provider.
func NewGemini(cfg config.ProviderConfig) Provider {
	return &Gemini{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Generate generates a response using the Gemini provider.
func (g *Gemini) Generate(_ context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	prompt := ""

	if req.System != "" {
		prompt += req.System + "\n\n"
	}

	for _, msg := range req.Messages {
		prompt += msg.Role + ": " + msg.Content + "\n"
	}

	return &GenerateResponse{
		Text:   "Gemini: " + prompt,
		Model:  g.model,
		Finish: "stop",
		Usage: Usage{
			InputTokens:  0,
			OutputTokens: 0,
			TotalTokens:  0,
		},
	}, nil
}