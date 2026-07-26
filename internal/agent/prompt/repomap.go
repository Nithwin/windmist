package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ignoredDirs is a list of directories we shouldn't map to avoid exploding tokens
var ignoredDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	"target":       true,
	".next":        true,
	".gemini":      true,
	"__pycache__":  true,
}

// RepoMap dynamically builds a file tree of the workspace to inject into the system prompt.
// This prevents the AI from needing to run `list_dir` manually and stops hallucinated file paths.
func RepoMap(cwd string) string {
	var sb strings.Builder
	sb.WriteString("## Repository Map\n\n")
	sb.WriteString("Below is a tree representation of the files in the current workspace.\n")
	sb.WriteString("Use this map to understand the project structure and locate files instantly.\n")
	sb.WriteString("This map is dynamically updated on every turn to reflect the current state of the filesystem.\n\n")
	sb.WriteString("```\n.\n")

	fileCount := 0
	maxFiles := 2000 // hard limit to prevent token blowout

	err := filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}

		if path == cwd {
			return nil
		}

		// Check if we should ignore this directory
		if d.IsDir() {
			if ignoredDirs[d.Name()] {
				return filepath.SkipDir
			}
		}

		fileCount++
		if fileCount > maxFiles {
			// Don't stop WalkDir completely, just skip directories if over limit
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate depth
		relPath, _ := filepath.Rel(cwd, path)
		depth := strings.Count(relPath, string(os.PathSeparator))

		indent := strings.Repeat("  ", depth)
		prefix := "├── "

		sb.WriteString(indent + prefix + d.Name())
		if d.IsDir() {
			sb.WriteString("/")
		}
		sb.WriteString("\n")

		return nil
	})

	if err != nil {
		sb.WriteString("  [Error generating map]\n")
	}

	if fileCount > maxFiles {
		sb.WriteString(fmt.Sprintf("\n... [Map truncated. Project has >%d files. Use list_dir tool to explore further.]\n", maxFiles))
	}

	sb.WriteString("```\n")

	return sb.String()
}
