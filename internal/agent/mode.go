package agent

import (
	"github.com/Nithwin/WindMist/internal/tools"
)

// Mode represents an operating mode for the agent.
type Mode string

const (
	// ModeAuto automatically decides between Plan and Build based on the prompt.
	ModeAuto Mode = "auto"
	// ModeBuild has full read/write access and autonomy.
	ModeBuild Mode = "build"
	// ModePlan is read-only. It can search and analyze, but cannot write files.
	ModePlan Mode = "plan"
)

// ModeConfig defines the behavior and permissions of a specific mode.
type ModeConfig struct {
	Name           Mode
	Description    string
	SystemPrompt   string
	AllowFileEdits bool
	AllowCommands  bool
}

// GetModeConfig returns the configuration for a given mode.
func GetModeConfig(mode Mode) ModeConfig {
	switch mode {
	case ModePlan:
		return ModeConfig{
			Name:           ModePlan,
			Description:    "Read-only architect mode. Analyzes and plans but cannot edit files.",
			SystemPrompt:   "You are an expert software architect in PLAN mode. Your job is to analyze the user's request, search the codebase, read files, and output a detailed, step-by-step implementation plan. YOU CANNOT MODIFY FILES OR WRITE CODE TO DISK. Do not attempt to use any write tools. If a command must be run, ask the user first. Focus on architectural decisions, edge cases, and producing a clear numbered plan.",
			AllowFileEdits: false,
			AllowCommands:  false,
		}
	default:
		// Default to build (even if auto, the actual execution mode resolves to build/plan)
		return ModeConfig{
			Name:           ModeBuild,
			Description:    "Full autonomy mode. Can read, write, and execute.",
			SystemPrompt:   "You are WindMist, an expert autonomous coding agent in BUILD mode. Your job is to implement features, fix bugs, and refactor code directly. You have full access to the filesystem. When asked to complete a task, you should read relevant files, make the necessary edits using your tools, and run commands to verify your work. Act surgically and efficiently.",
			AllowFileEdits: true,
			AllowCommands:  true,
		}
	}
}

// FilterTools returns only the tools allowed by the given ModeConfig.
func FilterTools(manager *tools.Manager, config ModeConfig) []tools.Definition {
	var allowed []tools.Definition

	for _, tool := range manager.List() {
		def := tool.Definition()
		
		// If edits are denied, filter out PermWrite and PermDangerous
		if !config.AllowFileEdits && (def.Category == tools.CategoryEditing || def.Permission == tools.PermWrite || def.Permission == tools.PermDangerous) {
			continue
		}
		
		// Wait, if commands are denied, we could filter out system/command tools, 
		// but maybe we just require permission instead of completely filtering.
		// For now, in plan mode, we completely disable editing tools.
		
		allowed = append(allowed, def)
	}

	return allowed
}
