package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestReplaceTextSingle(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "config.go")
	original := `package config

type Provider interface {
	Generate() error
}`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:    path,
		OldText: "Provider",
		NewText: "Model",
	}

	res, err := ReplaceText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Replacements != 1 || len(res.Operations) != 1 {
		t.Fatalf("expected 1 replacement, got %+v", res)
	}
	if res.Operations[0].StartLine != 3 {
		t.Errorf("expected match on line 3, got line %d", res.Operations[0].StartLine)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	expected := `package config

type Model interface {
	Generate() error
}`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}

	backupData, err := os.ReadFile(path + ".backup")
	if err != nil {
		t.Fatalf("backup file missing: %v", err)
	}
	if string(backupData) != original {
		t.Errorf("expected backup %q, got %q", original, string(backupData))
	}
}

func TestReplaceTextAmbiguousError(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "ambiguous.go")
	original := `func f1() error { return nil }
func f2() error { return nil }
func f3() error { return nil }`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:    path,
		OldText: "return nil",
		NewText: "return err",
	}

	_, err := ReplaceText(ctx, opts)
	if err == nil {
		t.Fatal("expected error for ambiguous replacement without replace_all or max_replacements, got nil")
	}
}

func TestReplaceTextReplaceAll(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "all.go")
	original := `func f1() error { return nil }
func f2() error { return nil }
func f3() error { return nil }`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:       path,
		OldText:    "return nil",
		NewText:    "return err",
		ReplaceAll: true,
	}

	res, err := ReplaceText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Replacements != 3 || len(res.Operations) != 3 {
		t.Fatalf("expected 3 operations, got %+v", res)
	}
	if res.Operations[0].StartLine != 1 || res.Operations[1].StartLine != 2 || res.Operations[2].StartLine != 3 {
		t.Errorf("unexpected operation lines: %+v", res.Operations)
	}

	data, _ := os.ReadFile(path)
	expected := `func f1() error { return err }
func f2() error { return err }
func f3() error { return err }`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestReplaceTextMaxReplacements(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "max.go")
	original := `func f1() error { return nil }
func f2() error { return nil }
func f3() error { return nil }`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:            path,
		OldText:         "return nil",
		NewText:         "return err",
		MaxReplacements: 2,
	}

	res, err := ReplaceText(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Replacements != 2 || len(res.Operations) != 2 {
		t.Fatalf("expected 2 replacements, got %+v", res)
	}

	data, _ := os.ReadFile(path)
	expected := `func f1() error { return err }
func f2() error { return err }
func f3() error { return nil }`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestReplaceTextNotFound(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "missing.go")
	if err := os.WriteFile(path, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:    path,
		OldText: "Provider",
		NewText: "Model",
	}

	_, err := ReplaceText(ctx, opts)
	if err == nil {
		t.Fatal("expected error when old_text is not found, got nil")
	}
}

func TestReplaceTextBinary(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "binary.bin")
	if err := os.WriteFile(path, []byte{0x00, 0x01, 'a', 'b', 'c'}, 0644); err != nil {
		t.Fatal(err)
	}

	opts := ReplaceOptions{
		File:    path,
		OldText: "abc",
		NewText: "xyz",
	}

	_, err := ReplaceText(ctx, opts)
	if err == nil {
		t.Fatal("expected error when attempting to replace in binary file, got nil")
	}
}

func TestReplaceTextToolRun(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()
	path := filepath.Join(tempDir, "tool.go")
	if err := os.WriteFile(path, []byte("package old"), 0644); err != nil {
		t.Fatal(err)
	}

	tool := NewReplaceTextTool()
	call := tools.Call{
		Name: "replace_text",
		Args: map[string]any{
			"file":     path,
			"old_text": "old",
			"new_text": "new",
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	res, ok := result.Output.(*ReplaceResult)
	if !ok {
		t.Fatalf("expected *ReplaceResult output, got %T", result.Output)
	}
	if res.Replacements != 1 {
		t.Fatalf("expected 1 replacement via tool, got %d", res.Replacements)
	}
}

func TestReplaceTextDefinition(t *testing.T) {
	tool := NewReplaceTextTool()
	def := tool.Definition()

	if def.Name != "replace_text" {
		t.Errorf("unexpected tool name: %s", def.Name)
	}
	if len(def.Parameters) < 3 {
		t.Fatalf("unexpected parameter count: %d", len(def.Parameters))
	}
}
