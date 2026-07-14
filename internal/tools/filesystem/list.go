package filesystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nithwin/WindMist/internal/tools"
)

type ListTool struct{}

func NewListTool() *ListTool {
	return &ListTool{}
}

func (t *ListTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "list",
		Description: "Lists files and directories inside a specified directory.",
		Parameters: []tools.Parameter{
			{
				Name:        "path",
				Type:        "string",
				Description: "Path of the directory to list.",
				Required:    true,
			},
			{
				Name:        "recursive",
				Type:        "boolean",
				Description: "If true, recursively lists all files and directories in nested subdirectories.",
				Required:    false,
			},
		},
	}
}

func (t *ListTool) Run(ctx context.Context, call tools.Call) tools.Result {
	path, ok := call.Args["path"].(string)
	if !ok || path == "" {
		return tools.Result{
			Error: os.ErrInvalid,
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return tools.Result{
			Error: err,
		}
	}
	if !info.IsDir() {
		return tools.Result{
			Error: fmt.Errorf("path %q is not a directory", path),
		}
	}

	recursive := false
	if r, ok := call.Args["recursive"].(bool); ok {
		recursive = r
	} else if rStr, ok := call.Args["recursive"].(string); ok {
		recursive = (rStr == "true" || rStr == "1")
	}

	if !recursive {
		entries, err := os.ReadDir(path)
		if err != nil {
			return tools.Result{
				Error: err,
			}
		}

		result := make([]Entry, 0, len(entries))
		for _, e := range entries {
			entryType := "file"
			if e.IsDir() {
				entryType = "directory"
			}
			result = append(result, Entry{
				Name: e.Name(),
				Type: entryType,
			})
		}

		return tools.Result{
			Output: result,
		}
	}

	result := make([]Entry, 0)
	err = filepath.WalkDir(path, func(walkPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if walkPath == path {
			return nil
		}

		entryType := "file"
		if d.IsDir() {
			entryType = "directory"
		}

		entryName := walkPath
		if path == "." || path == "./" {
			entryName = filepath.Clean(walkPath)
		}

		result = append(result, Entry{
			Name: entryName,
			Type: entryType,
		})
		return nil
	})

	if err != nil {
		return tools.Result{
			Error: err,
		}
	}

	return tools.Result{
		Output: result,
	}
}
