package chat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nithwin/WindMist/internal/mcp"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectMCPCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		var options []selector.Option
		for i, name := range mcp.GetCatalogList() {
			options = append(options, selector.Option{
				Label: name,
				Value: fmt.Sprintf("%d", i),
			})
		}

		return showInlineSelectorMsg{
			Title:   "Select MCP Server",
			Options: options,
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(opt selector.Option) tea.Cmd {
				return func() tea.Msg {
					idx, _ := strconv.Atoi(opt.Value)
					entry, ok := mcp.GetCatalogEntry(idx)
					if !ok {
						return switchErrorMsg{Err: fmt.Errorf("invalid MCP selection")}
					}
					return mcpEnvPromptChain(m, entry, make(map[string]string), 0)()
				}
			},
		}
	}
}

func mcpEnvPromptChain(m *Model, entry *mcp.CatalogEntry, envValues map[string]string, index int) tea.Cmd {
	return func() tea.Msg {
		if index >= len(entry.RequiredEnv) {
			if err := mcp.Install(entry, envValues); err != nil {
				return switchErrorMsg{Err: fmt.Errorf("failed to save config: %w", err)}
			}
			return mcpInstallSuccessMsg{Name: entry.Name}
		}

		envKey := entry.RequiredEnv[index]

		if entry.ID == "github" && envKey == "GITHUB_PERSONAL_ACCESS_TOKEN" {
			// Notify user in chat
			m.conversation.AddAssistant("⏳ Initializing GitHub OAuth Flow...")
			m.refreshViewport()

			token, err := mcp.PerformGithubOAuth(func(uri, code string) {
				msg := fmt.Sprintf("🔒 **GitHub Authentication Required**\n\n1. Open this link: %s\n2. Enter this code: `%s`\n\n_Waiting for authorization..._", uri, code)
				m.conversation.AddAssistant(msg)
				m.refreshViewport()
			})

			if err == nil && token != "" {
				m.conversation.AddAssistant("✅ Successfully authenticated with GitHub!")
				m.refreshViewport()
				envValues[envKey] = token
				return mcpEnvPromptChain(m, entry, envValues, index+1)()
			}

			m.conversation.AddAssistant(fmt.Sprintf("⚠️ OAuth failed (%v), falling back to manual entry...", err))
			m.refreshViewport()
		}

		prompt := fmt.Sprintf("Enter %s:", envKey)
		if entry.EnvPrompt != nil && entry.EnvPrompt[envKey] != "" {
			prompt = entry.EnvPrompt[envKey]
		}

		return showInlinePromptMsg{
			Prompt:     prompt,
			IsPassword: true,
			OnSubmit: func(val string) tea.Cmd {
				return func() tea.Msg {
					val = strings.TrimSpace(val)
					if val == "" {
						return switchCancelMsg{}
					}
					envValues[envKey] = val
					return mcpEnvPromptChain(m, entry, envValues, index+1)()
				}
			},
		}
	}
}
