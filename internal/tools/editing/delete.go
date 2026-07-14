package editing

import (
	"context"
	"fmt"
)

// DeleteRange deletes lines from start_line to end_line (1-indexed, inclusive).
// It delegates directly to the shared applyLineEdit internal helper.
func DeleteRange(ctx context.Context, opts DeleteOptions) (*DeleteResult, error) {
	if opts.StartLine < 1 || opts.EndLine < opts.StartLine {
		return nil, fmt.Errorf("invalid line coordinates [%d, %d]: start_line must be >= 1 and end_line >= start_line", opts.StartLine, opts.EndLine)
	}

	op, err := applyLineEdit(ctx, opts.File, opts.StartLine, opts.EndLine, "", Delete)
	if err != nil {
		return nil, err
	}

	return &DeleteResult{
		Operation: *op,
	}, nil
}
