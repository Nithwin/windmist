package editing

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// prepareMatcher compiles the appropriate regex or prepares case-insensitive lookup string based on SearchOptions.
func prepareMatcher(opts SearchOptions) (*regexp.Regexp, string, error) {
	lowerQuery := opts.Query
	if !opts.CaseSensitive {
		lowerQuery = strings.ToLower(opts.Query)
	}

	if opts.Type == SearchRegex {
		pattern := opts.Query
		if !opts.CaseSensitive && !strings.HasPrefix(pattern, "(?i)") {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, "", fmt.Errorf("invalid regex %q: %w", pattern, err)
		}
		return re, lowerQuery, nil
	} else if opts.WholeWord {
		pattern := `\b` + regexp.QuoteMeta(opts.Query) + `\b`
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, "", err
		}
		return re, lowerQuery, nil
	}

	return nil, lowerQuery, nil
}

// matchFile inspects a single file for occurrences of the search query up to maxAllowed matches.
func matchFile(ctx context.Context, path string, opts SearchOptions, re *regexp.Regexp, lowerQuery string, maxAllowed int) ([]Match, error) {
	if maxAllowed <= 0 {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if binary, _ := isBinaryFile(f); binary {
		return nil, nil // skip binary files silently
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	scanBuf := make([]byte, 0, 64*1024)
	scanner.Buffer(scanBuf, 10*1024*1024)

	lineNum := 1
	var matches []Match

	for scanner.Scan() {
		if ctx.Err() != nil {
			break
		}
		if len(matches) >= maxAllowed {
			break
		}

		line := scanner.Text()
		if re != nil {
			locs := re.FindAllStringIndex(line, -1)
			for _, loc := range locs {
				if len(matches) >= maxAllowed {
					break
				}
				matches = append(matches, Match{
					Line:   lineNum,
					Column: loc[0] + 1,
					Text:   line,
				})
			}
		} else {
			searchLine := line
			if !opts.CaseSensitive {
				searchLine = strings.ToLower(line)
			}
			colOffset := 0
			for {
				if len(matches) >= maxAllowed {
					break
				}
				idx := strings.Index(searchLine[colOffset:], lowerQuery)
				if idx == -1 {
					break
				}
				matches = append(matches, Match{
					Line:   lineNum,
					Column: colOffset + idx + 1,
					Text:   line,
				})
				colOffset += idx + len(lowerQuery)
				if colOffset >= len(searchLine) || len(lowerQuery) == 0 {
					break
				}
			}
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return matches, nil
	}

	return matches, nil
}

func isBinaryFile(f *os.File) (bool, error) {
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return false, err
	}
	return bytes.IndexByte(buf[:n], 0) != -1, nil
}
