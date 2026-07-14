package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestInsertTextBeginning(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "start.go")
	original := `package main
func main() {}`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := InsertOptions{
		File:    path,
		Line:    1,
		NewText: "// comment 1\n// comment 2",
	}

	res, err := InsertText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 1 || res.Operation.EndLine != 2 {
		t.Errorf("expected inserted lines 1-2, got %+v", res.Operation)
	}

	data, _ := os.ReadFile(path)
	expected := `// comment 1
// comment 2
package main
func main() {}`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestInsertTextMiddle(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "mid.go")
	original := `line 1
line 2
line 3`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := InsertOptions{
		File:    path,
		Line:    2,
		NewText: "inserted A\ninserted B",
	}

	res, err := InsertText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 2 || res.Operation.EndLine != 3 {
		t.Errorf("expected end line 3, got %d", res.Operation.EndLine)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
inserted A
inserted B
line 2
line 3`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestInsertTextEnd(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "end.go")
	original := `line 1
line 2`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := InsertOptions{
		File:    path,
		Line:    3,
		NewText: "appended 3",
	}

	res, err := InsertText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
line 2
appended 3`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
	if res.Operation.StartLine != 3 || res.Operation.EndLine != 3 {
		t.Errorf("unexpected operation range: %+v", res.Operation)
	}
}

func TestInsertTextInvalidLine(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "inv.go")
	if err := os.WriteFile(path, []byte("1\n2"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := InsertText(ctx, InsertOptions{File: path, Line: 0, NewText: "x"})
	if err == nil {
		t.Error("expected error when line < 1, got nil")
	}

	_, err = InsertText(ctx, InsertOptions{File: path, Line: 4, NewText: "x"})
	if err == nil {
		t.Error("expected error when line > totalLines + 1, got nil")
	}
}

func TestInsertTextBinary(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "bin.bin")
	if err := os.WriteFile(path, []byte{0x00, 0x01}, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := InsertText(ctx, InsertOptions{File: path, Line: 1, NewText: "x"})
	if err == nil {
		t.Fatal("expected error when inserting into binary file, got nil")
	}
}

func TestInsertTextToolRun(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "tool.go")
	if err := os.WriteFile(path, []byte("line 1\nline 2"), 0644); err != nil {
		t.Fatal(err)
	}

	tool := NewInsertTextTool()
	call := tools.Call{
		Name: "insert_text",
		Args: map[string]any{
			"file":     path,
			"line":     2,
			"new_text": "inserted via tool",
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	res, ok := result.Output.(*InsertResult)
	if !ok {
		t.Fatalf("expected *InsertResult output, got %T", result.Output)
	}
	if res.Operation.StartLine != 2 {
		t.Fatalf("expected start line 2, got %d", res.Operation.StartLine)
	}
}

func TestInsertTextDefinition(t *testing.T) {
	tool := NewInsertTextTool()
	def := tool.Definition()

	if def.Name != "insert_text" {
		t.Errorf("unexpected tool name: %s", def.Name)
	}
	if len(def.Parameters) < 3 {
		t.Fatalf("unexpected parameter count: %d", len(def.Parameters))
	}
}
