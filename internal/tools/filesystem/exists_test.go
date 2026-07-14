package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestExistsFile(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewExistsTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	call := tools.Call{
		Name: "exists",
		Args: map[string]any{
			"path": path,
		},
	}
	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(ExistsResult)
	if !ok {
		t.Fatalf("expected ExistsResult output, got %T", result.Output)
	}

	if !output.Exists {
		t.Errorf("expected exists to be true, got %v", output.Exists)
	}

	if output.Type != "file" {
		t.Errorf("expected type to be 'file', got %v", output.Type)
	}
}

func TestExistsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewExistsTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	call := tools.Call{
		Name: "exists",
		Args: map[string]any{
			"path": path,
		},
	}
	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(ExistsResult)
	if !ok {
		t.Fatalf("expected ExistsResult output, got %T", result.Output)
	}

	if !output.Exists {
		t.Errorf("expected exists to be true, got %v", output.Exists)
	}

	if output.Type != "directory" {
		t.Errorf("expected type to be 'directory', got %v", output.Type)
	}
}

func TestDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewExistsTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "nonexistent.txt")

	call := tools.Call{
		Name: "exists",
		Args: map[string]any{
			"path": path,
		},
	}
	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(ExistsResult)
	if !ok {
		t.Fatalf("expected ExistsResult output, got %T", result.Output)
	}

	if output.Exists {
		t.Errorf("expected exists to be false, got %v", output.Exists)
	}

	if output.Type != "" {
		t.Errorf("expected no 'type' key for non-existent path, got %v", output.Type)
	}
}

func TestInvalidPath(t *testing.T) {
	tool := NewExistsTool()
	ctx := context.Background()

	call := tools.Call{
		Name: "exists",
		Args: map[string]any{
			"path": "",
		},
	}
	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDefinition(t *testing.T) {
	tool := NewExistsTool()
	def := tool.Definition()

	if def.Name != "exists" {
		t.Errorf("unexpected name: %s", def.Name)
	}
	if def.Description == "" {
		t.Error("empty description")
	}
	if len(def.Parameters) != 1 {
		t.Errorf("unexpected number of parameters: %d", len(def.Parameters))
	}
}
