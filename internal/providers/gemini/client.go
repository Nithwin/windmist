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

	maxRetries := 3
	baseDelay := 2 * time.Second

	var resp *http.Response
	var data []byte

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Reset the request body since it gets consumed on each call
		httpReq.Body = io.NopCloser(bytes.NewReader(body))

		resp, err = c.client.Do(httpReq)
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("send request: %w", err)
			}
			// Wait and retry
			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		data, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		// Handle retryable status codes: 429 (Too Many Requests), 503 (Service Unavailable)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if attempt == maxRetries-1 {
				var apiErr ErrorResponse
				if err := json.Unmarshal(data, &apiErr); err == nil {
					return nil, fmt.Errorf(
						"gemini api (%d): %s (after retries)",
						apiErr.Error.Code,
						apiErr.Error.Message,
					)
				}
				return nil, fmt.Errorf("gemini api returned status %d after retries: %s", resp.StatusCode, string(data))
			}

			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		// Non-retryable error
		var apiErr ErrorResponse
		if err := json.Unmarshal(data, &apiErr); err == nil {
			return nil, fmt.Errorf(
				"gemini api (%d): %s",
				apiErr.Error.Code,
				apiErr.Error.Message,
			)
		}
		return nil, fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(data))
	}

	var result GenerateContentResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
