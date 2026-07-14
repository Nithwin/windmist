package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type WriteTool struct{}

func NewWriteTool() *WriteTool {
	return &WriteTool{}
}

func (t *WriteTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "write",
		Description: "Overwrites the entire contents of an existing file with new content. WARNING: This replaces all existing code in the file. Prefer using replace_text or replace_range when making targeted edits or modifying existing code.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to write.",
				Required:    true,
			},
			{
				Name:        "content",
				Type:        "string",
				Description: "Content to write to the file.",
				Required:    true,
			},
		},
	}
}

func (t *WriteTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	content, ok := call.Args["content"].(string)
	if !ok {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return tools.Result{
			Error: err,
		}
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return tools.Result{
			Error: err,
		}
	}

	return tools.Result{
		Output: fmt.Sprintf("Wrote %d bytes to %q", len(content), path),
	}
}
