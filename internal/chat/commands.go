package chat

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/mcp"
	"github.com/Nithwin/WindMist/internal/ui"
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
/sessions   Load a previous session
/undo       Undo the last AI file edit
/redo       Redo the last undone file edit
/model      Change model
/mode       Change agent mode
/provider   Change provider
/subagent   Configure sub-agent (cheaper background model)
/theme      Change UI theme
/exit       Exit WindMist`,
			)
			return nil
		},
	},
	{
		Name:        "/new",
		Description: "Start a new conversation",
		Execute: func(m *Model) tea.Cmd {
			return func() tea.Msg {
				return createNewSessionMsg{}
			}
		},
	},
	{
		Name:        "/sessions",
		Description: "Load a previous session",
		Execute: func(m *Model) tea.Cmd {
			return selectSessionCmd(m)
		},
	},
	{
		Name:        "/undo",
		Description: "Undo the last AI file edit",
		Execute: func(m *Model) tea.Cmd {
			return func() tea.Msg {
				return undoFileChangeMsg{}
			}
		},
	},
	{
		Name:        "/redo",
		Description: "Redo the last undone file edit",
		Execute: func(m *Model) tea.Cmd {
			return func() tea.Msg {
				return redoFileChangeMsg{}
			}
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
		Name:        "/mode",
		Description: "Change agent mode",
		Execute: func(m *Model) tea.Cmd {
			return selectModeCmd(m)
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
		Name:        "/subagent",
		Description: "Configure sub-agent (cheaper background model)",
		Execute: func(m *Model) tea.Cmd {
			return selectSubagentCmd(m)
		},
	},
	{
		Name:        "/theme",
		Description: "Change UI theme",
		Execute: func(m *Model) tea.Cmd {
			return selectThemeCmd(m)
		},
	},
	{
		Name:        "/apikey",
		Description: "Set API Key for the current provider",
		Execute: func(m *Model) tea.Cmd {
			return setAPIKeyCmd(m)
		},
	},
	{
		Name:        "/mcp",
		Description: "Install an MCP server (e.g. GitHub, Postgres)",
		Execute: func(m *Model) tea.Cmd {
			return selectMCPCmd(m)
		},
	},
	{
		Name:        "/exit",
		Description: "Exit WindMist",
		Execute: func(m *Model) tea.Cmd {
			if m.agent != nil {
				m.agent.Close()
			}
			return tea.Quit
		},
	},
	{
		Name:        "/quit",
		Description: "Exit WindMist",
		Execute: func(m *Model) tea.Cmd {
			if m.agent != nil {
				m.agent.Close()
			}
			return tea.Quit
		},
	},
}

func selectSessionCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}
		if m.store == nil {
			return switchErrorMsg{Err: fmt.Errorf("database not initialized")}
		}

		cwd, _ := os.Getwd()
		sessions, err := m.store.ListSessionsByProject(cwd)
		if err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to fetch sessions: %w", err)}
		}

		if len(sessions) == 0 {
			return switchErrorMsg{Err: fmt.Errorf("no past sessions found in this project")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		var options []selector.Option
		for _, s := range sessions {
			desc := fmt.Sprintf("%s | Tokens: %d | Cost: $%.3f", s.UpdatedAt.Format("Jan 02 15:04"), s.TokenCount, s.CostEstimate)
			options = append(options, selector.Option{
				Label: s.Title,
				Desc:  desc,
				Value: s.ID,
			})
		}

		opt, err := selector.Run("Select Session", "Choose a previous session to resume:", options)
		if err != nil {
			return switchCancelMsg{}
		}

		return switchSessionSuccessMsg{
			SessionID: opt.Value,
		}
	}
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
			m.cfg.GetModelOptions(providerOpt.Value, ollamaBaseURL),
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

			// Save the custom model so it shows up next time
			m.cfg.AddCustomModel(providerOpt.Value, modelValue)
			_ = config.Save(m.cfg)
		}

		return switchProviderSuccessMsg{
			Provider: providerOpt.Value,
			Model:    modelValue,
		}
	}
}

func selectSubagentCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		providerOpt, err := selector.Run(
			"Select Sub-Agent Provider",
			"Choose which AI provider the Sub-Agent should use (Auto uses main config):",
			append([]selector.Option{{Label: "Auto (Use Main Config)", Value: "auto"}}, config.GetProviderOptions()...),
		)
		if err != nil {
			return switchCancelMsg{}
		}

		if providerOpt.Value == "auto" {
			return switchSubagentSuccessMsg{
				Provider: "",
				Model:    "",
			}
		}

		ollamaBaseURL := ""
		if pConfig, ok := m.cfg.Providers[providerOpt.Value]; ok {
			ollamaBaseURL = pConfig.BaseURL
		}
		modelOpt, err := selector.Run(
			fmt.Sprintf("Select Sub-Agent Model for %s", providerOpt.Value),
			"Choose the active model for this provider (Auto uses fast default):",
			append([]selector.Option{{Label: "Auto (Fast Default)", Value: "auto"}}, m.cfg.GetModelOptions(providerOpt.Value, ollamaBaseURL)...),
		)
		if err != nil {
			return switchCancelMsg{}
		}

		modelValue := modelOpt.Value
		if modelValue == "auto" {
			modelValue = ""
		} else if modelValue == "__CUSTOM__" {
			customVal, err := selector.RunInput("Custom Model ID", "Enter exact model ID (e.g. gpt-4o-mini)", "")
			if err != nil {
				return switchCancelMsg{}
			}
			modelValue = customVal
			m.cfg.AddCustomModel(providerOpt.Value, modelValue)
			_ = config.Save(m.cfg)
		}

		return switchSubagentSuccessMsg{
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
			m.cfg.GetModelOptions(m.cfg.AI.Provider, ollamaBaseURL),
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

			// Save the custom model so it shows up next time
			m.cfg.AddCustomModel(m.cfg.AI.Provider, modelValue)
			_ = config.Save(m.cfg)
		}

		return switchModelSuccessMsg{
			Model: modelValue,
		}
	}
}

func selectModeCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		options := []selector.Option{
			{Label: "Auto", Desc: "Dynamically switches between Build and Plan based on prompt", Value: "auto"},
			{Label: "Build", Desc: "Full autonomy mode with read/write access", Value: "build"},
			{Label: "Plan", Desc: "Read-only mode for architecture and analysis", Value: "plan"},
		}

		opt, err := selector.Run("Select Agent Mode", "Choose how the AI should behave:", options)
		if err != nil {
			return switchCancelMsg{}
		}

		return switchModeSuccessMsg{Mode: opt.Value}
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

func selectThemeCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		themes := ui.AvailableThemes()
		var options []selector.Option
		for _, t := range themes {
			options = append(options, selector.Option{
				Label: t,
				Value: t,
			})
		}

		opt, err := selector.RunWithDefault("Select Theme", "Choose a UI theme:", options, ui.CurrentThemeName)
		if err != nil {
			return switchCancelMsg{}
		}

		return switchThemeSuccessMsg{
			Theme: opt.Value,
		}
	}
}

func selectMCPCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		var options []selector.Option
		for i, name := range mcp.GetCatalogList() {
			options = append(options, selector.Option{
				Label: name,
				Value: fmt.Sprintf("%d", i), // Use index as value
			})
		}

		opt, err := selector.Run("Select MCP Server", "Choose an MCP Server to install:", options)
		if err != nil {
			return switchCancelMsg{}
		}

		// Show prompts for required env vars
		idx, _ := strconv.Atoi(opt.Value)
		entry, ok := mcp.GetCatalogEntry(idx)
		if !ok {
			return switchErrorMsg{Err: fmt.Errorf("invalid MCP selection")}
		}

		envValues := make(map[string]string)
		for _, envKey := range entry.RequiredEnv {
			prompt := fmt.Sprintf("Enter %s:", envKey)
			if entry.EnvPrompt != nil && entry.EnvPrompt[envKey] != "" {
				prompt = entry.EnvPrompt[envKey]
			}
			
			// Simple fallback prompt via terminal since we released the TUI
			fmt.Printf("\n%s\n> ", prompt)
			var val string
			fmt.Scanln(&val)
			envValues[envKey] = strings.TrimSpace(val)
		}

		if err := mcp.Install(entry, envValues); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to save config: %w", err)}
		}

		return mcpInstallSuccessMsg{
			Name: entry.Name,
		}
	}
}

func setAPIKeyCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}

		if err := program.ReleaseTerminal(); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to release terminal: %w", err)}
		}
		defer program.RestoreTerminal()

		provider := m.cfg.AI.Provider
		if provider == "" {
			provider = "default"
		}
		
		fmt.Printf("\n🔑 Enter API Key for [%s]:\n> ", provider)
		var val string
		fmt.Scanln(&val)
		val = strings.TrimSpace(val)
		
		if val == "" {
			return switchCancelMsg{}
		}

		if err := m.cfg.SetAPIKey(provider, val); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to set api key: %w", err)}
		}
		
		if err := config.Save(m.cfg); err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to save config: %w", err)}
		}

		return setAPIKeySuccessMsg{
			Provider: provider,
		}
	}
}
