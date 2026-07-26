package system

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
)

type GitTool struct {
	approvalCb ApprovalCallback
}

func NewGitTool(approvalCb ApprovalCallback) *GitTool {
	return &GitTool{approvalCb: approvalCb}
}

func (t *GitTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "git",
		Description: "Execute git operations. Safe read-only commands (status, log, diff, branch) auto-run. Write commands (commit, checkout, stash, push) require user approval.",
		Category:    tools.CategoryGit,
		Permission:  tools.PermDangerous,
		Parameters: []tools.Parameter{
			{
				Name:        "command",
				Type:        "string",
				Description: "The git subcommand to run (e.g. status, diff, log -n 5, commit -m 'msg').",
				Required:    true,
			},
		},
	}
}

func (t *GitTool) Run(ctx context.Context, call tools.Call) tools.Result {
	cmdStr, ok := call.Args["command"].(string)
	if !ok || cmdStr == "" {
		return tools.Result{Error: fmt.Errorf("command is required")}
	}

	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return tools.Result{Error: fmt.Errorf("empty command")}
	}

	subcommand := args[0]
	isReadOnly := false
	switch subcommand {
	case "status", "diff", "log", "show", "branch", "rev-parse", "ls-files":
		isReadOnly = true
	}

	if !isReadOnly && t.approvalCb != nil {
		approved := t.approvalCb("git " + cmdStr)
		if !approved {
			return tools.Result{Error: fmt.Errorf("user denied execution of git %s", cmdStr)}
		}
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return tools.Result{Error: fmt.Errorf("git %s failed: %v\nOutput: %s", cmdStr, err, string(out))}
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		output = "(Success: no output)"
	}

	return tools.Result{Output: output}
}
