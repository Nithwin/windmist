package system

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Nithwin/WindMist/internal/tools"
)

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
		Description: "Execute a bash command in the terminal. Use this to run tests, compile code, execute git commands, or check system state.",
		Parameters: []tools.Parameter{
			{
				Name:        "command",
				Type:        "string",
				Description: "The shell command to execute",
				Required:    true,
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

	// Apply a timeout to prevent hanging commands (e.g., waiting for input)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)

	output, err := cmd.CombinedOutput()

	// If it timed out, append a notice
	if ctx.Err() == context.DeadlineExceeded {
		return tools.Result{
			Output: string(output) + "\nCommand timed out after 30 seconds.",
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
