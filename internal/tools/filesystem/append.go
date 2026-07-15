package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

type AppendTool struct{}

func NewAppendTool() *AppendTool {
	return &AppendTool{}
}

func (t *AppendTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "append",
		Description: "Appends content to an existing file.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to append to.",
				Required:    true,
			},
			{
				Name:        "content",
				Type:        "string",
				Description: "Content to append.",
				Required:    true,
			},
		},
	}
}

func (t *AppendTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	content, ok := call.Args["content"].(string)
	if !ok {
		return tools.Result{Error: os.ErrInvalid}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return tools.Result{Error: err}
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{
		Output: fmt.Sprintf("Appended %d bytes to %q", len(content), path),
	}
}
