package filesystem

// ExistsResult represents the structured output of the exists tool.
type ExistsResult struct {
	Exists bool   `json:"exists"`
	Type   string `json:"type,omitempty"`
}

// ListEntry represents a single file or directory entry returned by the list tool.
type ListEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Entry is a type alias for ListEntry for clean and flexible usage.
type Entry = ListEntry

// InfoResult represents detailed information about a file or directory.
type InfoResult struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Mode string `json:"mode"`
}
