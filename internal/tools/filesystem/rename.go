package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type RenameTool struct{}

func NewRenameTool() *RenameTool {
	return &RenameTool{}
}

func (t *RenameTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "rename",
		Description: "Renames or moves a file or directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "old_path",
				Type:        "string",
				Description: "The original path of the file or directory.",
				Required:    true,
			},
			{
				Name:        "new_path",
				Type:        "string",
				Description: "The new path for the file or directory.",
				Required:    true,
			},
		},
	}
}

func (t *RenameTool) Run(ctx context.Context, call tools.Call) tools.Result {
	oldPath, ok := call.Args["old_path"].(string)
	if !ok || oldPath == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	newPath, ok := call.Args["new_path"].(string)
	if !ok || newPath == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{
		Output: fmt.Sprintf("Renamed %q to %q", oldPath, newPath),
	}
}
