package editing

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// DeleteRangeTool implements the AI tool interface for delete_range.
type DeleteRangeTool struct{}

func NewDeleteRangeTool() *DeleteRangeTool {
	return &DeleteRangeTool{}
}

func (t *DeleteRangeTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "delete_range",
		Description: "Deletes exact 1-indexed line ranges from a file.",
		Parameters: []tools.Parameter{
			{
				Name:        "file",
				Type:        "string",
				Description: "Path of the file to modify.",
				Required:    true,
			},
			{
				Name:        "start_line",
				Type:        "integer",
				Description: "1-indexed starting line number to delete (inclusive).",
				Required:    true,
			},
			{
				Name:        "end_line",
				Type:        "integer",
				Description: "1-indexed ending line number to delete (inclusive).",
				Required:    true,
			},
		},
	}
}

func (t *DeleteRangeTool) Run(ctx context.Context, call tools.Call) tools.Result {
	file, ok := call.Args["file"].(string)
	if !ok || file == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	startLine := 0
	if sl, ok := call.Args["start_line"].(int); ok {
		startLine = sl
	} else if slFloat, ok := call.Args["start_line"].(float64); ok {
		startLine = int(slFloat)
	}

	endLine := 0
	if el, ok := call.Args["end_line"].(int); ok {
		endLine = el
	} else if elFloat, ok := call.Args["end_line"].(float64); ok {
		endLine = int(elFloat)
	}

	opts := DeleteOptions{
		File:      file,
		StartLine: startLine,
		EndLine:   endLine,
	}

	result, err := DeleteRange(ctx, opts)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: result}
}
