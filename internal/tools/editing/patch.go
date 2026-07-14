package editing

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
)

// NewPatch creates a new Patch and groups the given operations by target file.
func NewPatch(ops ...Operation) *Patch {
	p := &Patch{}
	p.Add(ops...)
	return p
}

// Add appends one or more operations to the patch, automatically grouping them into their corresponding FilePatch.
func (p *Patch) Add(ops ...Operation) {
	if p == nil {
		return
	}
	for _, op := range ops {
		if op.File == "" {
			continue
		}
		found := false
		for i := range p.Files {
			if p.Files[i].Path == op.File {
				p.Files[i].Operations = append(p.Files[i].Operations, op)
				found = true
				break
			}
		}
		if !found {
			p.Files = append(p.Files, FilePatch{
				Path:       op.File,
				Operations: []Operation{op},
			})
		}
	}
}

// AllOperations returns a flattened slice of all operations across all grouped files in the patch.
func (p *Patch) AllOperations() []Operation {
	if p == nil {
		return nil
	}
	var total int
	for _, fp := range p.Files {
		total += len(fp.Operations)
	}
	ops := make([]Operation, 0, total)
	for _, fp := range p.Files {
		ops = append(ops, fp.Operations...)
	}
	return ops
}


// Apply executes every FilePatch sequentially. For each file, content is read exactly once,
// all operations are sorted in descending order of StartLine and applied in memory, and the result is written to disk once.
func (p *Patch) Apply(ctx context.Context) error {
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
				return fmt.Errorf("patch apply failed on %s: %w", fp.Path, err)
			}
			continue
		}

		f, err := os.Open(fp.Path)
		if err != nil {
			return fmt.Errorf("patch apply failed to open %s: %w", fp.Path, err)
		}
		binary, binErr := isBinaryFile(f)
		_ = f.Close()
		if binErr != nil {
			return fmt.Errorf("patch apply failed checking binary state on %s: %w", fp.Path, binErr)
		}
		if binary {
			return fmt.Errorf("patch apply failed: cannot perform line editing on binary file %q", fp.Path)
		}

		content, err := os.ReadFile(fp.Path)
		if err != nil {
			return fmt.Errorf("patch apply failed to read %s: %w", fp.Path, err)
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
					return fmt.Errorf("patch apply failed on %s: invalid replace coordinates [%d, %d] (file has %d lines)", fp.Path, op.StartLine, op.EndLine, totalLines)
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
					return fmt.Errorf("patch apply failed on %s: invalid insert line %d (file has %d lines)", fp.Path, op.StartLine, totalLines)
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
					return fmt.Errorf("patch apply failed on %s: invalid delete coordinates [%d, %d] (file has %d lines)", fp.Path, op.StartLine, op.EndLine, totalLines)
				}
				oldSlice := lines[op.StartLine-1 : op.EndLine]
				updated := make([]string, 0, len(lines)-len(oldSlice))
				updated = append(updated, lines[:op.StartLine-1]...)
				updated = append(updated, lines[op.EndLine:]...)
				lines = updated

			default:
				return fmt.Errorf("patch apply failed on %s: unsupported operation type %v", fp.Path, op.Type)
			}
		}

		updatedContent := strings.Join(lines, "\n")
		if hasTrailingNewline && (len(updatedContent) == 0 || !strings.HasSuffix(updatedContent, "\n")) {
			updatedContent += "\n"
		}

		if err := WriteFile(fp.Path, []byte(updatedContent), 0); err != nil {
			return fmt.Errorf("patch apply failed to write %s: %w", fp.Path, err)
		}
	}

	return nil
}

// Reverse returns a new Patch where every FilePatch and every Operation is precisely inverted,
// ordered in reverse so that applying the reversed patch cleanly restores the original state.
func (p *Patch) Reverse() *Patch {
	if p == nil || len(p.Files) == 0 {
		return NewPatch()
	}

	revPatch := &Patch{
		Files: make([]FilePatch, len(p.Files)),
	}

	for i, fp := range p.Files {
		n := len(fp.Operations)
		revOps := make([]Operation, n)
		for j, op := range fp.Operations {
			revIdx := n - 1 - j
			switch op.Type {
			case Replace:
				revOps[revIdx] = Operation{
					Type:      Replace,
					File:      op.File,
					StartLine: op.StartLine,
					EndLine:   op.EndLine,
					OldText:   op.NewText,
					NewText:   op.OldText,
				}
			case Insert:
				revOps[revIdx] = Operation{
					Type:      Delete,
					File:      op.File,
					StartLine: op.StartLine,
					EndLine:   op.EndLine,
					OldText:   op.NewText,
					NewText:   "",
				}
			case Delete:
				revOps[revIdx] = Operation{
					Type:      Insert,
					File:      op.File,
					StartLine: op.StartLine,
					EndLine:   op.StartLine,
					OldText:   "",
					NewText:   op.OldText,
				}
			default:
				revOps[revIdx] = op
			}
		}
		revPatch.Files[len(p.Files)-1-i] = FilePatch{
			Path:       fp.Path,
			Operations: revOps,
		}
	}

	return revPatch
}

// Undo applies the inverse of the provided patch, cleanly rolling back all modifications across all grouped files.
func Undo(ctx context.Context, p *Patch) error {
	if p == nil {
		return nil
	}
	return p.Reverse().Apply(ctx)
}
