package ollama

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
)

// translateTools converts ai.ToolDefinition structs into Ollama /v1/chat/completions Tool schemas.
func translateTools(tools []ai.ToolDefinition) []Tool {
	if len(tools) == 0 {
		return nil
	}

	ollamaTools := make([]Tool, 0, len(tools))
	for _, t := range tools {
		properties := make(map[string]*Property)
		required := make([]string, 0)

		for _, p := range t.Parameters {
			schemaType := "string"
			switch strings.ToLower(p.Type) {
			case "string":
				schemaType = "string"
			case "int", "integer":
				schemaType = "integer"
			case "float", "number":
				schemaType = "number"
			case "bool", "boolean":
				schemaType = "boolean"
			case "array":
				schemaType = "array"
			case "object":
				schemaType = "object"
			}

			properties[p.Name] = &Property{
				Type:        schemaType,
				Description: p.Description,
				Enum:        p.Enum,
			}

			if p.Required {
				required = append(required, p.Name)
			}
		}

		ollamaTools = append(ollamaTools, Tool{
			Type: "function",
			Function: Function{
				Name:        t.Name,
				Description: t.Description,
				Parameters: &Schema{
					Type:       "object",
					Properties: properties,
					Required:   required,
				},
			},
		})
	}

	return ollamaTools
}

// translateMessages converts []ai.Message into Ollama /v1/chat/completions []Message format.
func translateMessages(messages []ai.Message) []Message {
	ollamaMsgs := make([]Message, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case ai.RoleSystem:
			ollamaMsgs = append(ollamaMsgs, Message{
				Role:    "system",
				Content: msg.Content,
			})

		case ai.RoleUser:
			ollamaMsgs = append(ollamaMsgs, Message{
				Role:    "user",
				Content: msg.Content,
			})

		case ai.RoleAssistant:
			var toolCalls []ToolCall
			if len(msg.ToolCalls) > 0 {
				toolCalls = make([]ToolCall, 0, len(msg.ToolCalls))
				for _, tc := range msg.ToolCalls {
					argsJSON, _ := json.Marshal(tc.Args)
					toolCalls = append(toolCalls, ToolCall{
						ID:   tc.ID,
						Type: "function",
						Function: FunctionCall{
							Name:      tc.Name,
							Arguments: string(argsJSON),
						},
					})
				}
			}

			ollamaMsgs = append(ollamaMsgs, Message{
				Role:      "assistant",
				Content:   msg.Content,
				ToolCalls: toolCalls,
			})

		case ai.RoleTool:
			for _, res := range msg.ToolResults {
				ollamaMsgs = append(ollamaMsgs, Message{
					Role:       "tool",
					Content:    res.Content,
					ToolCallID: res.ID,
					Name:       res.Name,
				})
			}
		}
	}

	return ollamaMsgs
}

// translateResponse converts an Ollama ChatResponse into an ai.GenerateResponse.
func translateResponse(model string, resp *ChatResponse) (*ai.GenerateResponse, error) {
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("ollama returned no choices")
	}

	choice := resp.Choices[0]
	toolCalls := make([]ai.ToolCall, 0, len(choice.Message.ToolCalls))

	for i, tc := range choice.Message.ToolCalls {
		args := make(map[string]any)
		if tc.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				// If JSON parsing fails, put raw arguments in map under "raw"
				args["raw"] = tc.Function.Arguments
			}
		}

		id := tc.ID
		if id == "" {
			id = fmt.Sprintf("call_%s_%d", tc.Function.Name, i)
		}

		toolCalls = append(toolCalls, ai.ToolCall{
			ID:   id,
			Name: tc.Function.Name,
			Args: args,
		})
	}

	return &ai.GenerateResponse{
		Text:      choice.Message.Content,
		ToolCalls: toolCalls,
		Model:     model,
		Finish:    choice.FinishReason,
		Usage: ai.Usage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}, nil
}
