package system

import (
	"context"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestGitTool(t *testing.T) {
	tool := NewGitTool(func(cmd string) bool { return true }) // Auto approve for tests

	res := tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"command": "version",
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}

	out := res.Output.(string)
	if out == "" {
		t.Fatal("expected non-empty output from git version")
	}
}
