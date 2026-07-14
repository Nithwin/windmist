package editing

import (
	"context"
	"fmt"
)

// InsertText inserts new_text at the specified 1-indexed line number.
// It delegates to the shared applyLineEdit internal helper.
func InsertText(ctx context.Context, opts InsertOptions) (*InsertResult, error) {
	if opts.Line < 1 {
		return nil, fmt.Errorf("line must be >= 1 (got %d)", opts.Line)
	}

	op, err := applyLineEdit(ctx, opts.File, opts.Line, opts.Line, opts.NewText, Insert)
	if err != nil {
		return nil, err
	}

	return &InsertResult{
		Operation: *op,
	}, nil
}
