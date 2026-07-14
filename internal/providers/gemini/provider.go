package gemini

import (
	"context"
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
)

func init() {
	ai.Register("gemini", New)
}

// Provider implements the ai.Provider interface.
type Provider struct {
	client *Client
	model  string
}

// New creates a new Gemini provider.
func New(cfg config.ProviderConfig) ai.Provider {
	return &Provider{
		client: NewClient(cfg.APIKey, cfg.Model),
		model:  cfg.Model,
	}
}

// Generate generates content using the Gemini API.
func (p *Provider) Generate(
	ctx context.Context,
	req *ai.GenerateRequest,
) (*ai.GenerateResponse, error) {

	geminiReq := &GenerateContentRequest{
		Contents: translateMessages(req.Messages),
		Tools:    translateTools(req.Tools),
	}

	if req.System != "" {
		geminiReq.SystemInstruction = &SystemInstruction{
			Parts: []Part{
				{
					Text: req.System,
				},
			},
		}
	}

	if req.Temperature != 0 || req.MaxTokens != 0 {
		geminiReq.GenerationConfig = &GenerationConfig{
			Temperature:     req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		}
	}

	geminiResp, err := p.client.GenerateContent(ctx, geminiReq)
	if err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("gemini returned no candidates")
	}

	candidate := geminiResp.Candidates[0]

	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned empty response")
	}

	return translateResponse(candidate, p.model, geminiResp), nil
}

// Stream streams a response from Gemini.
func (p *Provider) Stream(
	ctx context.Context,
	req *ai.GenerateRequest,
	onChunk func(string),
) error {

	geminiReq := &GenerateContentRequest{
		Contents: translateMessages(req.Messages),
		Tools:    translateTools(req.Tools),
	}

	if req.System != "" {
		geminiReq.SystemInstruction = &SystemInstruction{
			Parts: []Part{
				{
					Text: req.System,
				},
			},
		}
	}

	if req.Temperature != 0 || req.MaxTokens != 0 {
		geminiReq.GenerationConfig = &GenerationConfig{
			Temperature:     req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		}
	}

	return p.client.StreamContent(
		ctx,
		geminiReq,
		func(resp *GenerateContentResponse) {

			if len(resp.Candidates) == 0 {
				return
			}

			candidate := resp.Candidates[0]

			if len(candidate.Content.Parts) == 0 {
				return
			}

			onChunk(candidate.Content.Parts[0].Text)
		},
	)
}
