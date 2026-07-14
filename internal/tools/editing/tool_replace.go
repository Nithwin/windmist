package editing

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// ReplaceTextTool implements the AI tool interface for replace_text.
type ReplaceTextTool struct{}

func NewReplaceTextTool() *ReplaceTextTool {
	return &ReplaceTextTool{}
}

func (t *ReplaceTextTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "replace_text",
		Description: "Replace a unique piece of text in an existing file. Use this when the target text is known exactly. Prefer range-based editing when exact line numbers are available.",
		Parameters: []tools.Parameter{
			{
				Name:        "file",
				Type:        "string",
				Description: "Path of the file to modify.",
				Required:    true,
			},
			{
				Name:        "old_text",
				Type:        "string",
				Description: "Exact string pattern to replace.",
				Required:    true,
			},
			{
				Name:        "new_text",
				Type:        "string",
				Description: "Replacement string.",
				Required:    true,
			},
			{
				Name:        "replace_all",
				Type:        "boolean",
				Description: "Whether to replace all occurrences if multiple exist (defaults to false).",
				Required:    false,
			},
			{
				Name:        "max_replacements",
				Type:        "integer",
				Description: "Maximum number of occurrences to replace if replace_all is false.",
				Required:    false,
			},
		},
	}
}

func (t *ReplaceTextTool) Run(ctx context.Context, call tools.Call) tools.Result {
	file, ok := call.Args["file"].(string)
	if !ok || file == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	oldText, ok := call.Args["old_text"].(string)
	if !ok || oldText == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	newText, _ := call.Args["new_text"].(string)

	replaceAll := false
	if ra, ok := call.Args["replace_all"].(bool); ok {
		replaceAll = ra
	} else if raStr, ok := call.Args["replace_all"].(string); ok {
		replaceAll = (raStr == "true" || raStr == "1")
	}

	maxReplacements := 0
	if mr, ok := call.Args["max_replacements"].(int); ok && mr > 0 {
		maxReplacements = mr
	} else if mrFloat, ok := call.Args["max_replacements"].(float64); ok && mrFloat > 0 {
		maxReplacements = int(mrFloat)
	}

	opts := ReplaceOptions{
		File:            file,
		OldText:         oldText,
		NewText:         newText,
		ReplaceAll:      replaceAll,
		MaxReplacements: maxReplacements,
	}

	result, err := ReplaceText(ctx, opts)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: result}
}
