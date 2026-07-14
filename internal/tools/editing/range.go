package editing

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ReplaceRange replaces lines from start_line to end_line (1-indexed, inclusive) with new_text.
// It enforces strict fail-fast validation without clamping invalid line coordinates.
func ReplaceRange(ctx context.Context, opts ReplaceRangeOptions) (*ReplaceRangeResult, error) {
	if opts.File == "" {
		return nil, fmt.Errorf("file cannot be empty")
	}
	if opts.StartLine < 1 || opts.EndLine < opts.StartLine {
		return nil, fmt.Errorf("invalid line coordinates [%d, %d]: start_line must be >= 1 and end_line >= start_line", opts.StartLine, opts.EndLine)
	}

	f, err := os.Open(opts.File)
	if err != nil {
		return nil, err
	}
	binary, binErr := isBinaryFile(f)
	_ = f.Close()
	if binErr != nil {
		return nil, binErr
	}
	if binary {
		return nil, fmt.Errorf("cannot perform range replacement on binary file %q", opts.File)
	}

	content, err := os.ReadFile(opts.File)
	if err != nil {
		return nil, err
	}

	lines := []string{}
	hasTrailingNewline := false
	if len(content) > 0 {
		hasTrailingNewline = strings.HasSuffix(string(content), "\n")
		lines = strings.Split(string(content), "\n")
		if hasTrailingNewline && len(lines) > 1 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
	}

	totalLines := len(lines)
	if opts.StartLine > totalLines || opts.EndLine > totalLines {
		return nil, fmt.Errorf("invalid line range [%d, %d]: file %q has %d lines. Call read_context to verify exact line coordinates before replacing", opts.StartLine, opts.EndLine, opts.File, totalLines)
	}

	oldSlice := lines[opts.StartLine-1 : opts.EndLine]
	oldText := strings.Join(oldSlice, "\n")

	var newSlice []string
	if opts.NewText != "" {
		trimmedNew := strings.TrimSuffix(opts.NewText, "\n")
		newSlice = strings.Split(trimmedNew, "\n")
	}

	updatedLines := make([]string, 0, len(lines)-len(oldSlice)+len(newSlice))
	updatedLines = append(updatedLines, lines[:opts.StartLine-1]...)
	if len(newSlice) > 0 {
		updatedLines = append(updatedLines, newSlice...)
	}
	updatedLines = append(updatedLines, lines[opts.EndLine:]...)

	updatedContent := strings.Join(updatedLines, "\n")
	if hasTrailingNewline && (len(updatedContent) == 0 || !strings.HasSuffix(updatedContent, "\n")) {
		updatedContent += "\n"
	}

	if err := WriteFile(opts.File, []byte(updatedContent), 0); err != nil {
		return nil, err
	}

	newEndLine := opts.StartLine - 1
	if len(newSlice) > 0 {
		newEndLine = opts.StartLine + len(newSlice) - 1
	}

	op := Operation{
		Type:      Replace,
		File:      opts.File,
		StartLine: opts.StartLine,
		EndLine:   newEndLine,
		OldText:   oldText,
		NewText:   opts.NewText,
	}

	return &ReplaceRangeResult{
		Operation: op,
	}, nil
}
