package editing

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// ReplaceRangeTool implements the AI tool interface for replace_range.
type ReplaceRangeTool struct{}

func NewReplaceRangeTool() *ReplaceRangeTool {
	return &ReplaceRangeTool{}
}

func (t *ReplaceRangeTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "replace_range",
		Description: "Replace a contiguous range of lines (1-indexed, inclusive) in an existing file with new text. Use this when you know the exact line numbers from reading context. Preferred over replace_text when the target string appears multiple times in the file.",
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
				Description: "1-indexed starting line number (inclusive).",
				Required:    true,
			},
			{
				Name:        "end_line",
				Type:        "integer",
				Description: "1-indexed ending line number (inclusive).",
				Required:    true,
			},
			{
				Name:        "new_text",
				Type:        "string",
				Description: "New text to replace the specified lines with.",
				Required:    true,
			},
		},
	}
}

func (t *ReplaceRangeTool) Run(ctx context.Context, call tools.Call) tools.Result {
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

	newText, _ := call.Args["new_text"].(string)

	opts := ReplaceRangeOptions{
		File:      file,
		StartLine: startLine,
		EndLine:   endLine,
		NewText:   newText,
	}

	result, err := ReplaceRange(ctx, opts)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: result}
}
