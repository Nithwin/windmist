package openai

import (
	"context"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
)

func init() {
	ai.Register("openai", New)
}

// Provider implements the ai.Provider interface for OpenAI.
type Provider struct {
	client *Client
	model  string
}

// New creates a new OpenAI provider instance.
func New(cfg config.ProviderConfig) ai.Provider {
	baseURL := "https://api.openai.com/v1"

	model := cfg.Model
	if model == "" {
		model = "gpt-4o"
	}

	return &Provider{
		client: NewClient(cfg.APIKey, baseURL, model),
		model:  model,
	}
}

// Generate sends a non-streaming completion request via Client.
func (p *Provider) Generate(
	ctx context.Context,
	req *ai.GenerateRequest,
) (*ai.GenerateResponse, error) {

	messages := translateMessages(req.Messages)
	if req.System != "" {
		messages = append([]Message{{
			Role:    "system",
			Content: req.System,
		}}, messages...)
	}

	chatReq := &ChatRequest{
		Model:       p.model,
		Messages:    messages,
		Tools:       translateTools(req.Tools),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      false,
	}

	chatResp, err := p.client.GenerateContent(ctx, chatReq)
	if err != nil {
		return nil, err
	}

	return translateResponse(p.model, chatResp)
}

// Stream streams a completion response chunk by chunk via Client.
func (p *Provider) Stream(
	ctx context.Context,
	req *ai.GenerateRequest,
	onChunk func(string),
) error {

	messages := translateMessages(req.Messages)
	if req.System != "" {
		messages = append([]Message{{
			Role:    "system",
			Content: req.System,
		}}, messages...)
	}

	chatReq := &ChatRequest{
		Model:       p.model,
		Messages:    messages,
		Tools:       translateTools(req.Tools),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	}

	return p.client.StreamContent(ctx, chatReq, func(resp *StreamResponse) {
		if len(resp.Choices) == 0 {
			return
		}
		delta := resp.Choices[0].Delta
		if delta.Content != "" {
			onChunk(delta.Content)
		}
	})
}
