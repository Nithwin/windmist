package filesystem

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/bmatcuk/doublestar/v4"
)

type GrepTool struct{}

func NewGrepTool() *GrepTool {
	return &GrepTool{}
}

func (t *GrepTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "grep",
		Description: "Searches for a regex pattern inside files across the workspace. Returns file path, line number, and the matching line.",
		Category:    tools.CategorySearch,
		Permission:  tools.PermReadOnly,
		Parameters: []tools.Parameter{
			{
				Name:        "pattern",
				Type:        "string",
				Description: "The regular expression pattern to search for.",
				Required:    true,
			},
			{
				Name:        "path",
				Type:        "string",
				Description: "The directory to search in (defaults to current directory).",
				Required:    false,
			},
			{
				Name:        "include",
				Type:        "string",
				Description: "Optional glob pattern to filter files (e.g., *.go).",
				Required:    false,
			},
		},
	}
}

type GrepMatch struct {
	File    string `json:"file"`
	LineNum int    `json:"line"`
	Content string `json:"content"`
}

func (t *GrepTool) Run(ctx context.Context, call tools.Call) tools.Result {
	pattern, ok := call.Args["pattern"].(string)
	if !ok || pattern == "" {
		return tools.Result{Error: fmt.Errorf("pattern is required")}
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return tools.Result{Error: fmt.Errorf("invalid regex pattern: %w", err)}
	}

	basePath := "."
	if p, ok := call.Args["path"].(string); ok && p != "" {
		basePath = p
	}

	includeGlob := ""
	if inc, ok := call.Args["include"].(string); ok && inc != "" {
		includeGlob = inc
	}

	var matches []GrepMatch
	maxMatches := 200 // Cap to prevent massive outputs

	err = filepath.WalkDir(basePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if d.IsDir() {
			// Skip .git and common vendor/binary folders
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == ".windmist" {
				return filepath.SkipDir
			}
			return nil
		}

		if includeGlob != "" {
			// Check if file matches include glob
			rel, _ := filepath.Rel(basePath, path)
			if rel == "" {
				rel = path
			}
			matched, _ := doublestar.Match(includeGlob, rel)
			if !matched {
				return nil
			}
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				matches = append(matches, GrepMatch{
					File:    path,
					LineNum: lineNum,
					Content: strings.TrimSpace(line),
				})
				if len(matches) >= maxMatches {
					return fmt.Errorf("max matches reached")
				}
			}
			lineNum++
		}
		return nil
	})

	if err != nil && err.Error() != "max matches reached" {
		return tools.Result{Error: err}
	}

	return tools.Result{
		Output: map[string]interface{}{
			"matches": matches,
			"limit":   len(matches) == maxMatches,
		},
	}
}
