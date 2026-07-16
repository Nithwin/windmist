package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui/selector"
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
			m.refreshViewport()
			return nil
		},
	},
	{
		Name:        "/clear",
		Description: "Clear conversation",
		Execute: func(m *Model) tea.Cmd {
			m.conversation.Clear()
			m.refreshViewport()
			return nil
		},
	},
	{
		Name:        "/model",
		Description: "Change model",
		Execute: func(m *Model) tea.Cmd {
			return selectModelCmd(m)
		},
	},
	{
		Name:        "/provider",
		Description: "Change provider",
		Execute: func(m *Model) tea.Cmd {
			return selectProviderCmd(m)
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

func selectProviderCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		// 1. Release terminal of main program so selector can render cleanly
		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		// 2. Select Provider
		providerOpt, err := selector.Run(
			"Select AI Provider",
			"Choose which AI provider you want WindMist to use:",
			config.GetProviderOptions(),
		)
		if err != nil {
			return switchCancelMsg{}
		}

		// 3. Select Model for this provider
		ollamaBaseURL := ""
		if pConfig, ok := m.cfg.Providers[providerOpt.Value]; ok {
			ollamaBaseURL = pConfig.BaseURL
		}
		modelOpt, err := selector.Run(
			fmt.Sprintf("Select Model for %s", providerOpt.Value),
			"Choose the active model for this provider:",
			config.GetModelOptions(providerOpt.Value, ollamaBaseURL),
		)
		if err != nil {
			return switchCancelMsg{}
		}

		modelValue := modelOpt.Value
		if modelValue == "__CUSTOM__" {
			customVal, err := selector.RunInput("Custom Model ID", "Enter exact model ID (e.g. gpt-4o)", "")
			if err != nil {
				return switchCancelMsg{}
			}
			modelValue = customVal
		}

		return switchProviderSuccessMsg{
			Provider: providerOpt.Value,
			Model:    modelValue,
		}
	}
}

func selectModelCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		// 1. Release terminal of main program
		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		// 2. Select Model for current provider
		ollamaBaseURL := ""
		if pConfig, ok := m.cfg.Providers[m.cfg.AI.Provider]; ok {
			ollamaBaseURL = pConfig.BaseURL
		}
		modelOpt, err := selector.Run(
			fmt.Sprintf("Select Model for %s", m.cfg.AI.Provider),
			"Choose the active model to use:",
			config.GetModelOptions(m.cfg.AI.Provider, ollamaBaseURL),
		)
		if err != nil {
			return switchCancelMsg{}
		}

		modelValue := modelOpt.Value
		if modelValue == "__CUSTOM__" {
			customVal, err := selector.RunInput("Custom Model ID", "Enter exact model ID (e.g. gpt-4o)", "")
			if err != nil {
				return switchCancelMsg{}
			}
			modelValue = customVal
		}

		return switchModelSuccessMsg{
			Model: modelValue,
		}
	}
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
