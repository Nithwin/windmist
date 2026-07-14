package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

// TestWriteFile - check wether it writes file properly
func TestWriteFile(t *testing.T) {
	tool := NewWriteTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	call := tools.Call{
		Name: "write",
		Args: map[string]any{
			"path":    path,
			"content": "hello",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	if string(data) != "hello" {
		t.Fatalf("expected file content to be %q, got %q", "hello", string(data))
	}
}

// TestWriteFileEmptyContent - check wether it handles empty content properly
func TestWriteFileEmptyContent(t *testing.T) {
	tool := NewWriteTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	call := tools.Call{
		Name: "write",
		Args: map[string]any{
			"path":    path,
			"content": "",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	if string(data) != "" {
		t.Fatalf("expected file content to be empty, got %q", string(data))
	}
}

// TestWriteFileEmptyPath - check wether it handles empty path error properly
func TestWriteFileEmptyPath(t *testing.T) {
	tool := NewWriteTool()

	call := tools.Call{
		Name: "write",
		Args: map[string]any{
			"path": "",
			"content": "hello",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatalf("expected error when path is empty")
	}
}

// TestWriteFileInvalidPath - check wether it handles invalid path error properly
func TestWriteFileInvalidPath(t *testing.T) {
	tool := NewWriteTool()

	call := tools.Call{
		Name: "write",
		Args: map[string]any{
			"path": "/invalid/directory/test.txt",
			"content": "hello",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatalf("expected error when directory does not exist")
	}
}