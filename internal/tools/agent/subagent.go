package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
)

type subAgentArgs struct {
	Task  string   `json:"task"`
	Files []string `json:"files"`
}

type subAgentTool struct {
	cfg *config.Config
}

// NewSubAgentTool creates a tool that delegates tasks to a smaller/faster LLM model.
func NewSubAgentTool(cfg *config.Config) tools.Tool {
	return &subAgentTool{cfg: cfg}
}

func (t *subAgentTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "spawn_subagent",
		Description: "Spawns a sub-agent to read multiple files and summarize or analyze them based on a specific task. Use this to prevent cluttering the main context window when you need to research a large codebase.",
		Parameters: []tools.Parameter{
			{
				Name:        "task",
				Type:        "string",
				Description: "The specific research or analysis task for the sub-agent. E.g., 'Analyze how authentication is implemented and list the JWT secret name'.",
				Required:    true,
			},
			{
				Name:        "files",
				Type:        "array",
				Description: "List of exact file paths to read and analyze.",
				Required:    true,
			},
		},
	}
}

func (t *subAgentTool) Run(ctx context.Context, call tools.Call) tools.Result {
	taskStr, ok := call.Args["task"].(string)
	if !ok {
		return tools.Result{Error: fmt.Errorf("task must be a string")}
	}

	filesRaw, ok := call.Args["files"].([]any)
	if !ok {
		return tools.Result{Error: fmt.Errorf("files must be an array")}
	}

	var files []string
	for _, f := range filesRaw {
		if fs, ok := f.(string); ok {
			files = append(files, fs)
		}
	}

	if len(files) == 0 {
		return tools.Result{Error: fmt.Errorf("no valid files provided")}
	}

	// Read all files
	var fileContents string
	var filesRead []string
	for _, file := range files {
		cleanPath := filepath.Clean(file)
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			fileContents += fmt.Sprintf("File: %s\nError reading file: %v\n\n", cleanPath, err)
			continue
		}
		fileContents += fmt.Sprintf("File: %s\n```\n%s\n```\n\n", cleanPath, string(content))
		filesRead = append(filesRead, cleanPath)
	}

	// Prepare AI config using the fast model
	providerName := t.cfg.ActiveSubAgentProvider()
	modelName := t.cfg.ActiveSubAgentModel()

	// Create a temporary config for the sub-agent
	subCfg := &config.Config{
		AI: config.AIConfig{Provider: providerName},
		Providers: map[string]config.ProviderConfig{
			providerName: {
				Model:  modelName,
			},
		},
	}
	
	// Copy API key/base url from original provider if it exists
	if origProvider, ok := t.cfg.Providers[providerName]; ok {
		p := subCfg.Providers[providerName]
		p.APIKey = origProvider.APIKey
		p.BaseURL = origProvider.BaseURL
		subCfg.Providers[providerName] = p
	}

	provider, err := ai.New(subCfg)
	if err != nil {
		return tools.Result{Error: fmt.Errorf("failed to initialize sub-agent AI provider (%s/%s): %w", providerName, modelName, err)}
	}

	systemPrompt := "You are a specialized sub-agent for an AI coding assistant. Your job is to read the provided files, analyze them, and fulfill the requested task concisely and accurately. Do not write full files, just provide the exact analysis requested."
	
	req := &ai.GenerateRequest{
		System: systemPrompt,
		Messages: []ai.Message{
			{
				Role:    ai.RoleUser,
				Content: fmt.Sprintf("Task: %s\n\nFiles Content:\n%s", taskStr, fileContents),
			},
		},
	}

	resp, err := provider.Generate(ctx, req)
	if err != nil {
		return tools.Result{Error: fmt.Errorf("sub-agent generation failed: %w", err)}
	}

	output := fmt.Sprintf("Sub-Agent Analysis (Model: %s/%s):\n\n%s", providerName, modelName, resp.Text)
	return tools.Result{
		Output:    output,
		FilesRead: filesRead,
	}
}
