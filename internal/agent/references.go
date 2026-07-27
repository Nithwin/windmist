package agent

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	// refRegex matches @filename or @url. It supports basic file paths and http(s) urls.
	refRegex = regexp.MustCompile(`@([a-zA-Z0-9_./\-:]+)`)
)

// ReferenceInjector parses a user prompt, extracts @ references, and fetches their content.
type ReferenceInjector struct {
	workspaceDir string
}

// NewReferenceInjector creates a new injector.
func NewReferenceInjector(workspaceDir string) *ReferenceInjector {
	return &ReferenceInjector{workspaceDir: workspaceDir}
}

// Inject scans the prompt for @ references and returns a new prompt with the content appended.
func (i *ReferenceInjector) Inject(prompt string) string {
	matches := refRegex.FindAllStringSubmatch(prompt, -1)
	if len(matches) == 0 {
		return prompt
	}

	var sb strings.Builder
	sb.WriteString(prompt)
	sb.WriteString("\n\n---\n**Injected Context References:**\n")

	injected := make(map[string]bool)

	for _, match := range matches {
		ref := match[1]
		if injected[ref] {
			continue
		}
		injected[ref] = true

		content, err := i.fetchReference(ref)
		if err != nil {
			sb.WriteString(fmt.Sprintf("\n> ⚠️ *Failed to inject `@%s`: %v*\n", ref, err))
			continue
		}

		sb.WriteString(fmt.Sprintf("\n<context source=\"%s\">\n%s\n</context>\n", ref, content))
	}

	return sb.String()
}

func (i *ReferenceInjector) fetchReference(ref string) (string, error) {
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return i.fetchURL(ref)
	}
	return i.fetchFile(ref)
}

func (i *ReferenceInjector) fetchURL(url string) (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Truncate if it's too large (e.g. > 100KB)
	content := string(body)
	if len(content) > 100000 {
		content = content[:100000] + "\n...(truncated)"
	}

	return content, nil
}

func (i *ReferenceInjector) fetchFile(path string) (string, error) {
	// Prevent directory traversal attacks
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid path")
	}

	fullPath := filepath.Join(i.workspaceDir, cleanPath)

	// Check if it's a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("is a directory, not a file")
	}

	// Read file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	text := string(content)
	if len(text) > 100000 {
		text = text[:100000] + "\n...(truncated)"
	}

	return text, nil
}
