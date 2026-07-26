package chat

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getWorkspaceFiles returns a list of files in the workspace.
// It prioritizes git ls-files if available, otherwise falls back to basic walk.
func getWorkspaceFiles() []string {
	out, err := exec.Command("git", "ls-files").Output()
	if err == nil {
		files := strings.Split(strings.TrimSpace(string(out)), "\n")
		var validFiles []string
		for _, f := range files {
			if f != "" {
				validFiles = append(validFiles, f)
			}
		}
		if len(validFiles) > 0 {
			return validFiles
		}
	}

	// Fallback
	var files []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

func FilterFiles(files []string, query string) []string {
	if query == "" {
		if len(files) > 10 {
			return files[:10]
		}
		return files
	}

	var filtered []string
	query = strings.ToLower(query)
	for _, f := range files {
		if strings.Contains(strings.ToLower(f), query) {
			filtered = append(filtered, f)
		}
		if len(filtered) >= 10 { // Limit to 10 results for performance
			break
		}
	}
	return filtered
}
