package editing

import (
	"context"
	"fmt"
)

// ReplaceRange replaces lines from start_line to end_line (1-indexed, inclusive) with new_text.
// It delegates to the shared applyLineEdit internal helper.
func ReplaceRange(ctx context.Context, opts ReplaceRangeOptions) (*ReplaceRangeResult, error) {
	if opts.StartLine < 1 || opts.EndLine < opts.StartLine {
		return nil, fmt.Errorf("invalid line coordinates [%d, %d]: start_line must be >= 1 and end_line >= start_line", opts.StartLine, opts.EndLine)
	}

	op, err := applyLineEdit(ctx, opts.File, opts.StartLine, opts.EndLine, opts.NewText, Replace)
	if err != nil {
		return nil, err
	}

	return &ReplaceRangeResult{
		Operation: *op,
	}, nil
}
