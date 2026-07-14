package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestReplaceRangeSingleLine(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "single.go")
	original := `line 1
func old() {
line 3`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceRangeOptions{
		File:      path,
		StartLine: 2,
		EndLine:   2,
		NewText:   "func new() {",
	}

	res, err := ReplaceRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 2 || res.Operation.EndLine != 2 {
		t.Errorf("unexpected operation range: %+v", res.Operation)
	}
	if res.Operation.OldText != "func old() {" || res.Operation.NewText != "func new() {" {
		t.Errorf("unexpected texts: %+v", res.Operation)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
func new() {
line 3`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestReplaceRangeMultiLine(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "multi.go")
	original := `line 1
old 2
old 3
line 4`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceRangeOptions{
		File:      path,
		StartLine: 2,
		EndLine:   3,
		NewText:   "new 2\nnew 3\nnew 4",
	}

	res, err := ReplaceRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Operation.StartLine != 2 || res.Operation.EndLine != 4 {
		t.Errorf("expected end line 4 (2 + 3 lines - 1), got %d", res.Operation.EndLine)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
new 2
new 3
new 4
line 4`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestReplaceRangeDeleteLines(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "delete.go")
	original := `line 1
delete me 1
delete me 2
line 4`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceRangeOptions{
		File:      path,
		StartLine: 2,
		EndLine:   3,
		NewText:   "",
	}

	res, err := ReplaceRange(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	expected := `line 1
line 4`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
	if res.Operation.EndLine != 1 {
		t.Errorf("expected end line 1 after deletion, got %d", res.Operation.EndLine)
	}
}

func TestReplaceRangeFailFastClamping(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "clamp.go")
	if err := os.WriteFile(path, []byte("1\n2\n3\n4\n5"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceRangeOptions{
		File:      path,
		StartLine: 4,
		EndLine:   10,
		NewText:   "new",
	}

	_, err := ReplaceRange(ctx, opts)
	if err == nil {
		t.Fatal("expected fail-fast error when end_line exceeds file length, got nil")
	}
}

func TestReplaceRangeInvalidCoordinates(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "invalid.go")
	if err := os.WriteFile(path, []byte("1\n2\n3"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ReplaceRange(ctx, ReplaceRangeOptions{File: path, StartLine: 0, EndLine: 2})
	if err == nil {
		t.Error("expected error when start_line < 1, got nil")
	}

	_, err = ReplaceRange(ctx, ReplaceRangeOptions{File: path, StartLine: 3, EndLine: 1})
	if err == nil {
		t.Error("expected error when end_line < start_line, got nil")
	}
}

func TestReplaceRangeBinaryFile(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "binary.bin")
	if err := os.WriteFile(path, []byte{0x00, 0x01, 'a', 'b'}, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ReplaceRange(ctx, ReplaceRangeOptions{File: path, StartLine: 1, EndLine: 1, NewText: "x"})
	if err == nil {
		t.Fatal("expected error on binary file range replacement, got nil")
	}
}

func TestReplaceRangeToolRun(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "tool.go")
	if err := os.WriteFile(path, []byte("line 1\nline 2\nline 3"), 0644); err != nil {
		t.Fatal(err)
	}

	tool := NewReplaceRangeTool()
	call := tools.Call{
		Name: "replace_range",
		Args: map[string]any{
			"file":       path,
			"start_line": 2,
			"end_line":   2,
			"new_text":   "updated 2",
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	res, ok := result.Output.(*ReplaceRangeResult)
	if !ok {
		t.Fatalf("expected *ReplaceRangeResult output, got %T", result.Output)
	}
	if res.Operation.StartLine != 2 {
		t.Fatalf("expected start line 2 via tool, got %d", res.Operation.StartLine)
	}
}

func TestReplaceRangeDefinition(t *testing.T) {
	tool := NewReplaceRangeTool()
	def := tool.Definition()

	if def.Name != "replace_range" {
		t.Errorf("unexpected tool name: %s", def.Name)
	}
	if len(def.Parameters) < 4 {
		t.Fatalf("unexpected parameter count: %d", len(def.Parameters))
	}
}
