package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true, "dist": true,
	"build": true, "target": true, ".next": true, ".gemini": true,
	"__pycache__": true, ".venv": true, "venv": true, "coverage": true,
}

// RepoMap builds a compact workspace tree for the system prompt.
func RepoMap(cwd string) string {
	var sb strings.Builder
	sb.WriteString("## Repo map (")
	sb.WriteString(cwd)
	sb.WriteString(")\n```\n")

	fileCount := 0
	maxFiles := 80 // keep free-tier prompts lean
	maxDepth := 2

	_ = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil || path == cwd {
			return nil
		}

		if d.IsDir() && ignoredDirs[d.Name()] {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(cwd, path)
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		fileCount++
		if fileCount > maxFiles {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		indent := strings.Repeat("  ", depth)
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		sb.WriteString(indent)
		sb.WriteString(name)
		sb.WriteByte('\n')
		return nil
	})

	if fileCount > maxFiles {
		sb.WriteString(fmt.Sprintf("…[+%d more; use list_dir]\n", fileCount-maxFiles))
	}
	sb.WriteString("```")
	return sb.String()
}
