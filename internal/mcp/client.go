package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Client struct {
	Name   string
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	nextID  int64
	mu      sync.Mutex
	pending map[int64]chan *JSONRPCResponse
}

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(name, command string, args []string, env map[string]string) *Client {
	cmd := exec.Command(command, args...)

	if len(env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return &Client{
		Name:    name,
		cmd:     cmd,
		pending: make(map[int64]chan *JSONRPCResponse),
	}
}

func (c *Client) Start(ctx context.Context) error {
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// We might also want to pipe stderr for debugging
	c.cmd.Stderr = os.Stderr

	c.stdin = stdin
	c.stdout = stdout

	if err := c.cmd.Start(); err != nil {
		return err
	}

	go c.readLoop()

	// Initialize MCP session
	type ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	type InitParams struct {
		ProtocolVersion string                 `json:"protocolVersion"`
		Capabilities    map[string]interface{} `json:"capabilities"`
		ClientInfo      ClientInfo             `json:"clientInfo"`
	}

	_, err = c.Call(ctx, "initialize", InitParams{
		ProtocolVersion: "2024-11-05", // Standard MCP protocol version
		Capabilities:    map[string]interface{}{},
		ClientInfo: ClientInfo{
			Name:    "WindMist",
			Version: "2.0.0",
		},
	})

	if err != nil {
		c.Close()
		return fmt.Errorf("MCP initialization failed: %w", err)
	}

	// Send initialized notification
	_ = c.Notify("notifications/initialized", map[string]interface{}{})

	return nil
}

func (c *Client) Call(ctx context.Context, method string, params interface{}) (*JSONRPCResponse, error) {
	id := atomic.AddInt64(&c.nextID, 1)
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *JSONRPCResponse, 1)
	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	// MCP usually uses newline-delimited JSON or HTTP-like headers depending on transport.
	// StdIO transport usually uses JSON-RPC directly with \n
	msg := string(data) + "\n"
	if _, err := c.stdin.Write([]byte(msg)); err != nil {
		return nil, err
	}

	select {
	case res := <-ch:
		if res.Error != nil {
			return nil, fmt.Errorf("MCP RPC Error %d: %s", res.Error.Code, res.Error.Message)
		}
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Notify(method string, params interface{}) error {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := string(data) + "\n"
	_, err = c.stdin.Write([]byte(msg))
	return err
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.stdout)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		// Some MCP servers might use Content-Length headers, check for that
		if strings.HasPrefix(string(line), "Content-Length:") {
			parts := strings.Split(string(line), ":")
			if len(parts) == 2 {
				contentLength, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				// read the extra \r\n
				_, _ = reader.ReadBytes('\n')

				body := make([]byte, contentLength)
				if _, err := io.ReadFull(reader, body); err != nil {
					return
				}
				line = body
			}
		}

		var res JSONRPCResponse
		if err := json.Unmarshal(line, &res); err == nil {
			if res.ID != 0 {
				c.mu.Lock()
				if ch, ok := c.pending[res.ID]; ok {
					ch <- &res
				}
				c.mu.Unlock()
			}
		}
	}
}

func (c *Client) Close() {
	if c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
}
