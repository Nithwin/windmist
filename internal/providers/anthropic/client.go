package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles low-level HTTP communication with the Anthropic API server.
type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewClient creates a new Anthropic HTTP client.
func NewClient(apiKey, baseURL, model string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateContent sends a non-streaming completion request to the Anthropic server.
func (c *Client) GenerateContent(
	ctx context.Context,
	req *MessagesRequest,
) (*MessagesResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/messages", c.baseURL)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	if c.apiKey != "" {
		httpReq.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic server returned status %d: %s", resp.StatusCode, string(data))
	}

	var result MessagesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
