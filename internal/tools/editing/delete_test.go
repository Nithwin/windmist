package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestDeleteRangeSingleLine(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "single.go")
	original := `line 1
delete me
line 3`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := DeleteOptions{
		File:      path,
		StartLine: 2,
		EndLine:   2,
	}

	res, err := DeleteRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 2 || res.Operation.EndLine != 2 {
		t.Errorf("expected deleted line 2, got %+v", res.Operation)
	}
	if res.Operation.OldText != "delete me" || res.Operation.NewText != "" {
		t.Errorf("unexpected old/new text: %+v", res.Operation)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
line 3`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestDeleteRangeMultiLine(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "multi.go")
	original := `line 1
delete 2
delete 3
line 4`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := DeleteOptions{
		File:      path,
		StartLine: 2,
		EndLine:   3,
	}

	res, err := DeleteRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 2 || res.Operation.EndLine != 3 {
		t.Errorf("unexpected operation range: %+v", res.Operation)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
line 4`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestDeleteRangeAllLines(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "all.go")
	original := `line 1
line 2
line 3`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := DeleteOptions{
		File:      path,
		StartLine: 1,
		EndLine:   3,
	}

	res, err := DeleteRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "" {
		t.Errorf("expected empty file after deleting all lines, got %q", string(data))
	}
	if res.Operation.StartLine != 1 || res.Operation.EndLine != 3 {
		t.Errorf("unexpected range: %+v", res.Operation)
	}
}

func TestDeleteRangeInvalidCoordinates(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "inv.go")
	if err := os.WriteFile(path, []byte("1\n2"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := DeleteRange(ctx, DeleteOptions{File: path, StartLine: 0, EndLine: 1})
	if err == nil {
		t.Error("expected error when start_line < 1, got nil")
	}

	_, err = DeleteRange(ctx, DeleteOptions{File: path, StartLine: 2, EndLine: 1})
	if err == nil {
		t.Error("expected error when end_line < start_line, got nil")
	}

	_, err = DeleteRange(ctx, DeleteOptions{File: path, StartLine: 1, EndLine: 5})
	if err == nil {
		t.Error("expected error when end_line > totalLines, got nil")
	}
}

func TestDeleteRangeBinary(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "bin.bin")
	if err := os.WriteFile(path, []byte{0x00, 0x01}, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := DeleteRange(ctx, DeleteOptions{File: path, StartLine: 1, EndLine: 1})
	if err == nil {
		t.Fatal("expected error when deleting from binary file, got nil")
	}
}

func TestDeleteRangeToolRun(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "tool.go")
	if err := os.WriteFile(path, []byte("line 1\nline 2\nline 3"), 0644); err != nil {
		t.Fatal(err)
	}

	tool := NewDeleteRangeTool()
	call := tools.Call{
		Name: "delete_range",
		Args: map[string]any{
			"file":       path,
			"start_line": 2,
			"end_line":   2,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	res, ok := result.Output.(*DeleteResult)
	if !ok {
		t.Fatalf("expected *DeleteResult output, got %T", result.Output)
	}
	if res.Operation.StartLine != 2 || res.Operation.EndLine != 2 {
		t.Fatalf("expected deleted range 2-2, got %+v", res.Operation)
	}
}

func TestDeleteRangeDefinition(t *testing.T) {
	tool := NewDeleteRangeTool()
	def := tool.Definition()

	if def.Name != "delete_range" {
		t.Errorf("unexpected tool name: %s", def.Name)
	}
	if len(def.Parameters) < 3 {
		t.Fatalf("unexpected parameter count: %d", len(def.Parameters))
	}
}
