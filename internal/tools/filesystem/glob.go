package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/bmatcuk/doublestar/v4"
)

type GlobTool struct{}

func NewGlobTool() *GlobTool {
	return &GlobTool{}
}

func (t *GlobTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "glob",
		Description: "Finds files by matching a pattern (e.g., *.go, **/*.js) across the workspace.",
		Category:    tools.CategorySearch,
		Permission:  tools.PermReadOnly,
		Parameters: []tools.Parameter{
			{
				Name:        "pattern",
				Type:        "string",
				Description: "The glob pattern to search for (supports ** for recursive).",
				Required:    true,
			},
			{
				Name:        "path",
				Type:        "string",
				Description: "The base directory to start searching from. Defaults to current directory.",
				Required:    false,
			},
		},
	}
}

func (t *GlobTool) Run(ctx context.Context, call tools.Call) tools.Result {
	pattern, ok := call.Args["pattern"].(string)
	if !ok || pattern == "" {
		return tools.Result{
			Error: fmt.Errorf("pattern is required"),
		}
	}

	basePath := "."
	if p, ok := call.Args["path"].(string); ok && p != "" {
		basePath = p
	}

	fsys := os.DirFS(basePath)
	matches, err := doublestar.Glob(fsys, pattern)
	if err != nil {
		return tools.Result{
			Error: fmt.Errorf("glob error: %w", err),
		}
	}

	return tools.Result{
		Output: matches,
	}
}
