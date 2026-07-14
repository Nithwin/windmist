package filesystem

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type InfoTool struct{}

func NewInfoTool() *InfoTool {
	return &InfoTool{}
}

func (t *InfoTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "info",
		Description: "Retrieves metadata and information about a file or directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path to the file or directory.",
				Required:    true,
			},
		},
	}
}

func (t *InfoTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return tools.Result{
			Error: err,
		}
	}

	entryType := "file"
	if info.IsDir() {
		entryType = "directory"
	}

	result := InfoResult{
		Name:         info.Name(),
		Path:         path,
		Size:         info.Size(),
		Type:         entryType,
		Mode:         info.Mode().String(),
		LastModified: info.ModTime(),
	}

	return tools.Result{
		Output: result,
	}
}
