package files

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

// TestCreateFile - check the wether it is creating the file properly
func TestCreateFile(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	call := tools.Call{
		Name: "create_file",
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

// TestCreateExistingFile - check wether File Already Exists functions works properly
func TestCreateExistingFile(t *testing.T) {
	tool := NewCreateTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	call := tools.Call{
		Name: "create_file",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected error when file already exists")
	}
}

// TestCreateFileInvalidPath - check wether function properly handles Invalid Path
func TestCreateFileInvalidPath(t *testing.T) {
	tool := NewCreateTool()

	call := tools.Call{
		Name: "create_file",
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
		Name: "create_file",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected missing directory error")
	}
}
