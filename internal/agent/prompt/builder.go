package prompt

import "strings"

// Build constructs the complete system prompt for WindMist.
// Additional prompt sections can be added here as the agent evolves.
func Build() string {
	sections := []string{
		System(),
		Developer(),
		Tools(),
	}

	return strings.Join(sections, "\n\n")
}
