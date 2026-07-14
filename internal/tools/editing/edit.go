package editing

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// applyLineEdit is a shared internal helper for line-buffer operations (ReplaceRange, InsertText, DeleteRange).
// It handles reading, binary validation, line splitting, boundary checks, slice modification, and atomic writing.
func applyLineEdit(ctx context.Context, file string, startLine int, endLine int, newText string, opType OperationType) (*Operation, error) {
	if file == "" {
		return nil, fmt.Errorf("file cannot be empty")
	}

	if startLine < 1 {
		return nil, fmt.Errorf("start_line must be >= 1 (got %d)", startLine)
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	binary, binErr := isBinaryFile(f)
	_ = f.Close()
	if binErr != nil {
		return nil, binErr
	}
	if binary {
		return nil, fmt.Errorf("cannot perform line editing on binary file %q", file)
	}

	content, err := os.ReadFile(file)
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

	var oldText string
	var newSlice []string
	if newText != "" {
		trimmedNew := strings.TrimSuffix(newText, "\n")
		newSlice = strings.Split(trimmedNew, "\n")
	}

	var updatedLines []string
	var opOperationEndLine int

	switch opType {
	case Replace:
		if endLine < startLine {
			return nil, fmt.Errorf("invalid line coordinates [%d, %d]: end_line must be >= start_line", startLine, endLine)
		}
		if startLine > totalLines || endLine > totalLines {
			return nil, fmt.Errorf("invalid line range [%d, %d]: file %q has %d lines. Call read_context to verify exact line coordinates before replacing", startLine, endLine, file, totalLines)
		}
		oldSlice := lines[startLine-1 : endLine]
		oldText = strings.Join(oldSlice, "\n")

		updatedLines = make([]string, 0, len(lines)-len(oldSlice)+len(newSlice))
		updatedLines = append(updatedLines, lines[:startLine-1]...)
		if len(newSlice) > 0 {
			updatedLines = append(updatedLines, newSlice...)
		}
		updatedLines = append(updatedLines, lines[endLine:]...)

		opOperationEndLine = startLine - 1
		if len(newSlice) > 0 {
			opOperationEndLine = startLine + len(newSlice) - 1
		}

	case Insert:
		// For insertion, startLine is the line AT WHICH text is inserted (before existing startLine).
		// If startLine == totalLines + 1, it appends at the end of the file.
		if startLine > totalLines+1 {
			return nil, fmt.Errorf("invalid insertion line %d: file %q only has %d lines", startLine, file, totalLines)
		}
		if len(newSlice) == 0 {
			return nil, fmt.Errorf("cannot insert empty text")
		}

		updatedLines = make([]string, 0, len(lines)+len(newSlice))
		if startLine <= totalLines {
			updatedLines = append(updatedLines, lines[:startLine-1]...)
			updatedLines = append(updatedLines, newSlice...)
			updatedLines = append(updatedLines, lines[startLine-1:]...)
		} else {
			updatedLines = append(updatedLines, lines...)
			updatedLines = append(updatedLines, newSlice...)
		}

		opOperationEndLine = startLine + len(newSlice) - 1

	case Delete:
		if endLine < startLine {
			return nil, fmt.Errorf("invalid line coordinates [%d, %d]: end_line must be >= start_line", startLine, endLine)
		}
		if startLine > totalLines || endLine > totalLines {
			return nil, fmt.Errorf("invalid line range [%d, %d]: file %q has %d lines", startLine, endLine, file, totalLines)
		}
		oldSlice := lines[startLine-1 : endLine]
		oldText = strings.Join(oldSlice, "\n")

		updatedLines = make([]string, 0, len(lines)-len(oldSlice))
		updatedLines = append(updatedLines, lines[:startLine-1]...)
		updatedLines = append(updatedLines, lines[endLine:]...)

		opOperationEndLine = endLine
	default:
		return nil, fmt.Errorf("unsupported edit operation type: %v", opType)
	}

	updatedContent := strings.Join(updatedLines, "\n")
	if hasTrailingNewline && (len(updatedContent) == 0 || !strings.HasSuffix(updatedContent, "\n")) {
		updatedContent += "\n"
	}

	if err := WriteFile(file, []byte(updatedContent), 0); err != nil {
		return nil, err
	}

	return &Operation{
		Type:      opType,
		File:      file,
		StartLine: startLine,
		EndLine:   opOperationEndLine,
		OldText:   oldText,
		NewText:   newText,
	}, nil
}
