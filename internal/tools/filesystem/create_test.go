package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

// TestCreateFile - check whether it creates the file properly
func TestCreateFile(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

// TestCreateDirectory - check whether it creates a directory properly
func TestCreateDirectory(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "testdir")

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": path,
			"type": "directory",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected path to be a directory")
	}
}

// TestCreateExistingFile - check whether File Already Exists check works properly
func TestCreateExistingFile(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected error when file already exists")
	}
}

// TestCreateExistingDirectory - check whether Directory Already Exists check works properly
func TestCreateExistingDirectory(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "testdir")

	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": path,
			"type": "directory",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected error when directory already exists")
	}
}

// TestCreateFileInvalidPath - check whether function properly handles Invalid Path
func TestCreateFileInvalidPath(t *testing.T) {
	tool := NewCreateTool()

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": "",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected invalid path error")
	}
}

func TestCreateFileMissingDirectory(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "missing-folder", "test.txt")

	call := tools.Call{
		Name: "create",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected missing directory error")
	}
}
