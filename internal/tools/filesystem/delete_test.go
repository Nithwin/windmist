package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestDeleteFile(t *testing.T) {
	tool := NewDeleteTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	// Create a file to delete
	if err := os.WriteFile(path, []byte("delete me"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	call := tools.Call{
		Name: "delete",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	// Verify the file was deleted
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file to not exist, but it does: %v", err)
	}
}

func TestDeleteDirectory(t *testing.T) {
	tool := NewDeleteTool()

	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "testdir")

	// Create a directory to delete
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	call := tools.Call{
		Name: "delete",
		Args: map[string]any{
			"path": dirPath,
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	// Verify the directory was deleted
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		t.Fatalf("expected directory to not exist, but it does: %v", err)
	}
}

func TestDeleteNonExistentFile(t *testing.T) {
	tool := NewDeleteTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "non_existent.txt")

	call := tools.Call{
		Name: "delete",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}
}

func TestDeleteEmptyPath(t *testing.T) {
	tool := NewDeleteTool()

	call := tools.Call{
		Name: "delete",
		Args: map[string]any{
			"path": "",
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error == nil {
		t.Fatalf("expected error for empty path, got nil")
	}
}
