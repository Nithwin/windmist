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
