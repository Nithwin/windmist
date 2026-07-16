package agent

import (
	"testing"

	"github.com/Nithwin/WindMist/internal/ai"
)

func TestPruneMessages(t *testing.T) {
	// Case 1: History is smaller or equal to maxKeep + 1 -> Should not prune
	shortHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial prompt"},
		{Role: ai.RoleAssistant, Content: "Step 1"},
		{Role: ai.RoleTool, Content: "Result 1"},
	}
	pruned := pruneMessages(shortHistory, 4)
	if len(pruned) != 3 {
		t.Errorf("expected length 3, got %d", len(pruned))
	}

	// Case 2: History is large -> Should keep index 0 + last maxKeep messages
	longHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial task goal"},     // index 0 (MUST BE KEPT)
		{Role: ai.RoleAssistant, Content: "Turn 1 Assistant"}, // dropped
		{Role: ai.RoleTool, Content: "Turn 1 Tool"},           // dropped
		{Role: ai.RoleAssistant, Content: "Turn 2 Assistant"}, // dropped
		{Role: ai.RoleTool, Content: "Turn 2 Tool"},           // dropped
		{Role: ai.RoleAssistant, Content: "Turn 3 Assistant"}, // kept (last 4 starts here)
		{Role: ai.RoleTool, Content: "Turn 3 Tool"},           // kept
		{Role: ai.RoleAssistant, Content: "Turn 4 Assistant"}, // kept
		{Role: ai.RoleTool, Content: "Turn 4 Tool"},           // kept
	}

	prunedLong := pruneMessages(longHistory, 4)
	if len(prunedLong) != 5 { // 1 initial + 4 recent = 5 total
		t.Fatalf("expected 5 messages after pruning, got %d", len(prunedLong))
	}

	if prunedLong[0].Content != "Initial task goal" {
		t.Errorf("expected first message to be preserved, got %q", prunedLong[0].Content)
	}

	if prunedLong[1].Content != "Turn 3 Assistant" {
		t.Errorf("expected second kept message to be 'Turn 3 Assistant', got %q", prunedLong[1].Content)
	}

	if prunedLong[4].Content != "Turn 4 Tool" {
		t.Errorf("expected last kept message to be 'Turn 4 Tool', got %q", prunedLong[4].Content)
	}

	// Verify the preserved slice matches exact expected order
	expectedRoles := []ai.Role{ai.RoleUser, ai.RoleAssistant, ai.RoleTool, ai.RoleAssistant, ai.RoleTool}
	for i, msg := range prunedLong {
		if msg.Role != expectedRoles[i] {
			t.Errorf("at index %d: expected role %s, got %s", i, expectedRoles[i], msg.Role)
		}
	}
}
