package anthropic

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
)

// translateTools converts ai.ToolDefinition structs into Anthropic /v1/messages Tool schemas.
func translateTools(tools []ai.ToolDefinition) []Tool {
	if len(tools) == 0 {
		return nil
	}

	anthropicTools := make([]Tool, 0, len(tools))
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

		anthropicTools = append(anthropicTools, Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: &Schema{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
		})
	}

	return anthropicTools
}

// translateMessages converts []ai.Message into Anthropic []Message blocks format.
// Note: System prompts should be extracted separately as Anthropic expects `system` at the top-level request.
func translateMessages(messages []ai.Message) []Message {
	anthropicMsgs := make([]Message, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case ai.RoleSystem:
			// System role is skipped here because Anthropic expects system instruction in req.System directly.
			continue

		case ai.RoleUser:
			anthropicMsgs = append(anthropicMsgs, Message{
				Role: "user",
				Content: []ContentBlock{
					{
						Type: "text",
						Text: msg.Content,
					},
				},
			})

		case ai.RoleAssistant:
			blocks := make([]ContentBlock, 0, 1+len(msg.ToolCalls))
			if msg.Content != "" {
				blocks = append(blocks, ContentBlock{
					Type: "text",
					Text: msg.Content,
				})
			}

			for _, tc := range msg.ToolCalls {
				blocks = append(blocks, ContentBlock{
					Type:  "tool_use",
					ID:    tc.ID,
					Name:  tc.Name,
					Input: tc.Args,
				})
			}

			anthropicMsgs = append(anthropicMsgs, Message{
				Role:    "assistant",
				Content: blocks,
			})

		case ai.RoleTool:
			// Anthropic requires tool results to be inside a user message
			blocks := make([]ContentBlock, 0, len(msg.ToolResults))
			for _, res := range msg.ToolResults {
				blocks = append(blocks, ContentBlock{
					Type:      "tool_result",
					ToolUseID: res.ID,
					Content:   res.Content,
				})
			}

			anthropicMsgs = append(anthropicMsgs, Message{
				Role:    "user",
				Content: blocks,
			})
		}
	}

	return anthropicMsgs
}

// translateResponse converts an Anthropic MessagesResponse into an ai.GenerateResponse.
func translateResponse(model string, resp *MessagesResponse) (*ai.GenerateResponse, error) {
	var textBuilder strings.Builder
	var toolCalls []ai.ToolCall

	for i, block := range resp.Content {
		switch block.Type {
		case "text":
			textBuilder.WriteString(block.Text)
		case "tool_use":
			id := block.ID
			if id == "" {
				id = fmt.Sprintf("call_%s_%d", block.Name, i)
			}
			toolCalls = append(toolCalls, ai.ToolCall{
				ID:   id,
				Name: block.Name,
				Args: block.Input,
			})
		}
	}

	return &ai.GenerateResponse{
		Text:      textBuilder.String(),
		ToolCalls: toolCalls,
		Model:     model,
		Finish:    resp.StopReason,
		Usage: ai.Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}
