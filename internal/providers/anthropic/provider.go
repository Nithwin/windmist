package anthropic

import (
	"context"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
)

func init() {
	ai.Register("anthropic", New)
}

// Provider implements the ai.Provider interface for Anthropic.
type Provider struct {
	client *Client
	model  string
}

// New creates a new Anthropic provider instance.
func New(cfg config.ProviderConfig) ai.Provider {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL = baseURL + "/v1"
	}

	model := cfg.Model
	if model == "" {
		model = "claude-3-5-sonnet-latest"
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

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096 // Anthropic requires max_tokens to be set > 0
	}

	messagesReq := &MessagesRequest{
		Model:       p.model,
		Messages:    translateMessages(req.Messages),
		System:      req.System,
		MaxTokens:   maxTokens,
		Temperature: req.Temperature,
		Tools:       translateTools(req.Tools),
		Stream:      false,
	}

	messagesResp, err := p.client.GenerateContent(ctx, messagesReq)
	if err != nil {
		return nil, err
	}

	return translateResponse(p.model, messagesResp)
}

// Stream streams a completion response chunk by chunk via Client.
func (p *Provider) Stream(
	ctx context.Context,
	req *ai.GenerateRequest,
	onChunk func(string),
) error {

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	messagesReq := &MessagesRequest{
		Model:       p.model,
		Messages:    translateMessages(req.Messages),
		System:      req.System,
		MaxTokens:   maxTokens,
		Temperature: req.Temperature,
		Tools:       translateTools(req.Tools),
		Stream:      true,
	}

	return p.client.StreamContent(ctx, messagesReq, func(event *StreamEvent) {
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" && event.Delta.Text != "" {
			onChunk(event.Delta.Text)
		}
	})
}
