package editing

import (
	"context"
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

// Apply executes every FilePatch sequentially using the default Executor.
// To run in custom or dry-run execution pipelines, use Executor directly.
func (p *Patch) Apply(ctx context.Context) error {
	return NewExecutor().Apply(ctx, p)
}

// Reverse returns a new Patch where every FilePatch and every Operation is precisely inverted,
// ordered in reverse so that applying the reversed patch cleanly restores the original state.
// Operation structs inside the returned Patch are entirely new, preserving immutability of the original operations.
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

// Undo applies the inverse of the provided patch using the default Executor.
func Undo(ctx context.Context, p *Patch) error {
	return NewExecutor().Undo(ctx, p)
}
