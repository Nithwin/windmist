package agent

import (
	"testing"

	"github.com/Nithwin/WindMist/internal/ai"
)

func TestPruneMessages(t *testing.T) {
	// Case 1: History is smaller or equal to budget -> Should not prune
	shortHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial prompt"}, // 14/4 = 3 tokens
		{Role: ai.RoleAssistant, Content: "Step 1"},    // 6/4 = 1 token
		{Role: ai.RoleTool, Content: "Result 1"},       // 8/4 = 2 tokens
	}
	pruned := pruneMessages(shortHistory, 100)
	if len(pruned) != 3 {
		t.Errorf("expected length 3, got %d", len(pruned))
	}

	// Case 2: History is large -> Should keep index 0 + last fitting messages
	longHistory := []ai.Message{
		{Role: ai.RoleUser, Content: "Initial task goal"},     // 17/4 = 4 tokens (always kept)
		{Role: ai.RoleAssistant, Content: "Turn 1 Assistant"}, // 16/4 = 4 tokens
		{Role: ai.RoleTool, Content: "Turn 1 Tool"},           // 11/4 = 2 tokens
		{Role: ai.RoleAssistant, Content: "Turn 2 Assistant"}, // 16/4 = 4 tokens
		{Role: ai.RoleTool, Content: "Turn 2 Tool"},           // 11/4 = 2 tokens
		{Role: ai.RoleAssistant, Content: "Turn 3 Assistant"}, // 16/4 = 4 tokens
		{Role: ai.RoleTool, Content: "Turn 3 Tool"},           // 11/4 = 2 tokens
		{Role: ai.RoleAssistant, Content: "Turn 4 Assistant"}, // 16/4 = 4 tokens
		{Role: ai.RoleTool, Content: "Turn 4 Tool"},           // 11/4 = 2 tokens
	}

	// budget = 16 - 4 (first) = 12 tokens
	// Turn 4 = 6 tokens, Turn 3 = 6 tokens. Both fit exactly (12 tokens).
	prunedLong := pruneMessages(longHistory, 16)
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

	// Case 3: Dangling Tool Result
	// budget = 14 - 4 (first) = 10 tokens
	// Turn 4 = 6 tokens. Remaining budget = 4.
	// Turn 3 Tool = 2 tokens. Remaining budget = 2.
	// Turn 3 Assistant = 4 tokens. Exceeds budget (2 < 4)! Break.
	// Keep list starts with "Turn 3 Tool", which is dangling. It should be stripped.
	// Final expected: First message + Turn 4 = 3 messages.
	prunedDangling := pruneMessages(longHistory, 14)
	if len(prunedDangling) != 3 {
		t.Fatalf("expected 3 messages after pruning dangling tool, got %d", len(prunedDangling))
	}
	if prunedDangling[1].Content != "Turn 4 Assistant" {
		t.Errorf("expected dangling tool to be removed and start with 'Turn 4 Assistant', got %q", prunedDangling[1].Content)
	}
}
