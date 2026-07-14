package file

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// ReadTool struct
type ReadTool struct{}

// NewReadTool creates a new read file tool
func NewReadTool() *ReadTool {
	return &ReadTool{}
}

// Definition returns the tool definition
func (t *ReadTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "read_file",
		Description: "Reads the contents of a file",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to read",
				Required:    true,
			},
		},
	}
}

// Run executes the read file tool
func (t *ReadTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return tools.Result{
			Error: err,
		}
	}

	return tools.Result{
		Output: string(data),
	}
}