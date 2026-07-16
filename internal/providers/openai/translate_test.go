package openai

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

	fn := translated[0].Function
	if fn.Name != "create_file" {
		t.Errorf("expected name 'create_file', got '%s'", fn.Name)
	}

	prop, ok := fn.Parameters.Properties["path"]
	if !ok {
		t.Fatalf("expected property 'path' in parameters")
	}
	if prop.Type != "string" {
		t.Errorf("expected property type 'string', got '%s'", prop.Type)
	}
}

func TestTranslateResponse(t *testing.T) {
	resp := &ChatResponse{
		ID: "test-id",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role: "assistant",
					ToolCalls: []ToolCall{
						{
							ID:   "call_1",
							Type: "function",
							Function: FunctionCall{
								Name:      "read_file",
								Arguments: `{"path": "test.txt"}`,
							},
						},
					},
				},
				FinishReason: "tool_calls",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}

	genResp, err := translateResponse("gpt-4o", resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
