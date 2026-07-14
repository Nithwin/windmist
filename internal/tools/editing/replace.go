package editing

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ReplaceText finds exact matches of old_text using the search engine and replaces them with new_text.
// It enforces strict ambiguity checks when old_text occurs multiple times and replace_all is not set.
func ReplaceText(ctx context.Context, opts ReplaceOptions) (*ReplaceResult, error) {
	if opts.File == "" || opts.OldText == "" {
		return nil, fmt.Errorf("file and old_text cannot be empty")
	}

	f, err := os.Open(opts.File)
	if err != nil {
		return nil, err
	}
	binary, binErr := isBinaryFile(f)
	_ = f.Close()
	if binErr != nil {
		return nil, binErr
	}
	if binary {
		return nil, fmt.Errorf("cannot perform text replacement on binary file %q", opts.File)
	}

	// Step 1: Search exact matches using the Search engine
	searchOpts := SearchOptions{
		Root:          opts.File,
		Query:         opts.OldText,
		Type:          SearchText,
		CaseSensitive: true,
	}
	results, err := Search(ctx, searchOpts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 || len(results[0].Matches) == 0 {
		return nil, fmt.Errorf("old_text not found in file %q", opts.File)
	}

	matches := results[0].Matches
	totalMatches := len(matches)

	// Step 2: Enforce ambiguity checks
	if !opts.ReplaceAll && opts.MaxReplacements <= 1 && totalMatches > 1 {
		return nil, fmt.Errorf("ambiguous replacement: old_text appears %d times in file %q. Set replace_all=true, specify max_replacements, or use replace_range with exact line coordinates", totalMatches, opts.File)
	}

	limit := totalMatches
	if !opts.ReplaceAll {
		if opts.MaxReplacements > 0 && opts.MaxReplacements < limit {
			limit = opts.MaxReplacements
		} else if opts.MaxReplacements <= 0 {
			limit = 1
		}
	}

	// Step 3: Generate Operations using exact match line coordinates from Search
	ops := make([]Operation, 0, limit)
	newlines := strings.Count(opts.OldText, "\n")
	for i := 0; i < limit; i++ {
		m := matches[i]
		ops = append(ops, Operation{
			Type:        Replace,
			File:        opts.File,
			StartLine:   m.Line,
			EndLine:     m.Line + newlines,
			StartColumn: m.Column,
			EndColumn:   m.Column + len(opts.OldText),
			OldText:     opts.OldText,
			NewText:     opts.NewText,
		})
	}

	// Step 4: Execute replacement and atomically write to disk
	content, err := os.ReadFile(opts.File)
	if err != nil {
		return nil, err
	}

	updatedStr := strings.Replace(string(content), opts.OldText, opts.NewText, limit)
	if err := WriteFile(opts.File, []byte(updatedStr), 0); err != nil {
		return nil, err
	}

	return &ReplaceResult{
		Operations:   ops,
		Replacements: len(ops),
	}, nil
}
