package filesystem

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

// Definition returns the tool definition
func (t *CreateTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "create",
		Description: "Creates a new file or directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file or directory to create.",
				Required:    true,
			},
			{
				Name:        "type",
				Type:        "string",
				Description: "Type to create: 'file' or 'directory' (defaults to 'file').",
				Required:    false,
				Enum:        []string{"file", "directory"},
			},
		},
	}
}

// Run executes the create tool
func (t *CreateTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	entryType, _ := call.Args["type"].(string)
	if entryType == "directory" || entryType == "dir" {
		if _, err := os.Stat(path); err == nil {
			return tools.Result{
				Error: fmt.Errorf("directory %q already exists", path),
			}
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return tools.Result{
				Error: err,
			}
		}
		return tools.Result{
			Output: "Directory created successfully.",
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
