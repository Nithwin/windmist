package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestGlobTool(t *testing.T) {
	// Setup test directory
	tempDir, err := os.MkdirTemp("", "windmist_glob_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "test1.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tempDir, "test2.txt"), []byte("hello"), 0644)
	os.Mkdir(filepath.Join(tempDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tempDir, "sub", "test3.go"), []byte("package sub"), 0644)

	tool := NewGlobTool()

	// Test 1: *.go in root
	res := tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"pattern": "*.go",
			"path":    tempDir,
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}

	matches, ok := res.Output.([]string)
	if !ok {
		t.Fatalf("expected []string, got %T", res.Output)
	}
	if len(matches) != 1 || matches[0] != "test1.go" {
		t.Fatalf("expected [test1.go], got %v", matches)
	}

	// Test 2: **/*.go (recursive)
	res = tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"pattern": "**/*.go",
			"path":    tempDir,
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}
	matches = res.Output.([]string)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(matches), matches)
	}
}
