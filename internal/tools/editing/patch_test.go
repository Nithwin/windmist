package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPatchApplyAndUndo(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "patch.go")
	original := `line 1
line 2
line 3`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// First, let's execute operations and build our patch from their returned Operations
	res1, err := ReplaceRange(ctx, ReplaceRangeOptions{
		File:      path,
		StartLine: 2,
		EndLine:   2,
		NewText:   "updated 2",
	})
	if err != nil {
		t.Fatal(err)
	}

	res2, err := InsertText(ctx, InsertOptions{
		File:    path,
		Line:    3,
		NewText: "inserted 2.5",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify the forward modifications took place
	modifiedData, _ := os.ReadFile(path)
	expectedModified := `line 1
updated 2
inserted 2.5
line 3`
	if string(modifiedData) != expectedModified {
		t.Fatalf("expected modified:\n%s\ngot:\n%s", expectedModified, string(modifiedData))
	}

	// Build a patch of what we just applied (`res1` then `res2`)
	patch := NewPatch(res1.Operation, res2.Operation)
	if len(patch.AllOperations()) != 2 {
		t.Fatalf("expected 2 ops, got %d", len(patch.AllOperations()))
	}
	if len(patch.Files) != 1 || patch.Files[0].Path != path {
		t.Fatalf("expected 1 FilePatch grouped by %s, got %+v", path, patch.Files)
	}

	// Now call Undo on the patch (`patch.Reverse().Apply(ctx)`)
	if err := Undo(ctx, patch); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	// Verify that the file is 100% restored to its original content!
	restoredData, _ := os.ReadFile(path)
	if string(restoredData) != original {
		t.Errorf("expected original content after Undo:\n%s\ngot:\n%s", original, string(restoredData))
	}
}

func TestPatchReverseMathematics(t *testing.T) {
	op1 := Operation{
		Type:      Replace,
		File:      "f.go",
		StartLine: 2,
		EndLine:   4,
		OldText:   "old",
		NewText:   "new",
	}
	op2 := Operation{
		Type:      Insert,
		File:      "f.go",
		StartLine: 10,
		EndLine:   11,
		OldText:   "",
		NewText:   "new line 10\nnew line 11",
	}
	op3 := Operation{
		Type:      Delete,
		File:      "f.go",
		StartLine: 20,
		EndLine:   22,
		OldText:   "deleted lines",
		NewText:   "",
	}

	patch := NewPatch(op1, op2, op3)
	rev := patch.Reverse()

	revOps := rev.AllOperations()
	if len(revOps) != 3 {
		t.Fatalf("expected 3 rev ops, got %d", len(revOps))
	}

	// Check reverse order: revOps[0] should be inverse of op3
	r0 := revOps[0]
	if r0.Type != Insert || r0.StartLine != 20 || r0.NewText != "deleted lines" {
		t.Errorf("unexpected rev[0]: %+v", r0)
	}

	// revOps[1] should be inverse of op2
	r1 := revOps[1]
	if r1.Type != Delete || r1.StartLine != 10 || r1.EndLine != 11 || r1.OldText != "new line 10\nnew line 11" {
		t.Errorf("unexpected rev[1]: %+v", r1)
	}

	// revOps[2] should be inverse of op1
	r2 := revOps[2]
	if r2.Type != Replace || r2.StartLine != 2 || r2.EndLine != 4 || r2.OldText != "new" || r2.NewText != "old" {
		t.Errorf("unexpected rev[2]: %+v", r2)
	}
}

func TestPatchApplyFailureRollback(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "fail.go")
	if err := os.WriteFile(path, []byte("1\n2\n3"), 0644); err != nil {
		t.Fatal(err)
	}

	patch := NewPatch(
		Operation{Type: Replace, File: path, StartLine: 1, EndLine: 1, OldText: "1", NewText: "one"},
		Operation{Type: Replace, File: path, StartLine: 9999, EndLine: 9999, OldText: "x", NewText: "y"}, // invalid line
	)

	err := patch.Apply(ctx)
	if err == nil {
		t.Fatal("expected error on invalid operation index, got nil")
	}
}
