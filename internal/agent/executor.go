package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/store"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/pmezard/go-difflib/difflib"
)

// execute runs a slice of tool calls against the tool manager and returns their results.
func (a *Agent) execute(ctx context.Context, calls []ai.ToolCall, onChunk func(string)) []ai.ToolResult {
	results := make([]ai.ToolResult, len(calls))
	var wg sync.WaitGroup

	batchID := fmt.Sprintf("batch_%d", time.Now().UnixNano())

	// Clear redo history when a new edit is made
	if a.config.Store != nil && a.config.SessionID != "" {
		_ = a.config.Store.ClearRedoHistory(a.config.SessionID)
	}

	for i, call := range calls {
		wg.Add(1)
		go func(i int, call ai.ToolCall) {
			defer wg.Done()

			// Route to MCP Manager if it's an MCP tool
			if strings.HasPrefix(call.Name, "mcp_") && a.mcpManager != nil {
				res, err := a.mcpManager.ExecuteTool(ctx, call.Name, call.Args)

				content := ""
				isError := false
				if err != nil {
					content = fmt.Sprintf("MCP error: %v", err)
					isError = true
				} else {
					content = fmt.Sprintf("%v", res)
				}

				results[i] = ai.ToolResult{
					ID:      call.ID,
					Name:    call.Name,
					Content: content,
					IsError: isError,
				}
				return
			}

			tool, ok := a.manager.Get(call.Name)
			if !ok {
				results[i] = ai.ToolResult{
					ID:      call.ID,
					Name:    call.Name,
					Content: fmt.Sprintf("error: tool %q not found or not registered", call.Name),
					IsError: true,
				}
				return
			}

			if onChunk != nil {
				onChunk(fmt.Sprintf("\n\n> ⏳ **Executing tool**: `%s`...", call.Name))
			}

			// Execute the tool.
			res := tool.Run(ctx, tools.Call{
				Name:    call.Name,
				Args:    call.Args,
				OnChunk: onChunk,
			})

			if onChunk != nil {
				onChunk(fmt.Sprintf(" ✅ Done (`%s`).\n\n", call.Name))
			}

			content := ""
			isError := false

			if res.Error != nil {
				content = fmt.Sprintf("error executing tool %s: %v", call.Name, res.Error)
				isError = true
			} else if len(res.FileStates) > 0 {
				var diffs strings.Builder
				diffs.WriteString(fmt.Sprintf("Successfully modified %d file(s):\n\n", len(res.FileStates)))

				for i := range res.FileStates {
					state := &res.FileStates[i]

					// Auto-format the file if possible
					if autoFormat(state.Path) {
						// Re-read the formatted content
						if contentBytes, err := os.ReadFile(state.Path); err == nil {
							state.AfterContent = string(contentBytes)
						}
					}

					// Now save the file change to the store (with formatted content)
					if a.config.Store != nil && a.config.SessionID != "" {
						_ = a.config.Store.SaveFileChange(&store.FileChange{
							SessionID:     a.config.SessionID,
							BatchID:       batchID,
							FilePath:      state.Path,
							ChangeType:    state.ChangeType,
							BeforeContent: state.BeforeContent,
							AfterContent:  state.AfterContent,
						})
					}

					diff := difflib.UnifiedDiff{
						A:        difflib.SplitLines(state.BeforeContent),
						B:        difflib.SplitLines(state.AfterContent),
						FromFile: "a/" + state.Path,
						ToFile:   "b/" + state.Path,
						Context:  3,
					}
					text, _ := difflib.GetUnifiedDiffString(diff)
					diffs.WriteString(fmt.Sprintf("```diff\n%s\n```\n", strings.TrimSpace(text)))

					// Connect to LSP and check for diagnostics
					if a.lspManager != nil {
						absPath, err := filepath.Abs(state.Path)
						if err == nil {
							client, err := a.lspManager.GetClient(ctx, ".", absPath)
							if err == nil && client != nil {
								uri := "file://" + absPath
								// Trigger a didOpen/didChange or simply wait for the server
								// to send diagnostics based on file watching, or explicitly send them
								_ = client.Notify("textDocument/didOpen", map[string]interface{}{
									"textDocument": map[string]interface{}{
										"uri":        uri,
										"languageId": "",
										"version":    1,
										"text":       state.AfterContent,
									},
								})

								// Wait for diagnostics to stream in
								time.Sleep(500 * time.Millisecond)

								diags := client.GetDiagnostics(uri)
								if len(diags) > 0 {
									diffs.WriteString("\n⚠️ **LSP Diagnostics Found:**\n")
									for _, d := range diags {
										if d.Severity == 1 { // Error only
											diffs.WriteString(fmt.Sprintf("- [%s] %s\n", d.Source, d.Message))
										}
									}
								}
							}
						}
					}
				}

				content = diffs.String()

				// Send the diff to the chat UI via onChunk so the user sees it immediately
				if onChunk != nil {
					onChunk("\n" + content + "\n")
				}
			} else if res.Output != nil {
				content = fmt.Sprintf("%v", res.Output)
			} else {
				content = "success"
			}

			results[i] = ai.ToolResult{
				ID:      call.ID,
				Name:    call.Name,
				Content: content,
				IsError: isError,
			}
		}(i, call)
	}

	wg.Wait()
	return results
}

// toolDefinitions converts the registered tool definitions from tools.Manager into ai.ToolDefinition format.
func (a *Agent) toolDefinitions(modeConfig ModeConfig) []ai.ToolDefinition {
	if a.manager == nil {
		return nil
	}

	toolsList := FilterTools(a.manager, modeConfig)

	defs := make([]ai.ToolDefinition, 0, len(toolsList))
	for _, def := range toolsList {
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

	if a.mcpManager != nil {
		defs = append(defs, a.mcpManager.GetTools()...)
	}

	return defs
}
