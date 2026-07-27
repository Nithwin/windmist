package prompt

import (
	"os"
	"path/filepath"
	"strings"
)

// Build constructs the complete system prompt for WindMist.
// It dynamically generates a map of the workspace if cwd is provided.
func Build(cwd string, modeSystemPrompt string) string {
	if modeSystemPrompt == "" {
		modeSystemPrompt = System()
	}

	sections := []string{
		modeSystemPrompt,
		Developer(),
		Tools(),
	}

	if cwd != "" {
		sections = append(sections, RepoMap(cwd))
		
		// Look for custom AGENTS.md or .windmist/prompt.md conventions
		if custom := loadCustomPrompt(cwd); custom != "" {
			sections = append(sections, "## Workspace Conventions\n\nThe following rules apply specifically to this workspace:\n\n"+custom)
		}
	}

	return strings.Join(sections, "\n\n")
}

// loadCustomPrompt checks for local workspace prompt files.
func loadCustomPrompt(cwd string) string {
	paths := []string{
		filepath.Join(cwd, "AGENTS.md"),
		filepath.Join(cwd, ".windmist", "prompt.md"),
	}

	for _, p := range paths {
		if content, err := os.ReadFile(p); err == nil {
			return string(content)
		}
	}
	return ""
}
