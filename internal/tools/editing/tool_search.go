package editing

import (
	"context"
	"os"

	"github.com/Nithwin/WindMist/internal/tools"
)

// SearchTool implements the AI tool interface for search_text.
type SearchTool struct{}

func NewSearchTool() *SearchTool {
	return &SearchTool{}
}

func (t *SearchTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "search_text",
		Description: "Searches for text or regex patterns across files in a directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "query",
				Type:        "string",
				Description: "Text or regex query to search for.",
				Required:    true,
			},
			{
				Name:        "path",
				Type:        "string",
				Description: "Root directory or file to search within (defaults to current directory).",
				Required:    false,
			},
			{
				Name:        "case_sensitive",
				Type:        "boolean",
				Description: "Whether the search is case-sensitive (defaults to false).",
				Required:    false,
			},
			{
				Name:        "whole_word",
				Type:        "boolean",
				Description: "Whether to match whole words only (defaults to false).",
				Required:    false,
			},
			{
				Name:        "regex",
				Type:        "boolean",
				Description: "Whether the query is a regular expression (defaults to false).",
				Required:    false,
			},
			{
				Name:        "max_matches",
				Type:        "integer",
				Description: "Maximum number of total matches to return across all files (defaults to 1000).",
				Required:    false,
			},
			{
				Name:        "max_files",
				Type:        "integer",
				Description: "Maximum number of distinct files to return (defaults to 500).",
				Required:    false,
			},
		},
	}
}

func (t *SearchTool) Run(ctx context.Context, call tools.Call) tools.Result {
	query, ok := call.Args["query"].(string)
	if !ok || query == "" {
		return tools.Result{Error: os.ErrInvalid}
	}

	root, _ := call.Args["path"].(string)
	if root == "" {
		root = "."
	}

	caseSensitive := false
	if cs, ok := call.Args["case_sensitive"].(bool); ok {
		caseSensitive = cs
	} else if csStr, ok := call.Args["case_sensitive"].(string); ok {
		caseSensitive = (csStr == "true" || csStr == "1")
	}

	wholeWord := false
	if ww, ok := call.Args["whole_word"].(bool); ok {
		wholeWord = ww
	} else if wwStr, ok := call.Args["whole_word"].(string); ok {
		wholeWord = (wwStr == "true" || wwStr == "1")
	}

	isRegex := false
	if r, ok := call.Args["regex"].(bool); ok {
		isRegex = r
	} else if rStr, ok := call.Args["regex"].(string); ok {
		isRegex = (rStr == "true" || rStr == "1")
	}

	maxMatches := 1000
	if mm, ok := call.Args["max_matches"].(int); ok && mm > 0 {
		maxMatches = mm
	} else if mmFloat, ok := call.Args["max_matches"].(float64); ok && mmFloat > 0 {
		maxMatches = int(mmFloat)
	}

	maxFiles := 500
	if mf, ok := call.Args["max_files"].(int); ok && mf > 0 {
		maxFiles = mf
	} else if mfFloat, ok := call.Args["max_files"].(float64); ok && mfFloat > 0 {
		maxFiles = int(mfFloat)
	}

	searchType := SearchText
	if isRegex {
		searchType = SearchRegex
	}

	opts := SearchOptions{
		Root:          root,
		Query:         query,
		Type:          searchType,
		CaseSensitive: caseSensitive,
		WholeWord:     wholeWord,
		MaxMatches:    maxMatches,
		MaxFiles:      maxFiles,
	}

	results, err := Search(ctx, opts)
	if err != nil {
		return tools.Result{Error: err}
	}

	return tools.Result{Output: results}
}
