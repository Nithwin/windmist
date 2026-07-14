package editing

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// ReadContext reads a window of lines around targetLine with 1-indexed line prefixes.
func ReadContext(ctx context.Context, path string, targetLine int, before int, after int) (string, error) {
	if targetLine < 1 {
		targetLine = 1
	}
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}

	startLine := targetLine - before
	if startLine < 1 {
		startLine = 1
	}
	endLine := targetLine + after
	if endLine < startLine {
		endLine = startLine
	}

	return ReadRange(ctx, path, startLine, endLine)
}

// ReadRange reads lines from startLine to endLine (inclusive) formatted with 1-indexed line numbers.
func ReadRange(ctx context.Context, path string, startLine int, endLine int) (string, error) {
	if startLine < 1 {
		startLine = 1
	}
	if endLine < startLine {
		endLine = startLine
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanBuf := make([]byte, 0, 64*1024)
	scanner.Buffer(scanBuf, 10*1024*1024)

	var builder strings.Builder
	currentLine := 1

	for scanner.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if currentLine > endLine {
			break
		}
		if currentLine >= startLine {
			builder.WriteString(fmt.Sprintf("%d: %s\n", currentLine, scanner.Text()))
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if startLine > currentLine-1 && currentLine > 1 {
		return "", fmt.Errorf("target start line %d exceeds file length (%d lines)", startLine, currentLine-1)
	}

	return builder.String(), nil
}


