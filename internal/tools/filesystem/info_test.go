package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestInfoFile(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewInfoTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "file.txt")
	content := "hello world"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	call := tools.Call{
		Name: "info",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(InfoResult)
	if !ok {
		t.Fatalf("expected InfoResult output, got %T", result.Output)
	}

	if output.Name != "file.txt" {
		t.Errorf("expected name 'file.txt', got %q", output.Name)
	}
	if output.Path != path {
		t.Errorf("expected path %q, got %q", path, output.Path)
	}
	if output.Size != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), output.Size)
	}
	if output.Type != "file" {
		t.Errorf("expected type 'file', got %q", output.Type)
	}
	if output.Mode == "" {
		t.Error("expected non-empty mode string")
	}
	if output.LastModified.IsZero() {
		t.Error("expected non-zero last modified time")
	}
}

func TestInfoDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewInfoTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	call := tools.Call{
		Name: "info",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(InfoResult)
	if !ok {
		t.Fatalf("expected InfoResult output, got %T", result.Output)
	}

	if output.Name != "testdir" {
		t.Errorf("expected name 'testdir', got %q", output.Name)
	}
	if output.Path != path {
		t.Errorf("expected path %q, got %q", path, output.Path)
	}
	if output.Type != "directory" {
		t.Errorf("expected type 'directory', got %q", output.Type)
	}
	if output.Mode == "" {
		t.Error("expected non-empty mode string")
	}
	if output.LastModified.IsZero() {
		t.Error("expected non-zero last modified time")
	}
}

func TestInfoMissingPath(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewInfoTool()
	ctx := context.Background()

	path := filepath.Join(tempDir, "nonexistent.txt")

	call := tools.Call{
		Name: "info",
		Args: map[string]any{
			"path": path,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestInfoInvalidPath(t *testing.T) {
	tool := NewInfoTool()
	ctx := context.Background()

	call := tools.Call{
		Name: "info",
		Args: map[string]any{
			"path": "",
		},
	}

	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error for empty path")
	}
}

func TestInfoDefinition(t *testing.T) {
	tool := NewInfoTool()
	def := tool.Definition()

	if def.Name != "info" {
		t.Errorf("unexpected name: %s", def.Name)
	}
	if def.Description == "" {
		t.Error("empty description")
	}
	if len(def.Parameters) != 1 {
		t.Fatalf("unexpected number of parameters: %d", len(def.Parameters))
	}
	if def.Parameters[0].Name != "path" {
		t.Errorf("unexpected parameter name: %s", def.Parameters[0].Name)
	}
	if !def.Parameters[0].Required {
		t.Error("expected path parameter to be required")
	}
}
