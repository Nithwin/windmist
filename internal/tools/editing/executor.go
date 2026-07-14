package editing

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Executor is responsible for applying or undoing Patches against target files or memory buffers.
type Executor struct{}

// NewExecutor creates a new Executor instance.
func NewExecutor() *Executor {
	return &Executor{}
}

// Apply executes every FilePatch sequentially inside the provided patch against disk.
// For each file, content is read exactly once, all operations are sorted in descending order
// of StartLine and applied in memory, and the result is written to disk once using atomic file writing.
func (e *Executor) Apply(ctx context.Context, p *Patch) error {
	if p == nil || len(p.Files) == 0 {
		return nil
	}

	for _, fp := range p.Files {
		if len(fp.Operations) == 0 {
			continue
		}

		// Handle string-based replacements where line coordinates are not tracked (StartLine == 0)
		if len(fp.Operations) == 1 && fp.Operations[0].StartLine == 0 && fp.Operations[0].EndLine == 0 && fp.Operations[0].Type == Replace {
			opts := ReplaceOptions{
				File:       fp.Path,
				OldText:    fp.Operations[0].OldText,
				NewText:    fp.Operations[0].NewText,
				ReplaceAll: true,
			}
			if _, err := ReplaceText(ctx, opts); err != nil {
				return fmt.Errorf("executor apply failed on %s: %w", fp.Path, err)
			}
			continue
		}

		f, err := os.Open(fp.Path)
		if err != nil {
			return fmt.Errorf("executor apply failed to open %s: %w", fp.Path, err)
		}
		binary, binErr := isBinaryFile(f)
		_ = f.Close()
		if binErr != nil {
			return fmt.Errorf("executor apply failed checking binary state on %s: %w", fp.Path, binErr)
		}
		if binary {
			return fmt.Errorf("executor apply failed: cannot perform line editing on binary file %q", fp.Path)
		}

		content, err := os.ReadFile(fp.Path)
		if err != nil {
			return fmt.Errorf("executor apply failed to read %s: %w", fp.Path, err)
		}

		lines := []string{}
		hasTrailingNewline := false
		if len(content) > 0 {
			hasTrailingNewline = strings.HasSuffix(string(content), "\n")
			lines = strings.Split(string(content), "\n")
			if hasTrailingNewline && len(lines) > 1 && lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
		}

		opsCopy := make([]Operation, len(fp.Operations))
		copy(opsCopy, fp.Operations)
		sort.Slice(opsCopy, func(i, j int) bool {
			if opsCopy[i].StartLine == opsCopy[j].StartLine {
				return opsCopy[i].EndLine > opsCopy[j].EndLine
			}
			return opsCopy[i].StartLine > opsCopy[j].StartLine
		})

		for _, op := range opsCopy {
			totalLines := len(lines)
			switch op.Type {
			case Replace:
				if op.StartLine < 1 || op.EndLine < op.StartLine || op.StartLine > totalLines || op.EndLine > totalLines {
					return fmt.Errorf("executor apply failed on %s: invalid replace coordinates [%d, %d] (file has %d lines)", fp.Path, op.StartLine, op.EndLine, totalLines)
				}
				oldSlice := lines[op.StartLine-1 : op.EndLine]
				var newSlice []string
				if op.NewText != "" {
					newSlice = strings.Split(strings.TrimSuffix(op.NewText, "\n"), "\n")
				}
				updated := make([]string, 0, len(lines)-len(oldSlice)+len(newSlice))
				updated = append(updated, lines[:op.StartLine-1]...)
				if len(newSlice) > 0 {
					updated = append(updated, newSlice...)
				}
				updated = append(updated, lines[op.EndLine:]...)
				lines = updated

			case Insert:
				if op.StartLine < 1 || op.StartLine > totalLines+1 {
					return fmt.Errorf("executor apply failed on %s: invalid insert line %d (file has %d lines)", fp.Path, op.StartLine, totalLines)
				}
				if op.NewText == "" {
					continue
				}
				newSlice := strings.Split(strings.TrimSuffix(op.NewText, "\n"), "\n")
				if len(newSlice) == 0 {
					continue
				}
				updated := make([]string, 0, len(lines)+len(newSlice))
				if op.StartLine <= totalLines {
					updated = append(updated, lines[:op.StartLine-1]...)
					updated = append(updated, newSlice...)
					updated = append(updated, lines[op.StartLine-1:]...)
				} else {
					updated = append(updated, lines...)
					updated = append(updated, newSlice...)
				}
				lines = updated

			case Delete:
				if op.StartLine < 1 || op.EndLine < op.StartLine || op.StartLine > totalLines || op.EndLine > totalLines {
					return fmt.Errorf("executor apply failed on %s: invalid delete coordinates [%d, %d] (file has %d lines)", fp.Path, op.StartLine, op.EndLine, totalLines)
				}
				oldSlice := lines[op.StartLine-1 : op.EndLine]
				updated := make([]string, 0, len(lines)-len(oldSlice))
				updated = append(updated, lines[:op.StartLine-1]...)
				updated = append(updated, lines[op.EndLine:]...)
				lines = updated

			default:
				return fmt.Errorf("executor apply failed on %s: unsupported operation type %v", fp.Path, op.Type)
			}
		}

		updatedContent := strings.Join(lines, "\n")
		if hasTrailingNewline && (len(updatedContent) == 0 || !strings.HasSuffix(updatedContent, "\n")) {
			updatedContent += "\n"
		}

		if err := WriteFile(fp.Path, []byte(updatedContent), 0); err != nil {
			return fmt.Errorf("executor apply failed to write %s: %w", fp.Path, err)
		}
	}

	return nil
}

// Undo applies the exact mathematical inverse of the provided patch using this executor.
func (e *Executor) Undo(ctx context.Context, p *Patch) error {
	if p == nil {
		return nil
	}
	return e.Apply(ctx, p.Reverse())
}
