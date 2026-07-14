package gemini

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
)

// translateTools converts ai.ToolDefinitions into Gemini Tool schemas.
func translateTools(tools []ai.ToolDefinition) []Tool {
	if len(tools) == 0 {
		return nil
	}

	funcDecls := make([]FunctionDeclaration, 0, len(tools))
	for _, tool := range tools {
		properties := make(map[string]*Schema)
		required := make([]string, 0)

		for _, p := range tool.Parameters {
			schemaType := "STRING"
			switch strings.ToLower(p.Type) {
			case "string":
				schemaType = "STRING"
			case "int", "integer":
				schemaType = "INTEGER"
			case "float", "number":
				schemaType = "NUMBER"
			case "bool", "boolean":
				schemaType = "BOOLEAN"
			case "array":
				schemaType = "ARRAY"
			case "object":
				schemaType = "OBJECT"
			}

			properties[p.Name] = &Schema{
				Type:        schemaType,
				Description: p.Description,
				Enum:        p.Enum,
			}
			if p.Required {
				required = append(required, p.Name)
			}
		}

		funcDecls = append(funcDecls, FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters: &Schema{
				Type:       "OBJECT",
				Properties: properties,
				Required:   required,
			},
		})
	}

	return []Tool{
		{
			FunctionDeclarations: funcDecls,
		},
	}
}

// translateMessages converts ai.Messages into Gemini Content items.
func translateMessages(messages []ai.Message) []Content {
	contents := make([]Content, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case ai.RoleUser:
			contents = append(contents, Content{
				Role: "user",
				Parts: []Part{
					{
						Text: msg.Content,
					},
				},
			})

		case ai.RoleAssistant:
			parts := make([]Part, 0, 1+len(msg.ToolCalls))
			if msg.Content != "" {
				parts = append(parts, Part{Text: msg.Content})
			}
			for _, call := range msg.ToolCalls {
				parts = append(parts, Part{
					FunctionCall: &FunctionCall{
						Name: call.Name,
						Args: call.Args,
					},
				})
			}
			if len(parts) > 0 {
				contents = append(contents, Content{
					Role:  "model",
					Parts: parts,
				})
			}

		case ai.RoleTool:
			parts := make([]Part, 0, len(msg.ToolResults))
			for _, res := range msg.ToolResults {
				parts = append(parts, Part{
					FunctionResponse: &FunctionResponse{
						Name: res.Name,
						Response: map[string]any{
							"content":  res.Content,
							"is_error": res.IsError,
						},
					},
				})
			}
			if len(parts) > 0 {
				contents = append(contents, Content{
					Role:  "function",
					Parts: parts,
				})
			}
		}
	}

	return contents
}

// translateResponse converts Gemini response into ai.GenerateResponse.
func translateResponse(candidate Candidate, model string, resp *GenerateContentResponse) *ai.GenerateResponse {
	var textBuilder strings.Builder
	toolCalls := make([]ai.ToolCall, 0)

	for i, part := range candidate.Content.Parts {
		if part.Text != "" {
			textBuilder.WriteString(part.Text)
		}
		if part.FunctionCall != nil {
			toolCalls = append(toolCalls, ai.ToolCall{
				ID:   fmt.Sprintf("call_%s_%d", part.FunctionCall.Name, i),
				Name: part.FunctionCall.Name,
				Args: part.FunctionCall.Args,
			})
		}
	}

	return &ai.GenerateResponse{
		Text:      textBuilder.String(),
		ToolCalls: toolCalls,
		Model:     model,
		Finish:    candidate.FinishReason,
		Usage: ai.Usage{
			InputTokens:  resp.Usage.PromptTokenCount,
			OutputTokens: resp.Usage.CandidatesTokenCount,
			TotalTokens:  resp.Usage.TotalTokenCount,
		},
	}
}
