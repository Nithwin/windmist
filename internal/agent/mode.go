package agent

import (
	"github.com/Nithwin/WindMist/internal/tools"
)

// Mode represents an operating mode for the agent.
type Mode string

const (
	// ModeAuto automatically decides between Chat, Plan, and Build based on the prompt.
	ModeAuto Mode = "auto"
	// ModeChat is a lightweight conversational mode: no tools, minimal prompt.
	ModeChat Mode = "chat"
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
	// AllowTools controls whether any tools (including MCP) are sent to the model.
	AllowTools bool
	// IncludeRepoMap injects a workspace file tree into the system prompt.
	IncludeRepoMap bool
}

// GetModeConfig returns the configuration for a given mode.
func GetModeConfig(mode Mode) ModeConfig {
	switch mode {
	case ModeChat:
		return ModeConfig{
			Name:           ModeChat,
			Description:    "Lightweight chat. No tools; minimal prompt for greetings and simple questions.",
			SystemPrompt:   "You are WindMist, a concise coding assistant. Answer briefly and helpfully. Do not invent file contents or pretend to have edited code. If the user wants changes in their project, tell them to ask you to implement it.",
			AllowFileEdits: false,
			AllowCommands:  false,
			AllowTools:     false,
			IncludeRepoMap: false,
		}
	case ModePlan:
		return ModeConfig{
			Name:           ModePlan,
			Description:    "Read-only analysis and planning. Cannot edit files.",
			SystemPrompt:   "You are WindMist in PLAN mode. Answer questions, search/read the codebase, and propose plans. Do NOT modify files or run destructive commands.",
			AllowFileEdits: false,
			AllowCommands:  false,
			AllowTools:     true,
			IncludeRepoMap: true,
		}
	default:
		// Default to build (even if auto, execution resolves to chat/plan/build)
		return ModeConfig{
			Name:           ModeBuild,
			Description:    "Full autonomy mode. Can read, write, and execute.",
			SystemPrompt:   "You are WindMist in BUILD mode. Implement features, fix bugs, and edit files directly. Inspect before editing. Prefer the smallest safe change. Verify when practical.",
			AllowFileEdits: true,
			AllowCommands:  true,
			AllowTools:     true,
			IncludeRepoMap: true,
		}
	}
}

// FilterTools returns only the tools allowed by the given ModeConfig.
func FilterTools(manager *tools.Manager, config ModeConfig) []tools.Definition {
	if !config.AllowTools {
		return nil
	}

	var allowed []tools.Definition

	for _, tool := range manager.List() {
		def := tool.Definition()

		// If edits are denied, filter out write/dangerous tools
		if !config.AllowFileEdits && (def.Category == tools.CategoryEditing || def.Permission == tools.PermWrite || def.Permission == tools.PermDangerous) {
			continue
		}

		allowed = append(allowed, def)
	}

	return allowed
}
