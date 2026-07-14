package chat

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Command struct {
	Name        string
	Description string
	Execute     func(*Model) tea.Cmd
}

// Registry keeps commands in display order.
var Registry = []Command{
	{
		Name:        "/help",
		Description: "Show available commands",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.AddAssistant(
				`Available commands:

/help       Show available commands
/new        Start a new conversation
/model      Change model
/provider   Change provider
/clear      Clear conversation
/exit       Exit WindMist`,
			)
			return nil
		},
	},
	{
		Name:        "/new",
		Description: "Start a new conversation",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.Clear()
			return nil
		},
	},
	{
		Name:        "/clear",
		Description: "Clear conversation",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.Clear()
			return nil
		},
	},
	{
		Name:        "/model",
		Description: "Change model",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.AddAssistant("Model selection coming soon.")
			return nil
		},
	},
	{
		Name:        "/provider",
		Description: "Change provider",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.AddAssistant("Provider selection coming soon.")
			return nil
		},
	},
	{
		Name:        "/exit",
		Description: "Exit WindMist",
		Execute: func(m *Model) tea.Cmd {
			return tea.Quit
		},
	},
	{
		Name:        "/quit",
		Description: "Exit WindMist",
		Execute: func(m *Model) tea.Cmd {
			return tea.Quit
		},
	},
}

func FilterCommands(input string) []Command {
	if input == "/" {
		return Registry
	}

	var filtered []Command

	for _, cmd := range Registry {
		if strings.HasPrefix(cmd.Name, input) {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}

func FindCommand(name string) (Command, bool) {
	for _, cmd := range Registry {
		if cmd.Name == name {
			return cmd, true
		}
	}

	return Command{}, false
}
