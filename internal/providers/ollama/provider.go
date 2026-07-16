package ollama

import (
	"context"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
)

func init() {
	ai.Register("ollama", New)
}

// Provider implements the ai.Provider interface for Ollama.
type Provider struct {
	client *Client
	model  string
}

// New creates a new Ollama provider instance.
func New(cfg config.ProviderConfig) ai.Provider {
	baseURL := "http://localhost:11434/v1"

	model := cfg.Model
	if model == "" {
		model = "qwen2.5:8b"
	}

	return &Provider{
		client: NewClient(baseURL, model),
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
		// If the Ollama model does not support tool calling, retry without tools
		if strings.Contains(err.Error(), "does not support tools") {
			chatReq.Tools = nil
			chatResp, err = p.client.GenerateContent(ctx, chatReq)
		}
		if err != nil {
			return nil, err
		}
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

	err := p.client.StreamContent(ctx, chatReq, func(resp *StreamResponse) {
		if len(resp.Choices) == 0 {
			return
		}
		delta := resp.Choices[0].Delta
		if delta.Content != "" {
			onChunk(delta.Content)
		}
	})
	if err != nil && strings.Contains(err.Error(), "does not support tools") {
		chatReq.Tools = nil
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
	return err
}
