package editing

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
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

// ReadContextTool implements the AI tool interface for read_context.
type ReadContextTool struct{}

func NewReadContextTool() *ReadContextTool {
	return &ReadContextTool{}
}

func (t *ReadContextTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "read_context",
		Description: "Reads a specific range of lines from a file with 1-indexed line numbers formatted for editing context.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the file to read context from.",
				Required:    true,
			},
			{
				Name:        "line",
				Type:        "integer",
				Description: "Target line number (1-indexed).",
				Required:    true,
			},
			{
				Name:        "before",
				Type:        "integer",
				Description: "Number of context lines to include before the target line (defaults to 10).",
				Required:    false,
			},
			{
				Name:        "after",
				Type:        "integer",
				Description: "Number of context lines to include after the target line (defaults to 10).",
				Required:    false,
			},
		},
	}
}

func (t *ReadContextTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	targetLine := 0
	if l, ok := call.Args["line"].(int); ok {
		targetLine = l
	} else if lFloat, ok := call.Args["line"].(float64); ok {
		targetLine = int(lFloat)
	}
	if targetLine < 1 {
		return tools.Result{Error: fmt.Errorf("line must be >= 1")}
	}

	before := 10
	if b, ok := call.Args["before"].(int); ok && b >= 0 {
		before = b
	} else if bFloat, ok := call.Args["before"].(float64); ok && bFloat >= 0 {
		before = int(bFloat)
	}

	after := 10
	if a, ok := call.Args["after"].(int); ok && a >= 0 {
		after = a
	} else if aFloat, ok := call.Args["after"].(float64); ok && aFloat >= 0 {
		after = int(aFloat)
	}

	output, err := ReadContext(ctx, path, targetLine, before, after)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: output}
}
