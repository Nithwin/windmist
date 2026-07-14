package editing

import (
	"context"
	"fmt"
	"os"
)

// SearchStream executes the search engine and streams each file's SearchResult immediately to onResult.
func SearchStream(ctx context.Context, opts SearchOptions, onResult func(SearchResult) error) error {
	if opts.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	if opts.Root == "" {
		opts.Root = "."
	}
	if opts.MaxMatches <= 0 {
		opts.MaxMatches = 1000
	}
	if opts.MaxFiles <= 0 {
		opts.MaxFiles = 500
	}

	if opts.Type == SearchRegex && opts.WholeWord {
		return fmt.Errorf("cannot combine whole_word with regex search type")
	}

	_, err := os.Stat(opts.Root)
	if err != nil {
		return err
	}

	re, lowerQuery, err := prepareMatcher(opts)
	if err != nil {
		return err
	}

	return walkDirectory(ctx, opts, re, lowerQuery, onResult)
}

// Search executes the search engine and accumulates all results into a slice.
func Search(ctx context.Context, opts SearchOptions) (Results, error) {
	results := make(Results, 0)
	err := SearchStream(ctx, opts, func(res SearchResult) error {
		results = append(results, res)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}
