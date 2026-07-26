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
		Category:    tools.CategoryFilesystem,
		Permission:  tools.PermWrite,
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

	info, err := os.Stat(path)
	if err != nil {
		return tools.Result{Error: err}
	}

	beforeContent := ""
	if !info.IsDir() {
		beforeBytes, err := os.ReadFile(path)
		if err == nil {
			beforeContent = string(beforeBytes)
		}
	}

	if err := os.RemoveAll(path); err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{
		Output:       fmt.Sprintf("Deleted %q", path),
		FilesChanged: []string{path},
		FileStates: []tools.FileState{
			{
				Path:          path,
				BeforeContent: beforeContent,
				AfterContent:  "",
				ChangeType:    "delete",
			},
		},
	}
}
