package chat

import (
	"fmt"
	"os"
	"strings"
	"time"

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
/compact    Summarize old messages to save tokens
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
		Name:        "/index",
		Description: "Index workspace for semantic search (RAG)",
		Execute: func(m *Model) tea.Cmd {
			return func() tea.Msg {
				return indexWorkspaceMsg{}
			}
		},
	},
	{
		Name:        "/compact",
		Description: "Summarize old messages to save tokens",
		Execute: func(m *Model) tea.Cmd {
			return func() tea.Msg {
				return compactConversationMsg{}
			}
		},
	},
	{
		Name:        "/remote",
		Description: "Configure and manage remote control (Telegram, Web)",
		Execute: func(m *Model) tea.Cmd {
			return selectRemoteCmd(m)
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
		Name:        "/export",
		Description: "Export the conversation to a markdown file",
		Execute: func(m *Model) tea.Cmd {
			filename := fmt.Sprintf("windmist_export_%d.md", time.Now().Unix())
			var b strings.Builder
			b.WriteString(fmt.Sprintf("# WindMist Conversation Export\n\nDate: %s\n\n", time.Now().Format(time.RFC1123)))

			for _, msg := range m.conversation.Messages {
				role := "User"
				if msg.Role == "assistant" {
					role = "WindMist"
				}
				b.WriteString(fmt.Sprintf("## %s\n\n%s\n\n---\n\n", role, msg.Content))
			}

			err := os.WriteFile(filename, []byte(b.String()), 0644)
			if err != nil {
				m.conversation.AddAssistant(fmt.Sprintf("❌ Failed to export conversation: %v", err))
			} else {
				m.conversation.AddAssistant(fmt.Sprintf("✅ Conversation exported to `%s`", filename))
			}
			return nil
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

// mcpEnvPromptChain returns a tea.Msg that recursively prompts for each required env variable.
