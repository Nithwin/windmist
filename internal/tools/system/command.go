package system

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Nithwin/WindMist/internal/tools"
)

type streamWriter struct {
	callback func(string)
	buf      []byte
}

func (s *streamWriter) Write(p []byte) (n int, err error) {
	if s.callback != nil {
		s.callback(string(p))
	}
	s.buf = append(s.buf, p...)
	return len(p), nil
}

// ApprovalCallback is a function that asks the user for permission.
type ApprovalCallback func(cmd string) bool

// CommandTool allows the agent to execute terminal commands.
type CommandTool struct {
	AskApproval ApprovalCallback
}

// NewCommandTool creates a new instance of the CommandTool.
func NewCommandTool(cb ApprovalCallback) *CommandTool {
	return &CommandTool{AskApproval: cb}
}

// Definition returns the schema for this tool.
func (t *CommandTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "run_command",
		Description: "Execute a bash command in the terminal. Use this to run tests, compile code, execute git commands, or check system state. Commands have a default timeout of 120 seconds.",
		Category:    tools.CategorySystem,
		Permission:  tools.PermDangerous,
		Parameters: []tools.Parameter{
			{
				Name:        "command",
				Type:        "string",
				Description: "The shell command to execute",
				Required:    true,
			},
			{
				Name:        "timeout_seconds",
				Type:        "integer",
				Description: "Optional timeout in seconds (default: 120, max: 600). Use higher values for long-running builds or test suites.",
				Required:    false,
			},
		},
	}
}

// Run executes the bash command and returns the combined stdout and stderr.
func (t *CommandTool) Run(ctx context.Context, call tools.Call) tools.Result {
	cmdStr, ok := call.Args["command"].(string)
	if !ok || cmdStr == "" {
		return tools.Result{Error: fmt.Errorf("invalid or missing 'command' parameter")}
	}

	if t.AskApproval != nil {
		approved := t.AskApproval(cmdStr)
		if !approved {
			return tools.Result{Error: fmt.Errorf("user rejected execution of command")}
		}
	}

	// Parse optional timeout (default 120s, max 600s)
	timeout := 120 * time.Second
	if t, ok := call.Args["timeout_seconds"]; ok {
		if ts, ok := t.(float64); ok && ts > 0 {
			timeout = time.Duration(ts) * time.Second
			if timeout > 600*time.Second {
				timeout = 600 * time.Second
			}
		}
	}

	// Apply a timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)

	sw := &streamWriter{callback: call.OnChunk}
	cmd.Stdout = sw
	cmd.Stderr = sw

	err := cmd.Run()
	output := sw.buf

	// If it timed out, append a notice
	if ctx.Err() == context.DeadlineExceeded {
		return tools.Result{
			Output: string(output) + fmt.Sprintf("\nCommand timed out after %d seconds.", int(timeout.Seconds())),
		}
	}

	if err != nil {
		return tools.Result{
			Output: fmt.Sprintf("Command failed with error: %v\nOutput: %s", err, string(output)),
		}
	}

	// If output is empty but successful
	if len(output) == 0 {
		return tools.Result{Output: "Command executed successfully with no output."}
	}

	return tools.Result{Output: string(output)}
}
