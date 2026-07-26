package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Client represents an LSP JSON-RPC client connected via stdio.
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	projectPath string

	nextID  int64
	mu      sync.Mutex
	pending map[int64]chan *JSONRPCMessage

	diagMu      sync.Mutex
	diagnostics map[string][]Diagnostic // URI -> Diagnostics

	idleTimer  *time.Timer
	idleMu     sync.Mutex
	onIdleFunc func()
	idleDur    time.Duration
}

// JSONRPCRequest represents a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCMessage can be a request, response, or notification.
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Diagnostic struct {
	Message  string `json:"message"`
	Severity int    `json:"severity"` // 1: Error, 2: Warning, 3: Info, 4: Hint
	Source   string `json:"source"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// NewClient creates a new LSP client.
func NewClient(command string, args []string, projectPath string) *Client {
	cmd := exec.Command(command, args...)
	cmd.Dir = projectPath

	return &Client{
		cmd:         cmd,
		projectPath: projectPath,
		pending:     make(map[int64]chan *JSONRPCMessage),
		diagnostics: make(map[string][]Diagnostic),
	}
}

// Start launches the LSP server and the read loop.
func (c *Client) Start(ctx context.Context) error {
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	c.stdin = stdin
	c.stdout = stdout

	if err := c.cmd.Start(); err != nil {
		return err
	}

	go c.readLoop()

	// Send initialize request
	type InitParams struct {
		ProcessID int    `json:"processId"`
		RootURI   string `json:"rootUri"`
	}

	_, err = c.Call(ctx, "initialize", InitParams{
		ProcessID: c.cmd.Process.Pid,
		RootURI:   "file://" + c.projectPath,
	})
	if err != nil {
		c.Close()
		return fmt.Errorf("LSP initialization failed: %w", err)
	}

	// Send initialized notification
	_ = c.Notify("initialized", map[string]interface{}{})

	return nil
}

// ResetIdleTimer resets the idle countdown.
func (c *Client) ResetIdleTimer() {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	if c.idleTimer != nil {
		c.idleTimer.Reset(c.idleDur)
	}
}

// OnIdle sets a callback to be called when the client is idle.
func (c *Client) OnIdle(duration time.Duration, callback func()) {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	c.idleDur = duration
	c.onIdleFunc = callback
	c.idleTimer = time.AfterFunc(duration, callback)
}

// Call sends a JSON-RPC request and waits for the response.
func (c *Client) Call(ctx context.Context, method string, params interface{}) (*JSONRPCMessage, error) {
	c.ResetIdleTimer()

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

	ch := make(chan *JSONRPCMessage, 1)
	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	if _, err := c.stdin.Write([]byte(msg)); err != nil {
		return nil, err
	}

	select {
	case res := <-ch:
		if res.Error != nil {
			return nil, fmt.Errorf("RPC Error %d: %s", res.Error.Code, res.Error.Message)
		}
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Notify sends a JSON-RPC notification (no response expected).
func (c *Client) Notify(method string, params interface{}) error {
	c.ResetIdleTimer()

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	_, err = c.stdin.Write([]byte(msg))
	return err
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.stdout)
	for {
		// Read headers
		var contentLength int
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return // EOF or closed
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			if strings.HasPrefix(line, "Content-Length:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					contentLength, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				}
			}
		}

		if contentLength == 0 {
			continue
		}

		// Read body
		body := make([]byte, contentLength)
		if _, err := io.ReadFull(reader, body); err != nil {
			return
		}

		var res JSONRPCMessage
		if err := json.Unmarshal(body, &res); err == nil {
			// If it's a response to a request we made
			if res.ID != 0 {
				c.mu.Lock()
				if ch, ok := c.pending[res.ID]; ok {
					ch <- &res
				}
				c.mu.Unlock()
			} else if res.Method == "textDocument/publishDiagnostics" {
				var params PublishDiagnosticsParams
				if err := json.Unmarshal(res.Params, &params); err == nil {
					c.diagMu.Lock()
					c.diagnostics[params.URI] = params.Diagnostics
					c.diagMu.Unlock()
				}
			}
		}
	}
}

// GetDiagnostics returns the latest collected diagnostics for a given file URI.
func (c *Client) GetDiagnostics(uri string) []Diagnostic {
	c.diagMu.Lock()
	defer c.diagMu.Unlock()

	// Create a copy to avoid race conditions
	if diags, ok := c.diagnostics[uri]; ok {
		cpy := make([]Diagnostic, len(diags))
		copy(cpy, diags)
		return cpy
	}
	return nil
}

// Close gracefully terminates the LSP server.
func (c *Client) Close() {
	c.idleMu.Lock()
	if c.idleTimer != nil {
		c.idleTimer.Stop()
	}
	c.idleMu.Unlock()

	_ = c.Notify("exit", nil)
	if c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
}
