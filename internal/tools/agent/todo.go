package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
)

type TodoTool struct {
	tasks []string
}

func NewTodoTool() *TodoTool {
	return &TodoTool{
		tasks: make([]string, 0),
	}
}

func (t *TodoTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "todo",
		Description: "Maintains an in-memory checklist to keep track of multi-step tasks. You can add, complete, remove, or list tasks.",
		Category:    tools.CategoryAgent,
		Permission:  tools.PermWrite,
		Parameters: []tools.Parameter{
			{
				Name:        "action",
				Type:        "string",
				Description: "Action to perform: 'add', 'complete', 'remove', or 'list'.",
				Required:    true,
			},
			{
				Name:        "task",
				Type:        "string",
				Description: "The task text. Required for 'add', 'complete', and 'remove'. For 'complete' or 'remove', it must match part of the task string.",
				Required:    false,
			},
		},
	}
}

func (t *TodoTool) Run(ctx context.Context, call tools.Call) tools.Result {
	action, ok := call.Args["action"].(string)
	if !ok || action == "" {
		return tools.Result{Error: fmt.Errorf("action is required")}
	}

	task := ""
	if v, ok := call.Args["task"].(string); ok {
		task = v
	}

	switch action {
	case "add":
		if task == "" {
			return tools.Result{Error: fmt.Errorf("task text is required for add")}
		}
		t.tasks = append(t.tasks, "[ ] "+task)
		return tools.Result{Output: "Added task. Current list:\n" + t.list()}
	case "complete":
		if task == "" {
			return tools.Result{Error: fmt.Errorf("task text is required for complete")}
		}
		found := false
		for i, v := range t.tasks {
			if strings.Contains(v, task) && strings.HasPrefix(v, "[ ]") {
				t.tasks[i] = strings.Replace(v, "[ ]", "[x]", 1)
				found = true
				break
			}
		}
		if !found {
			return tools.Result{Error: fmt.Errorf("no incomplete task matching %q found", task)}
		}
		return tools.Result{Output: "Completed task. Current list:\n" + t.list()}
	case "remove":
		if task == "" {
			return tools.Result{Error: fmt.Errorf("task text is required for remove")}
		}
		found := false
		var newTasks []string
		for _, v := range t.tasks {
			if !found && strings.Contains(v, task) {
				found = true
				continue
			}
			newTasks = append(newTasks, v)
		}
		if !found {
			return tools.Result{Error: fmt.Errorf("no task matching %q found", task)}
		}
		t.tasks = newTasks
		return tools.Result{Output: "Removed task. Current list:\n" + t.list()}
	case "list":
		return tools.Result{Output: "Current tasks:\n" + t.list()}
	default:
		return tools.Result{Error: fmt.Errorf("invalid action: %s", action)}
	}
}

func (t *TodoTool) list() string {
	if len(t.tasks) == 0 {
		return "(Empty)"
	}
	return strings.Join(t.tasks, "\n")
}
