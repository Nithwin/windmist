package editing

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEditingPipelineIntegration tests the full end-to-end AI editing workflow:
// Search -> ReadContext -> ReplaceRange -> InsertText -> DeleteRange -> Verify final file state.
func TestEditingPipelineIntegration(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "service.go")

	initialCode := `package service

import "fmt"

func ProcessData(input string) error {
	if input == "" {
		return fmt.Errorf("empty input")
	}
	// TODO: add validation logic
	fmt.Println("Processing:", input)
	return nil
}`
	if err := os.WriteFile(path, []byte(initialCode), 0644); err != nil {
		t.Fatal(err)
	}

	// Step 1: Search for where "TODO: add validation logic" is located.
	searchOpts := SearchOptions{
		Root:          tempDir,
		Query:         "TODO: add validation logic",
		Type:          SearchText,
		CaseSensitive: true,
	}
	searchResults, err := Search(ctx, searchOpts)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(searchResults) != 1 || len(searchResults[0].Matches) != 1 {
		t.Fatalf("expected 1 match, got %v", searchResults)
	}
	todoLine := searchResults[0].Matches[0].Line
	if todoLine != 9 {
		t.Fatalf("expected TODO at line 9, got %d", todoLine)
	}

	// Step 2: Read context around the TODO line using ReadContext.
	contextOutput, err := ReadContext(ctx, path, todoLine, 2, 2)
	if err != nil {
		t.Fatalf("read_context failed: %v", err)
	}
	if !strings.Contains(contextOutput, "9: \t// TODO: add validation logic") || !strings.Contains(contextOutput, "7: \t\treturn fmt.Errorf(\"empty input\")") {
		t.Fatalf("unexpected context output:\n%s", contextOutput)
	}

	// Step 3: Replace the TODO line (line 9) with real validation using ReplaceRange.
	replaceOpts := ReplaceRangeOptions{
		File:      path,
		StartLine: todoLine,
		EndLine:   todoLine,
		NewText:   "\tif len(input) > 100 {\n\t\treturn fmt.Errorf(\"input too long\")\n\t}",
	}
	replaceRes, err := ReplaceRange(ctx, replaceOpts)
	if err != nil {
		t.Fatalf("replace_range failed: %v", err)
	}
	if replaceRes.Operation.StartLine != 9 || replaceRes.Operation.EndLine != 11 {
		t.Errorf("expected replaced lines 9-11, got %+v", replaceRes.Operation)
	}

	// Step 4: Insert a log entry right before line 12 (formerly line 10 before our 3-line replacement).
	insertOpts := InsertOptions{
		File:    path,
		Line:    12,
		NewText: "\tfmt.Println(\"Validation passed\")",
	}
	insertRes, err := InsertText(ctx, insertOpts)
	if err != nil {
		t.Fatalf("insert_text failed: %v", err)
	}
	if insertRes.Operation.StartLine != 12 || insertRes.Operation.EndLine != 12 {
		t.Errorf("expected inserted line 12, got %+v", insertRes.Operation)
	}

	// Step 5: Delete the old "Processing:" print line (now pushed to line 13).
	deleteOpts := DeleteOptions{
		File:      path,
		StartLine: 13,
		EndLine:   13,
	}
	deleteRes, err := DeleteRange(ctx, deleteOpts)
	if err != nil {
		t.Fatalf("delete_range failed: %v", err)
	}
	if deleteRes.Operation.StartLine != 13 || deleteRes.Operation.EndLine != 13 {
		t.Errorf("expected deleted line 13, got %+v", deleteRes.Operation)
	}

	// Step 6: Verify final code structure and backup preservation.
	finalData, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	expectedCode := `package service

import "fmt"

func ProcessData(input string) error {
	if input == "" {
		return fmt.Errorf("empty input")
	}
	if len(input) > 100 {
		return fmt.Errorf("input too long")
	}
	fmt.Println("Validation passed")
	return nil
}`
	if string(finalData) != expectedCode {
		t.Errorf("unexpected final file state after pipeline:\nExpected:\n%s\nGot:\n%s", expectedCode, string(finalData))
	}

	// Verify backup file exists from the atomic writer
	if _, err := os.Stat(path + ".backup"); err != nil {
		t.Errorf("expected .backup file to be preserved: %v", err)
	}
}

// TestEditingMultiLineInsertionAndCleanup tests inserting multiple helpers and cleaning up via ReplaceRange.
func TestEditingMultiLineInsertionAndCleanup(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "helpers.go")

	initialCode := `package helpers

func Main() {}`
	if err := os.WriteFile(path, []byte(initialCode), 0644); err != nil {
		t.Fatal(err)
	}

	// Insert two helper functions after line 3 (at line 4)
	insertRes, err := InsertText(ctx, InsertOptions{
		File:    path,
		Line:    4,
		NewText: "func HelperA() int { return 1 }\nfunc HelperB() int { return 2 }",
	})
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if insertRes.Operation.StartLine != 4 || insertRes.Operation.EndLine != 5 {
		t.Errorf("expected inserted lines 4-5, got %+v", insertRes.Operation)
	}

	// Replace HelperB with HelperC
	replaceRes, err := ReplaceText(ctx, ReplaceOptions{
		File:    path,
		OldText: "func HelperB() int { return 2 }",
		NewText: "func HelperC() string { return \"C\" }",
	})
	if err != nil {
		t.Fatalf("replace_text failed: %v", err)
	}
	if len(replaceRes.Operations) != 1 || replaceRes.Operations[0].StartLine != 5 {
		t.Errorf("expected replace on line 5, got %+v", replaceRes.Operations)
	}

	finalData, _ := os.ReadFile(path)
	expected := `package helpers

func Main() {}
func HelperA() int { return 1 }
func HelperC() string { return "C" }`
	if string(finalData) != expected {
		t.Errorf("unexpected content:\n%s", string(finalData))
	}
}

// TestEditingRollbackSimulation tests that any failed or aborted workflow can cleanly revert using the automatic .backup file.
func TestEditingRollbackSimulation(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "config.go")

	originalCode := `package config
const Version = "v1.0.0"`
	if err := os.WriteFile(path, []byte(originalCode), 0644); err != nil {
		t.Fatal(err)
	}

	// Perform a risky replacement
	_, err := ReplaceText(ctx, ReplaceOptions{
		File:    path,
		OldText: "v1.0.0",
		NewText: "v2.0.0-broken",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified state
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "v2.0.0-broken") {
		t.Fatalf("expected modified file, got %s", string(data))
	}

	// Simulate rollback using backup
	backupData, err := os.ReadFile(path + ".backup")
	if err != nil {
		t.Fatalf("missing backup: %v", err)
	}
	if err := WriteFile(path, backupData, 0); err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	restoredData, _ := os.ReadFile(path)
	if string(restoredData) != originalCode {
		t.Errorf("expected original code after rollback, got %s", string(restoredData))
	}
}

