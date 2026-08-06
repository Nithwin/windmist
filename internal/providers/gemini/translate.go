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

			var itemsSchema *Schema
			if schemaType == "ARRAY" {
				itemType := "STRING"
				switch strings.ToLower(p.ItemsType) {
				case "int", "integer":
					itemType = "INTEGER"
				case "float", "number":
					itemType = "NUMBER"
				case "bool", "boolean":
					itemType = "BOOLEAN"
				case "object":
					itemType = "OBJECT"
				case "array":
					itemType = "ARRAY"
				case "string", "":
					itemType = "STRING"
				default:
					itemType = "STRING"
				}
				itemsSchema = &Schema{Type: itemType}
			}

			properties[p.Name] = &Schema{
				Type:        schemaType,
				Description: p.Description,
				Enum:        p.Enum,
				Items:       itemsSchema,
			}
			if p.Required {
				required = append(required, p.Name)
			}
		}

		if len(required) == 0 {
			required = nil
		}

		var paramsSchema *Schema
		if len(properties) > 0 {
			paramsSchema = &Schema{
				Type:       "OBJECT",
				Properties: properties,
				Required:   required,
			}
		}

		funcDecls = append(funcDecls, FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  paramsSchema,
		})
	}

	return []Tool{
		{
			FunctionDeclarations: funcDecls,
		},
	}
}

// isThoughtSignature reports whether an ID is a Gemini thought signature
// (as opposed to a synthetic call_* tool-call ID).
func isThoughtSignature(id string) bool {
	return id != "" && !strings.HasPrefix(id, "call_")
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
			// Gemini 3 requires thoughtSignature on the same Part as the
			// first functionCall — never as a standalone Part, and never
			// on functionResponse Parts.
			sigAttached := false
			for _, call := range msg.ToolCalls {
				part := Part{
					FunctionCall: &FunctionCall{
						Name: call.Name,
						Args: call.Args,
					},
				}
				if !sigAttached && isThoughtSignature(call.ID) {
					part.ThoughtSignature = call.ID
					sigAttached = true
				}
				parts = append(parts, part)
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
					Role:  "user",
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

	var thoughtSig string

	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			textBuilder.WriteString(part.Text)
		}
		// Prefer signature attached to the functionCall part itself.
		if part.FunctionCall != nil && part.ThoughtSignature != "" {
			thoughtSig = part.ThoughtSignature
		} else if part.ThoughtSignature != "" && thoughtSig == "" {
			thoughtSig = part.ThoughtSignature
		}
		if part.FunctionCall != nil {
			toolCalls = append(toolCalls, ai.ToolCall{
				ID:   "",
				Name: part.FunctionCall.Name,
				Args: part.FunctionCall.Args,
			})
		}
	}

	for i := range toolCalls {
		if i == 0 && thoughtSig != "" {
			// Only the first parallel function call carries the signature.
			toolCalls[i].ID = thoughtSig
		} else {
			toolCalls[i].ID = fmt.Sprintf("call_%s_%d", toolCalls[i].Name, i)
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
