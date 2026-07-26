package agent

import (
	"testing"

	"github.com/Nithwin/WindMist/internal/ai"
)

func TestPruneMessages(t *testing.T) {
	shortHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial prompt"},
		{Role: ai.RoleAssistant, Content: "Step 1"},
		{Role: ai.RoleTool, Content: "Result 1"},
	}
	mem := SlidingWindowMemory{}
	pruned := mem.Prune(shortHistory, 1000)
	if len(pruned) != 3 {
		t.Errorf("expected length 3, got %d", len(pruned))
	}

	longHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial task goal"},
		{Role: ai.RoleAssistant, Content: "Turn 1 Assistant"},
		{Role: ai.RoleTool, Content: "Turn 1 Tool"},
		{Role: ai.RoleAssistant, Content: "Turn 2 Assistant"},
		{Role: ai.RoleTool, Content: "Turn 2 Tool"},
		{Role: ai.RoleAssistant, Content: "Turn 3 Assistant"},
		{Role: ai.RoleTool, Content: "Turn 3 Tool"},
		{Role: ai.RoleAssistant, Content: "Turn 4 Assistant"},
		{Role: ai.RoleTool, Content: "Turn 4 Tool"},
	}

	prunedLong := mem.Prune(longHistory, 50)
	if len(prunedLong) < 2 {
		t.Fatalf("expected at least 2 messages after pruning, got %d", len(prunedLong))
	}

	if prunedLong[0].Content != "Initial task goal" {
		t.Errorf("expected first message to be preserved, got %q", prunedLong[0].Content)
	}

	prunedDangling := mem.Prune(longHistory, 20)
	if len(prunedDangling) > 0 {
		for i := 1; i < len(prunedDangling); i++ {
			if prunedDangling[i].Role == ai.RoleTool {
				if prunedDangling[i-1].Role != ai.RoleAssistant {
					t.Errorf("Dangling tool found! Tool result at index %d has no preceding Assistant msg", i)
				}
			}
		}
	}
}
