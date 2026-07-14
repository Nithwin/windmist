package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestListDirectoryNonRecursive(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewListTool()
	ctx := context.Background()

	// Create files and directory inside tempDir
	if err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create file1.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "file2.go"), []byte("package main"), 0644); err != nil {
		t.Fatalf("failed to create file2.go: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tempDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	call := tools.Call{
		Name: "list",
		Args: map[string]any{
			"path": tempDir,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	entries, ok := result.Output.([]Entry)
	if !ok {
		t.Fatalf("expected []Entry output, got %T", result.Output)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d: %+v", len(entries), entries)
	}

	entryMap := make(map[string]string)
	for _, e := range entries {
		entryMap[e.Name] = e.Type
	}

	if entryMap["file1.txt"] != "file" {
		t.Errorf("expected file1.txt to be type 'file', got %q", entryMap["file1.txt"])
	}
	if entryMap["file2.go"] != "file" {
		t.Errorf("expected file2.go to be type 'file', got %q", entryMap["file2.go"])
	}
	if entryMap["subdir"] != "directory" {
		t.Errorf("expected subdir to be type 'directory', got %q", entryMap["subdir"])
	}
}

func TestListDirectoryRecursive(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewListTool()
	ctx := context.Background()

	// Create nested structure:
	// tempDir/a.txt
	// tempDir/subdir/b.txt
	// tempDir/subdir/nested/c.txt
	if err := os.WriteFile(filepath.Join(tempDir, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatalf("failed to create a.txt: %v", err)
	}
	subdir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("failed to create b.txt: %v", err)
	}
	nested := filepath.Join(subdir, "nested")
	if err := os.Mkdir(nested, 0755); err != nil {
		t.Fatalf("failed to create nested: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nested, "c.txt"), []byte("c"), 0644); err != nil {
		t.Fatalf("failed to create c.txt: %v", err)
	}

	call := tools.Call{
		Name: "list",
		Args: map[string]any{
			"path":      tempDir,
			"recursive": true,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	entries, ok := result.Output.([]Entry)
	if !ok {
		t.Fatalf("expected []Entry output, got %T", result.Output)
	}

	entryMap := make(map[string]string)
	for _, e := range entries {
		entryMap[e.Name] = e.Type
	}

	expectedPaths := map[string]string{
		filepath.Join(tempDir, "a.txt"):                "file",
		filepath.Join(tempDir, "subdir"):               "directory",
		filepath.Join(tempDir, "subdir", "b.txt"):      "file",
		filepath.Join(tempDir, "subdir", "nested"):     "directory",
		filepath.Join(tempDir, "subdir", "nested", "c.txt"): "file",
	}

	if len(entries) != len(expectedPaths) {
		t.Fatalf("expected %d entries, got %d: %+v", len(expectedPaths), len(entries), entries)
	}

	for expectedName, expectedType := range expectedPaths {
		if gotType, exists := entryMap[expectedName]; !exists || gotType != expectedType {
			t.Errorf("expected entry %q with type %q, got exists=%v type=%q", expectedName, expectedType, exists, gotType)
		}
	}
}

func TestListNonExistentDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewListTool()
	ctx := context.Background()

	call := tools.Call{
		Name: "list",
		Args: map[string]any{
			"path": filepath.Join(tempDir, "nonexistent"),
		},
	}

	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error for non-existent directory")
	}
}

func TestListPathNotDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewListTool()
	ctx := context.Background()

	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	call := tools.Call{
		Name: "list",
		Args: map[string]any{
			"path": filePath,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error when listing a regular file path")
	}
}

func TestListEmptyPath(t *testing.T) {
	tool := NewListTool()
	ctx := context.Background()

	call := tools.Call{
		Name: "list",
		Args: map[string]any{
			"path": "",
		},
	}

	result := tool.Run(ctx, call)
	if result.Error == nil {
		t.Error("expected error for empty path")
	}
}

func TestListDefinition(t *testing.T) {
	tool := NewListTool()
	def := tool.Definition()

	if def.Name != "list" {
		t.Errorf("unexpected name: %s", def.Name)
	}
	if def.Description == "" {
		t.Error("empty description")
	}
	if len(def.Parameters) != 2 {
		t.Errorf("unexpected number of parameters: %d", len(def.Parameters))
	}
}
