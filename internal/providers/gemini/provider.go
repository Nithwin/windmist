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
		Contents: make([]Content, 0, len(req.Messages)),
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

	for _, msg := range req.Messages {
		geminiReq.Contents = append(geminiReq.Contents, Content{
			Role: string(msg.Role),
			Parts: []Part{
				{
					Text: msg.Content,
				},
			},
		})
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

	return &ai.GenerateResponse{
		Text:   candidate.Content.Parts[0].Text,
		Model:  p.model,
		Finish: candidate.FinishReason,
		Usage: ai.Usage{
			InputTokens:  geminiResp.Usage.PromptTokenCount,
			OutputTokens: geminiResp.Usage.CandidatesTokenCount,
			TotalTokens:  geminiResp.Usage.TotalTokenCount,
		},
	}, nil
}