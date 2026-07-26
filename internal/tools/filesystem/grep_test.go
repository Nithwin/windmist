package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestGrepTool(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windmist_grep_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("hello world\nthis is a test\nend"), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.go"), []byte("package main\nfunc hello() {}\n"), 0644)

	tool := NewGrepTool()

	res := tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"pattern": "hello",
			"path":    tempDir,
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}

	output, ok := res.Output.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", res.Output)
	}

	matches := output["matches"].([]GrepMatch)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	// Test include glob
	res = tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"pattern": "hello",
			"path":    tempDir,
			"include": "*.go",
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}

	output = res.Output.(map[string]interface{})
	matches = output["matches"].([]GrepMatch)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
}
