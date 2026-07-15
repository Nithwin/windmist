package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type DeleteTool struct{}

func NewDeleteTool() *DeleteTool {
	return &DeleteTool{}
}

func (t *DeleteTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "delete",
		Description: "Deletes a file or directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "The path to the file or directory to delete.",
				Required:    true,
			},
		},
	}
}

func (t *DeleteTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	if _, err := os.Stat(path); err != nil {
		return tools.Result{Error: err}
	}

	if err := os.RemoveAll(path); err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{
		Output: fmt.Sprintf("Deleted %q", path),
	}
}
