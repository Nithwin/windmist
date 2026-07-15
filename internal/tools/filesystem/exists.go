package filesystem

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type ExistsTool struct{}

func NewExistsTool() *ExistsTool {
	return &ExistsTool{}
}

func (t *ExistsTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "exists",
		Description: "Checks if a file or directory exists.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path to check.",
				Required:    true,
			},
		},
	}
}

func (t *ExistsTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	info, err := os.Stat(path)
	if err == nil {
		entryType := "file"
		if info.IsDir() {
			entryType = "directory"
		}
		return tools.Result{
			Output: ExistsResult{
				Exists: true,
				Type:   entryType,
			},
		}
	}

	if os.IsNotExist(err) {
		return tools.Result{
			Output: ExistsResult{
				Exists: false,
			},
		}
	}

	return tools.Result{
		Error: err,
	}
}
