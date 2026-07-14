package editing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicNew(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "newfile.txt")
	content := []byte("hello world")

	if err := WriteFile(path, content, 0644); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(data))
	}
}

func TestWriteFileAtomicBackup(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "existing.go")
	if err := os.WriteFile(path, []byte("package original"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := WriteFile(path, []byte("package modified"), 0644); err != nil {
		t.Fatalf("unexpected error on overwrite: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "package modified" {
		t.Errorf("expected modified content, got %q", string(data))
	}

	backupPath := path + ".backup"
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("backup file not created: %v", err)
	}
	if string(backupData) != "package original" {
		t.Errorf("expected backup to contain original content, got %q", string(backupData))
	}
}
