package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://generativelanguage.googleapis.com/v1beta"
)

// Client handles communication with the Gemini API.
type Client struct {
	apiKey string
	model  string
	client *http.Client
}

// NewClient creates a new Gemini HTTP client.
func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateContent sends a request to the Gemini API.
func (c *Client) GenerateContent(
	ctx context.Context,
	req *GenerateContentRequest,
) (*GenerateContentResponse, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf(
		"%s/models/%s:generateContent?key=%s",
		baseURL,
		c.model,
		url.QueryEscape(c.apiKey),
	)

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
		var apiErr ErrorResponse

		if err := json.Unmarshal(data, &apiErr); err == nil {
			return nil, fmt.Errorf(
				"gemini api (%d): %s",
				apiErr.Error.Code,
				apiErr.Error.Message,
			)
		}

		return nil, fmt.Errorf(
			"gemini api returned status %d: %s",
			resp.StatusCode,
			string(data),
		)
	}

	var result GenerateContentResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
