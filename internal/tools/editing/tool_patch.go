package editing

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
)

type PatchTool struct{}

func NewPatchTool() *PatchTool {
	return &PatchTool{}
}

func (t *PatchTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "patch",
		Description: "Applies a unified diff patch to the workspace. Useful for making complex modifications to files efficiently.",
		Parameters: []tools.Parameter{
			{
				Name:        "diff",
				Type:        "string",
				Description: "The unified diff string to apply.",
				Required:    true,
			},
		},
	}
}

func (t *PatchTool) Run(ctx context.Context, call tools.Call) tools.Result {
	diff, ok := call.Args["diff"].(string)
	if !ok || diff == "" {
		return tools.Result{Error: fmt.Errorf("diff is required")}
	}

	// Create temp file for the patch
	tmpFile, err := os.CreateTemp("", "windmist-patch-*.diff")
	if err != nil {
		return tools.Result{Error: fmt.Errorf("failed to create temp file: %w", err)}
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(diff); err != nil {
		return tools.Result{Error: fmt.Errorf("failed to write patch: %w", err)}
	}
	tmpFile.Close()

	// Try patch command first (standard on linux/mac)
	cmd := exec.CommandContext(ctx, "patch", "-p1", "-i", tmpFile.Name())
	out, err := cmd.CombinedOutput()
	if err == nil {
		return tools.Result{Output: "Patch applied successfully:\n" + string(out)}
	}

	// If patch fails, try git apply
	gitCmd := exec.CommandContext(ctx, "git", "apply", tmpFile.Name())
	gitOut, gitErr := gitCmd.CombinedOutput()
	if gitErr == nil {
		return tools.Result{Output: "Patch applied successfully via git apply:\n" + string(gitOut)}
	}

	return tools.Result{
		Error: fmt.Errorf("failed to apply patch.\npatch error: %v, out: %s\ngit apply error: %v, out: %s",
			err, strings.TrimSpace(string(out)),
			gitErr, strings.TrimSpace(string(gitOut))),
	}
}
