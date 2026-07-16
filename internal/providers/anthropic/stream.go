package anthropic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamContent streams Server-Sent Events (SSE) chunks from the Anthropic API server.
func (c *Client) StreamContent(
	ctx context.Context,
	req *MessagesRequest,
	onChunk func(*StreamEvent),
) error {

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/messages", c.baseURL)

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
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	if c.apiKey != "" {
		httpReq.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("anthropic server returned status %d: %s", resp.StatusCode, string(data))
	}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}

		var chunk StreamEvent
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
