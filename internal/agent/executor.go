package agent

import (
	"context"
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/tools"
)

// execute runs a slice of tool calls against the tool manager and returns their results.
func (a *Agent) execute(ctx context.Context, calls []ai.ToolCall) []ai.ToolResult {
	results := make([]ai.ToolResult, 0, len(calls))

	for _, call := range calls {
		tool, ok := a.manager.Get(call.Name)
		if !ok {
			results = append(results, ai.ToolResult{
				ID:      call.ID,
				Name:    call.Name,
				Content: fmt.Sprintf("error: tool %q not found or not registered", call.Name),
				IsError: true,
			})
			continue
		}

		// Execute the tool.
		res := tool.Run(ctx, tools.Call{
			Name: call.Name,
			Args: call.Args,
		})

		content := ""
		isError := false

		if res.Error != nil {
			content = fmt.Sprintf("error executing tool %s: %v", call.Name, res.Error)
			isError = true
		} else if res.Output != nil {
			content = fmt.Sprintf("%v", res.Output)
		} else {
			content = "success"
		}

		results = append(results, ai.ToolResult{
			ID:      call.ID,
			Name:    call.Name,
			Content: content,
			IsError: isError,
		})
	}

	return results
}

// toolDefinitions converts the registered tool definitions from tools.Manager into ai.ToolDefinition format.
func (a *Agent) toolDefinitions() []ai.ToolDefinition {
	if a.manager == nil {
		return nil
	}
	toolsList := a.manager.List()
	defs := make([]ai.ToolDefinition, 0, len(toolsList))
	for _, t := range toolsList {
		def := t.Definition()
		params := make([]ai.ToolParameter, 0, len(def.Parameters))
		for _, p := range def.Parameters {
			params = append(params, ai.ToolParameter{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Description,
				Required:    p.Required,
				Enum:        p.Enum,
			})
		}
		defs = append(defs, ai.ToolDefinition{
			Name:        def.Name,
			Description: def.Description,
			Parameters:  params,
		})
	}
	return defs
}
