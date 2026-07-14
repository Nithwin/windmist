package editing

import (
	"context"
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// InsertTextTool implements the AI tool interface for insert_text.
type InsertTextTool struct{}

func NewInsertTextTool() *InsertTextTool {
	return &InsertTextTool{}
}

func (t *InsertTextTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "insert_text",
		Description: "Inserts text at a specific 1-indexed line number.",
		Parameters: []tools.Parameter{
			{
				Name:        "file",
				Type:        "string",
				Description: "Path of the file to modify.",
				Required:    true,
			},
			{
				Name:        "line",
				Type:        "integer",
				Description: "1-indexed target line number where text will be inserted.",
				Required:    true,
			},
			{
				Name:        "new_text",
				Type:        "string",
				Description: "New text to insert.",
				Required:    true,
			},
		},
	}
}

func (t *InsertTextTool) Run(ctx context.Context, call tools.Call) tools.Result {
	file, ok := call.Args["file"].(string)
	if !ok || file == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	line := 0
	if l, ok := call.Args["line"].(int); ok {
		line = l
	} else if lFloat, ok := call.Args["line"].(float64); ok {
		line = int(lFloat)
	}
	if line < 1 {
		return tools.Result{Error: fmt.Errorf("line must be >= 1")}
	}

	newText, ok := call.Args["new_text"].(string)
	if !ok {
		return tools.Result{Error: fmt.Errorf("new_text must be a string")}
	}

	opts := InsertOptions{
		File:    file,
		Line:    line,
		NewText: newText,
	}

	result, err := InsertText(ctx, opts)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: result}
}
