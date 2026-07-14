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

// NewGemini creates a Gemini provider.
func NewGemini(cfg config.ProviderConfig) Provider {
	return &Gemini{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Generate generates a response using Gemini.
func (g *Gemini) Generate(_ context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	return &GenerateResponse{
		Text: "Gemini: " + req.Prompt,
	}, nil
}