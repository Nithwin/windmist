package editing

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// ReadContextTool implements the AI tool interface for read_context.
type ReadContextTool struct{}

func NewReadContextTool() *ReadContextTool {
	return &ReadContextTool{}
}

func (t *ReadContextTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "read_context",
		Description: "Reads a specific range of lines from a file with 1-indexed line numbers formatted for editing context.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to read context from.",
				Required:    true,
			},
			{
				Name:        "line",
				Type:        "integer",
				Description: "Target line number (1-indexed).",
				Required:    true,
			},
			{
				Name:        "before",
				Type:        "integer",
				Description: "Number of context lines to include before the target line (defaults to 10).",
				Required:    false,
			},
			{
				Name:        "after",
				Type:        "integer",
				Description: "Number of context lines to include after the target line (defaults to 10).",
				Required:    false,
			},
		},
	}
}

func (t *ReadContextTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	targetLine := 0
	if l, ok := call.Args["line"].(int); ok {
		targetLine = l
	} else if lFloat, ok := call.Args["line"].(float64); ok {
		targetLine = int(lFloat)
	}
	if targetLine < 1 {
		return tools.Result{Error: fmt.Errorf("line must be >= 1")}
	}

	before := 10
	if b, ok := call.Args["before"].(int); ok && b >= 0 {
		before = b
	} else if bFloat, ok := call.Args["before"].(float64); ok && bFloat >= 0 {
		before = int(bFloat)
	}

	after := 10
	if a, ok := call.Args["after"].(int); ok && a >= 0 {
		after = a
	} else if aFloat, ok := call.Args["after"].(float64); ok && aFloat >= 0 {
		after = int(aFloat)
	}

	output, err := ReadContext(ctx, path, targetLine, before, after)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: output}
}
