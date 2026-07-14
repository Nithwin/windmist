package editing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func createTestNumberedFile(t *testing.T, dir string, lines int) string {
	path := filepath.Join(dir, "testfile.go")
	var sb strings.Builder
	for i := 1; i <= lines; i++ {
		sb.WriteString(fmt.Sprintf("line content %d\n", i))
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	return path
}

func TestReadContext(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := createTestNumberedFile(t, tempDir, 20)

	output, err := ReadContext(ctx, path, 10, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `8: line content 8
9: line content 9
10: line content 10
11: line content 11
12: line content 12
`
	if output != expected {
		t.Errorf("expected output:\n%q\ngot:\n%q", expected, output)
	}
}

func TestReadContextBoundary(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := createTestNumberedFile(t, tempDir, 20)

	output, err := ReadContext(ctx, path, 2, 5, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `1: line content 1
2: line content 2
3: line content 3
4: line content 4
`
	if output != expected {
		t.Errorf("expected output:\n%q\ngot:\n%q", expected, output)
	}
}

func TestReadRange(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := createTestNumberedFile(t, tempDir, 20)

	output, err := ReadRange(ctx, path, 15, 17)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `15: line content 15
16: line content 16
17: line content 17
`
	if output != expected {
		t.Errorf("expected output:\n%q\ngot:\n%q", expected, output)
	}
}

func TestReadContextOutOfRange(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := createTestNumberedFile(t, tempDir, 10)

	_, err := ReadContext(ctx, path, 100, 2, 2)
	if err == nil {
		t.Fatal("expected error when target line exceeds file length, got nil")
	}
}

func TestReadContextToolRun(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := createTestNumberedFile(t, tempDir, 15)
	tool := NewReadContextTool()

	call := tools.Call{
		Name: "read_context",
		Args: map[string]any{
			"path":   path,
			"line":   5,
			"before": 1,
			"after":  1,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	output, ok := result.Output.(string)
	if !ok {
		t.Fatalf("expected string output, got %T", result.Output)
	}

	expected := `4: line content 4
5: line content 5
6: line content 6
`
	if output != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestReadContextDefinition(t *testing.T) {
	tool := NewReadContextTool()
	def := tool.Definition()

	if def.Name != "read_context" {
		t.Errorf("unexpected name: %s", def.Name)
	}
	if len(def.Parameters) < 2 {
		t.Fatalf("unexpected parameters count: %d", len(def.Parameters))
	}
	if def.Parameters[0].Name != "path" || !def.Parameters[0].Required {
		t.Error("expected path to be required")
	}
	if def.Parameters[1].Name != "line" || !def.Parameters[1].Required {
		t.Error("expected line to be required")
	}
}
