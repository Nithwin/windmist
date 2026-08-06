package gemini

import (
	"testing"

	"github.com/Nithwin/WindMist/internal/ai"
)

func TestTranslateMessages_ThoughtSignatureOnFunctionCallPart(t *testing.T) {
	sig := "thought-sig-abc"
	messages := []ai.Message{
		{
			Role:    ai.RoleAssistant,
			Content: "",
			ToolCalls: []ai.ToolCall{
				{ID: sig, Name: "read_file", Args: map[string]any{"path": "a.go"}},
				{ID: "call_write_1", Name: "write_file", Args: map[string]any{"path": "b.go"}},
			},
		},
		{
			Role: ai.RoleTool,
			ToolResults: []ai.ToolResult{
				{Name: "read_file", Content: "ok"},
				{Name: "write_file", Content: "ok"},
			},
		},
	}

	contents := translateMessages(messages)
	if len(contents) != 2 {
		t.Fatalf("expected 2 contents, got %d", len(contents))
	}

	modelParts := contents[0].Parts
	if len(modelParts) != 2 {
		t.Fatalf("expected 2 functionCall parts, got %d", len(modelParts))
	}
	if modelParts[0].FunctionCall == nil {
		t.Fatal("first part missing functionCall")
	}
	if modelParts[0].ThoughtSignature != sig {
		t.Fatalf("expected thought signature on first functionCall part, got %q", modelParts[0].ThoughtSignature)
	}
	if modelParts[1].ThoughtSignature != "" {
		t.Fatalf("second functionCall must not carry thought signature, got %q", modelParts[1].ThoughtSignature)
	}

	toolParts := contents[1].Parts
	for i, p := range toolParts {
		if p.ThoughtSignature != "" {
			t.Fatalf("functionResponse part %d must not carry thought signature", i)
		}
		if p.FunctionResponse == nil {
			t.Fatalf("part %d missing functionResponse", i)
		}
	}
}

func TestTranslateResponse_OnlyFirstToolCallGetsSignature(t *testing.T) {
	sig := "sig-xyz"
	candidate := Candidate{
		Content: Content{
			Parts: []Part{
				{
					FunctionCall:     &FunctionCall{Name: "a", Args: map[string]any{}},
					ThoughtSignature: sig,
				},
				{
					FunctionCall: &FunctionCall{Name: "b", Args: map[string]any{}},
				},
			},
		},
	}
	resp := translateResponse(candidate, "gemini-test", &GenerateContentResponse{})
	if len(resp.ToolCalls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(resp.ToolCalls))
	}
	if resp.ToolCalls[0].ID != sig {
		t.Fatalf("first tool call ID = %q, want %q", resp.ToolCalls[0].ID, sig)
	}
	if resp.ToolCalls[1].ID != "call_b_1" {
		t.Fatalf("second tool call ID = %q, want call_b_1", resp.ToolCalls[1].ID)
	}
}

func TestTranslateTools_ArrayItemsType(t *testing.T) {
	tools := translateTools([]ai.ToolDefinition{
		{
			Name: "demo",
			Parameters: []ai.ToolParameter{
				{Name: "tags", Type: "array", ItemsType: "number", Required: true},
			},
		},
	})
	if len(tools) != 1 || len(tools[0].FunctionDeclarations) != 1 {
		t.Fatal("unexpected tools shape")
	}
	schema := tools[0].FunctionDeclarations[0].Parameters.Properties["tags"]
	if schema == nil || schema.Items == nil {
		t.Fatal("missing array items schema")
	}
	if schema.Items.Type != "NUMBER" {
		t.Fatalf("items type = %q, want NUMBER", schema.Items.Type)
	}
}
