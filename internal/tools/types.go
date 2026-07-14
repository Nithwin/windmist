package tools

import "context"

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
	Parameters  []Parameter
}

type Call struct {
	Name string
	Args map[string]any
}

type Result struct {
	Output any
	Error  error
}

type Tool interface {
	Definition() Definition
	Run(ctx context.Context, call Call) Result
}