package anthropic

import (
	"testing"

	"github.com/Nithwin/WindMist/internal/ai"
)

func TestTranslateTools(t *testing.T) {
	tools := []ai.ToolDefinition{
		{
			Name:        "create_file",
			Description: "Creates a file",
			Parameters: []ai.ToolParameter{
				{
					Name:        "path",
					Type:        "string",
					Description: "Path to file",
					Required:    true,
				},
			},
		},
	}

	translated := translateTools(tools)
	if len(translated) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(translated))
	}

	fn := translated[0]
	if fn.Name != "create_file" {
		t.Errorf("expected name 'create_file', got '%s'", fn.Name)
	}

	prop, ok := fn.InputSchema.Properties["path"]
	if !ok {
		t.Fatalf("expected property 'path' in parameters")
	}
	if prop.Type != "string" {
		t.Errorf("expected property type 'string', got '%s'", prop.Type)
	}
}

func TestTranslateResponse(t *testing.T) {
	resp := &MessagesResponse{
		ID:   "msg_test",
		Type: "message",
		Role: "assistant",
		Content: []ContentBlock{
			{
				Type: "text",
				Text: "I am calling a tool.",
			},
			{
				Type: "tool_use",
				ID:   "call_1",
				Name: "read_file",
				Input: map[string]any{
					"path": "test.txt",
				},
			},
		},
		StopReason: "tool_use",
		Usage: Usage{
			InputTokens:  10,
			OutputTokens: 5,
		},
	}

	genResp, err := translateResponse("claude-3-5-sonnet-latest", resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if genResp.Text != "I am calling a tool." {
		t.Errorf("expected text 'I am calling a tool.', got '%s'", genResp.Text)
	}

	if len(genResp.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(genResp.ToolCalls))
	}

	call := genResp.ToolCalls[0]
	if call.Name != "read_file" {
		t.Errorf("expected tool call name 'read_file', got '%s'", call.Name)
	}
	if path, ok := call.Args["path"].(string); !ok || path != "test.txt" {
		t.Errorf("expected path 'test.txt', got '%v'", call.Args["path"])
	}
}
