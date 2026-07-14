package editing

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFile performs an atomic write to disk by writing to a temporary file and atomically renaming it.
// It also creates a backup of the original file (.backup) if it already exists, ensuring rollback safety.
func WriteFile(path string, content []byte, perm os.FileMode) error {
	info, err := os.Stat(path)
	if err == nil {
		if perm == 0 {
			perm = info.Mode()
		}
		// Create backup before modifying
		backupPath := path + ".backup"
		originalData, readErr := os.ReadFile(path)
		if readErr == nil {
			_ = os.WriteFile(backupPath, originalData, perm)
		}
	} else if perm == 0 {
		perm = 0644
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %q: %w", dir, err)
	}

	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary file in %q: %w", dir, err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(content); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to write content to temporary file %q: %w", tmpPath, err)
	}
	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions on temporary file %q: %w", tmpPath, err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to close temporary file %q: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically rename %q to %q: %w", tmpPath, path, err)
	}

	return nil
}
