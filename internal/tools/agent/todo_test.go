package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestTodoTool(t *testing.T) {
	tool := NewTodoTool()

	// Test Add
	res := tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"action": "add",
			"task":   "fix tests",
		},
	})
	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}
	out := res.Output.(string)
	if !strings.Contains(out, "[ ] fix tests") {
		t.Fatalf("expected output to contain task, got: %s", out)
	}

	// Test Complete
	res = tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"action": "complete",
			"task":   "fix tests",
		},
	})
	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}
	out = res.Output.(string)
	if !strings.Contains(out, "[x] fix tests") {
		t.Fatalf("expected output to contain completed task, got: %s", out)
	}

	// Test Remove
	res = tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"action": "remove",
			"task":   "fix tests",
		},
	})
	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}
	out = res.Output.(string)
	if strings.Contains(out, "fix tests") {
		t.Fatalf("expected task to be removed, got: %s", out)
	}
}
