package prompt

import "strings"

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
	}

	return strings.Join(sections, "\n\n")
}
