package editing

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var errSearchLimitReached = errors.New("search limit reached")

// walkDirectory traverses the directory tree or checks a single file, streaming results as soon as matches are found.
func walkDirectory(ctx context.Context, opts SearchOptions, re *regexp.Regexp, lowerQuery string, onResult func(SearchResult) error) error {
	info, err := os.Stat(opts.Root)
	if err != nil {
		return err
	}

	totalMatches := 0
	totalFiles := 0

	// Handle single file root
	if !info.IsDir() {
		if opts.Type == SearchFileName {
			baseName := filepath.Base(opts.Root)
			matched := false
			if re != nil {
				matched = re.MatchString(baseName)
			} else if opts.CaseSensitive {
				matched = strings.Contains(baseName, opts.Query)
			} else {
				matched = strings.Contains(strings.ToLower(baseName), lowerQuery)
			}
			if matched {
				totalMatches++
				totalFiles++
				return onResult(SearchResult{
					Path: opts.Root,
					Matches: []Match{
						{Line: 1, Column: 1, Text: baseName},
					},
				})
			}
			return nil
		}

		matches, fileErr := matchFile(ctx, opts.Root, opts, re, lowerQuery, opts.MaxMatches)
		if fileErr != nil {
			return fileErr
		}
		if len(matches) > 0 {
			totalMatches += len(matches)
			totalFiles++
			return onResult(SearchResult{
				Path:    opts.Root,
				Matches: matches,
			})
		}
		return nil
	}

	// Walk directory tree
	err = filepath.WalkDir(opts.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if totalMatches >= opts.MaxMatches || totalFiles >= opts.MaxFiles {
			return errSearchLimitReached
		}

		if d.IsDir() {
			if !opts.IncludeHidden && strings.HasPrefix(d.Name(), ".") && path != opts.Root {
				return filepath.SkipDir
			}
			if (d.Name() == "node_modules" || d.Name() == "vendor") && !opts.IncludeHidden && path != opts.Root {
				return filepath.SkipDir
			}
			return nil
		}

		if !opts.IncludeHidden && strings.HasPrefix(d.Name(), ".") && path != opts.Root {
			return nil
		}

		cleanPath := path
		if opts.Root == "." || opts.Root == "./" {
			cleanPath = filepath.Clean(path)
		}

		if opts.Type == SearchFileName {
			matched := false
			if re != nil {
				matched = re.MatchString(d.Name())
			} else if opts.CaseSensitive {
				matched = strings.Contains(d.Name(), opts.Query)
			} else {
				matched = strings.Contains(strings.ToLower(d.Name()), lowerQuery)
			}
			if matched {
				totalMatches++
				totalFiles++
				resErr := onResult(SearchResult{
					Path: cleanPath,
					Matches: []Match{
						{Line: 1, Column: 1, Text: d.Name()},
					},
				})
				if resErr != nil {
					return resErr
				}
				if totalMatches >= opts.MaxMatches || totalFiles >= opts.MaxFiles {
					return errSearchLimitReached
				}
			}
			return nil
		}

		matches, fileErr := matchFile(ctx, path, opts, re, lowerQuery, opts.MaxMatches-totalMatches)
		if fileErr != nil {
			return nil
		}
		if len(matches) > 0 {
			totalMatches += len(matches)
			totalFiles++
			resErr := onResult(SearchResult{
				Path:    cleanPath,
				Matches: matches,
			})
			if resErr != nil {
				return resErr
			}
			if totalMatches >= opts.MaxMatches || totalFiles >= opts.MaxFiles {
				return errSearchLimitReached
			}
		}
		return nil
	})

	if errors.Is(err, errSearchLimitReached) {
		return nil
	}
	return err
}
