package files

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

// TestReadFile - check wether it reads file properly
func TestReadFile(t *testing.T) {
	tool := NewReadTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	call := tools.Call{
		Name: "read_file",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	if result.Output != "hello" {
		t.Fatalf("expected 'hello', got %s", result.Output)
	}
}

// TestReadFileNonExistent - check wether it handles file not found error properly
func TestReadFileNonExistent(t *testing.T) {
	tool := NewReadTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "nonexistent.txt")

	call := tools.Call{
		Name: "read_file",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatalf("expected error when file does not exist")
	}
}

// TestReadFileInvalidPath - check wether it handles invalid path error properly
func TestReadFileInvalidPath(t *testing.T) {
	tool := NewReadTool()


	call := tools.Call{
		Name: "read_file",
		Args: map[string]any{
			"path": "",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatalf("expected error when directory is invalid")
	}
}