package prompt

import (
	"os"
	"path/filepath"
	"strings"
)

// Options controls which optional sections are included in the system prompt.
type Options struct {
	IncludeRepoMap bool
	IncludeGuides  bool // developer + tool workflow guidance
}

// Build constructs the system prompt for WindMist.
func Build(cwd string, modeSystemPrompt string, opts Options) string {
	if modeSystemPrompt == "" {
		modeSystemPrompt = System()
	}

	sections := []string{modeSystemPrompt}

	if opts.IncludeGuides {
		sections = append(sections, Workflow())
	}

	if opts.IncludeRepoMap && cwd != "" {
		sections = append(sections, RepoMap(cwd))
	}

	if cwd != "" {
		if custom := loadCustomPrompt(cwd); custom != "" {
			// Cap custom prompts to avoid blowing free-tier context.
			if len(custom) > 4000 {
				custom = custom[:4000] + "\n…[truncated]"
			}
			sections = append(sections, "## Workspace Conventions\n\n"+custom)
		}
	}

	return strings.Join(sections, "\n\n")
}

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
