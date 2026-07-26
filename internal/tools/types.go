package tools

import (
	"context"
	"time"
)

type Category string

const (
	CategoryFilesystem Category = "filesystem"
	CategoryEditing    Category = "editing"
	CategorySystem     Category = "system"
	CategorySearch     Category = "search"
	CategoryGit        Category = "git"
	CategoryWeb        Category = "web"
	CategoryAgent      Category = "agent"
)

type PermissionLevel int

const (
	PermReadOnly  PermissionLevel = iota // Auto-approved
	PermWrite                            // Needs approval first time
	PermDangerous                        // Always needs approval
)

type Parameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Enum        []string
}

type Definition struct {
	Name        string
	Description string
	Category    Category
	Permission  PermissionLevel
	Parameters  []Parameter
}

type Call struct {
	Name string
	Args map[string]any
}

type FileState struct {
	Path          string
	BeforeContent string
	AfterContent  string
	ChangeType    string // create, edit, delete
}

type Result struct {
	Output       any
	Error        error
	Duration     time.Duration // How long the tool took
	FilesRead    []string      // Files accessed
	FilesChanged []string      // Files modified
	FileStates   []FileState   // Exact before/after for Undo/Redo
	BytesChanged int64         // Total bytes changed
}

type Tool interface {
	Definition() Definition
	Run(ctx context.Context, call Call) Result
}
