package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestAppendFile(t *testing.T) {
	tool := NewAppendTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	// Create initial file
	if err := os.WriteFile(path, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	call := tools.Call{
		Name: "append",
		Args: map[string]any{
			"path":    path,
			"content": "more content\n",
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

	expected := "initial content\nmore content\n"
	if string(data) != expected {
		t.Fatalf("expected file content to be %q, got %q", expected, string(data))
	}
}

func TestAppendFileEmptyContent(t *testing.T) {
	tool := NewAppendTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(path, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	call := tools.Call{
		Name: "append",
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

	expected := "initial content\n"
	if string(data) != expected {
		t.Fatalf("expected file content to be %q, got %q", expected, string(data))
	}
}

func TestAppendFileMissingFile(t *testing.T) {
	tool := NewAppendTool()

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "nonexistent.txt")

	call := tools.Call{
		Name: "append",
		Args: map[string]any{
			"path":    path,
			"content": "hello\n",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected error when file does not exist")
	}
}

func TestAppendFileInvalidPath(t *testing.T) {
	tool := NewAppendTool()

	call := tools.Call{
		Name: "append",
		Args: map[string]any{
			"path":    "",
			"content": "hello\n",
		},
	}

	result := tool.Run(context.Background(), call)

	if result.Error == nil {
		t.Fatal("expected invalid path error")
	}
}
