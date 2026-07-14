package file

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type CreateTool struct{}

func NewCreateTool() *CreateTool {
	return &CreateTool{}
}

func (t *CreateTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "create_file",
		Description: "Creates a new empty file.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to create.",
				Required:    true,
			},
		},
	}
}

func (t *CreateTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return tools.Result{
			Error: fmt.Errorf("file %q already exists", path),
		}
	}

	defer file.Close()

	return tools.Result{
		Output: "File created successfully.",
	}
}
