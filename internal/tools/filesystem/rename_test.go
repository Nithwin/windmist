package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestRenameFile(t *testing.T) {
	tool := NewRenameTool()

	tempDir := t.TempDir()

	oldPath := filepath.Join(tempDir, "old.txt")
	newPath := filepath.Join(tempDir, "new.txt")

	if err := os.WriteFile(oldPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create original file: %v", err)
	}

	call := tools.Call{
		Name: "rename",
		Args: map[string]any{
			"old_path": oldPath,
			"new_path": newPath,
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	// Verify the file was renamed
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Fatalf("expected new file to exist, but it doesn't: %v", err)
	}

	// Verify the old file no longer exists
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old file to not exist, but it does: %v", err)
	}
}

func TestRenameFileInvalidOldPath(t *testing.T) {
	tool := NewRenameTool()

	call := tools.Call{
		Name: "rename",
		Args: map[string]any{
			"old_path": "",
			"new_path": "new.txt",
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error == nil {
		t.Fatalf("expected error for invalid old_path, got nil")
	}
}

func TestRenameFileInvalidNewPath(t *testing.T) {
	tool := NewRenameTool()

	call := tools.Call{
		Name: "rename",
		Args: map[string]any{
			"old_path": "old.txt",
			"new_path": "",
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error == nil {
		t.Fatalf("expected error for invalid new_path, got nil")
	}
}

func TestRenameNonExistentFile(t *testing.T) {
	tool := NewRenameTool()

	tempDir := t.TempDir()
	oldPath := filepath.Join(tempDir, "non_existent.txt")
	newPath := filepath.Join(tempDir, "new.txt")

	call := tools.Call{
		Name: "rename",
		Args: map[string]any{
			"old_path": oldPath,
			"new_path": newPath,
		},
	}

	result := tool.Run(context.Background(), call)
	if result.Error == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}
}