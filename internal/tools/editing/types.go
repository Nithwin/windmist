package editing

type SearchType int

const (
	SearchText SearchType = iota
	SearchFileName
	SearchRegex
	SearchSymbol
)

// Match represents a single occurrence of a search query within a file.
type Match struct {
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Text   string `json:"text"`
}

// SearchResult groups all matches found inside a specific file.
type SearchResult struct {
	Path    string  `json:"path"`
	Matches []Match `json:"matches"`
}

// Results represents a collection of search results across multiple files.
type Results []SearchResult

// SearchOptions configures the behavior of the search engine.
type SearchOptions struct {
	Root          string     `json:"root"`
	Query         string     `json:"query"`
	Type          SearchType `json:"type,omitempty"`
	CaseSensitive bool       `json:"case_sensitive"`
	WholeWord     bool       `json:"whole_word"`
	IncludeHidden bool       `json:"include_hidden"`
	MaxMatches    int        `json:"max_matches"`
	MaxFiles      int        `json:"max_files"`
}

type OperationType int

const (
	Insert OperationType = iota
	Replace
	Delete
)

// Operation represents a common model for every file modification produced by editing functions.
type Operation struct {
	Type        OperationType `json:"type"`
	File        string        `json:"file"`
	StartLine   int           `json:"start_line"`
	EndLine     int           `json:"end_line"`
	StartColumn int           `json:"start_column,omitempty"`
	EndColumn   int           `json:"end_column,omitempty"`
	OldText     string        `json:"old_text,omitempty"`
	NewText     string        `json:"new_text,omitempty"`
}

type ReplaceOptions struct {
	File            string `json:"file"`
	OldText         string `json:"old_text"`
	NewText         string `json:"new_text"`
	ReplaceAll      bool   `json:"replace_all"`
	MaxReplacements int    `json:"max_replacements"`
}

type ReplaceResult struct {
	Operations []Operation `json:"operations"`
	Replacements int       `json:"replacements"`
}

type ReplaceRangeOptions struct {
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	NewText   string `json:"new_text"`
}

type ReplaceRangeResult struct {
	Operation Operation `json:"operation"`
}

type InsertOptions struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	NewText string `json:"new_text"`
}

type InsertResult struct {
	Operation Operation `json:"operation"`
}




