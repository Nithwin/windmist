package gemini

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// StreamContent streams content from the Gemini API.
func (c *Client) StreamContent(
	ctx context.Context,
	req *GenerateContentRequest,
	onChunk func(*GenerateContentResponse),
) error {

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	actualModel := c.model

	// Handle users who have cached -preview models in their config
	if actualModel == "gemini-3.5-flash-preview" {
		actualModel = "gemini-3.5-flash"
	}
	if actualModel == "gemini-3.6-flash-preview" {
		actualModel = "gemini-3.6-flash"
	}

	if actualModel == "gemini-3.5-lite" {
		actualModel = "gemini-3.5-flash-lite"
	}

	// 3.1-pro requires -preview
	if actualModel == "gemini-3.1-pro" {
		actualModel = "gemini-3.1-pro-preview"
	}

	endpoint := fmt.Sprintf(
		"%s/models/%s:streamGenerateContent?alt=sse&key=%s",
		baseURLBeta,
		actualModel,
		url.QueryEscape(c.apiKey),
	)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	maxRetries := 3
	baseDelay := 2 * time.Second

	var resp *http.Response

	for attempt := 0; attempt < maxRetries; attempt++ {
		httpReq.Body = io.NopCloser(bytes.NewReader(body))

		enforceRateLimit(actualModel)

		resp, err = c.client.Do(httpReq)
		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("send request: %w", err)
			}
			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if attempt == maxRetries-1 {
				var apiErr ErrorResponse
				if err := json.Unmarshal(data, &apiErr); err == nil {
					return fmt.Errorf("gemini api (%d): %s (after retries)", apiErr.Error.Code, apiErr.Error.Message)
				}
				return fmt.Errorf("gemini api returned status %d after retries: %s", resp.StatusCode, string(data))
			}

			// Sleep before retry
			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		var apiErr ErrorResponse
		if err := json.Unmarshal(data, &apiErr); err == nil {
			return fmt.Errorf("gemini api (%d): %s", apiErr.Error.Code, apiErr.Error.Message)
		}
		return fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(data))
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		payload := strings.TrimPrefix(line, "data: ")

		var chunk GenerateContentResponse

		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			continue
		}

		onChunk(&chunk)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	return nil
}
