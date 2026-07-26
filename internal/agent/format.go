package agent

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// autoFormat attempts to format the given file using standard language formatters.
// It returns true if a formatter was successfully run, or false if no formatter was found or it failed.
func autoFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	var cmd *exec.Cmd

	switch ext {
	case ".go":
		if _, err := exec.LookPath("gofmt"); err == nil {
			cmd = exec.Command("gofmt", "-w", path)
		}
	case ".js", ".ts", ".jsx", ".tsx", ".json", ".css", ".md", ".html":
		if _, err := exec.LookPath("prettier"); err == nil {
			cmd = exec.Command("prettier", "--write", path)
		}
	case ".py":
		if _, err := exec.LookPath("black"); err == nil {
			cmd = exec.Command("black", path)
		} else if _, err := exec.LookPath("ruff"); err == nil {
			cmd = exec.Command("ruff", "format", path)
		}
	case ".rs":
		if _, err := exec.LookPath("rustfmt"); err == nil {
			cmd = exec.Command("rustfmt", path)
		}
	}

	if cmd == nil {
		return false
	}

	// We don't care about the output right now, just run it silently
	err := cmd.Run()
	return err == nil
}
