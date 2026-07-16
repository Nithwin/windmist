package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all user interactions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(m.width - 12)
		m.refreshViewport()
		return m, nil

	case tea.KeyMsg:
		// Scroll conversation when command palette is closed.
		if !m.showCommands {
			switch msg.String() {

			case "up":
				m.viewport.ScrollUp(1)
				return m, nil

			case "down":
				m.viewport.ScrollDown(1)
				return m, nil

			case "pgup":
				m.viewport.ScrollUp(m.viewport.Height / 2)
				return m, nil

			case "pgdown":
				m.viewport.ScrollDown(m.viewport.Height / 2)
				return m, nil

			case "home":
				m.viewport.GotoTop()
				return m, nil

			case "end":
				m.viewport.GotoBottom()
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		// Hide splash on first key press.
		if m.showSplash {
			m.showSplash = false

			// Preserve the first typed character.
			if len(msg.String()) == 1 {
				m.input.SetValue(msg.String())
				m.input.CursorEnd()
			}

			m.refreshViewport()
			return m, nil
		}

		// Update slash command suggestions (check first line only).
		value := m.input.Value()
		firstLine := strings.SplitN(value, "\n", 2)[0]

		if strings.HasPrefix(firstLine, "/") {
			m.showCommands = true
			m.filteredCommands = FilterCommands(firstLine)
		} else {
			m.showCommands = false
			m.filteredCommands = nil
			m.selectedCommand = 0
		}
		m.updateViewportSize()

		// Navigate the command palette.
		if m.showCommands {
			switch msg.String() {

			case "up":
				if m.selectedCommand > 0 {
					m.selectedCommand--
				}
				return m, nil

			case "down":
				if m.selectedCommand < len(m.filteredCommands)-1 {
					m.selectedCommand++
				}
				return m, nil

			case "esc":
				m.showCommands = false
				m.filteredCommands = nil
				m.selectedCommand = 0
				return m, nil
			}
		}
		switch msg.String() {

		case "enter":
			prompt := strings.TrimSpace(m.input.Value())

			if prompt == "" {
				return m, nil
			}

			// Execute selected command from palette.
			if m.showCommands && len(m.filteredCommands) > 0 {
				cmd := m.filteredCommands[m.selectedCommand]

				m.showCommands = false
				m.filteredCommands = nil
				m.selectedCommand = 0
				m.input.SetValue("")

				return m, cmd.Execute(&m)
			}

			// Execute typed slash command.
			if strings.HasPrefix(prompt, "/") {
				if command, ok := FindCommand(prompt); ok {
					m.input.SetValue("")
					return m, command.Execute(&m)
				}

				m.conversation.AddAssistant("Unknown command: " + prompt)
				m.input.SetValue("")
				return m, nil
			}

			// Normal AI message.
			m.conversation.AddUser(prompt)
			m.refreshViewport()
			m.loading = true

			m.input.SetValue("")

			// Create an empty assistant message.
			// Streaming chunks will be appended to this.
			m.conversation.AddAssistant("")
			m.refreshViewport()

			m.sendMessage(prompt)

			return m, nil
		}

	case StreamingMsg:

		if msg.Err != nil {
			m.loading = false

			if len(m.conversation.Messages) > 0 {
				m.conversation.Messages[len(m.conversation.Messages)-1].Content =
					"Error: " + msg.Err.Error()

				m.refreshViewport()
			}

			return m, nil
		}

		if len(m.conversation.Messages) > 0 {
			last := &m.conversation.Messages[len(m.conversation.Messages)-1]

			if last.Role == "assistant" {
				last.Content += msg.Text
				m.refreshViewport()
			}
		}

		if msg.Done {
			m.loading = false
		}

		return m, nil

	case switchProviderSuccessMsg:
		m.cfg.SetProvider(msg.Provider)
		m.cfg.SetModel(msg.Provider, msg.Model)
		_ = config.Save(m.cfg)

		provider, err := ai.New(m.cfg)
		if err == nil {
			m.provider = provider
			manager := tools.NewManager()
			defaults.RegisterAll(manager)
			m.agent = agent.New(provider, manager, agent.Config{})
		}

		m.conversation.AddAssistant(fmt.Sprintf("✨ Provider switched to **%s** (model: `%s`)", msg.Provider, msg.Model))
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchModelSuccessMsg:
		m.cfg.SetModel(m.cfg.AI.Provider, msg.Model)
		_ = config.Save(m.cfg)

		provider, err := ai.New(m.cfg)
		if err == nil {
			m.provider = provider
			manager := tools.NewManager()
			defaults.RegisterAll(manager)
			m.agent = agent.New(provider, manager, agent.Config{})
		}

		m.conversation.AddAssistant(fmt.Sprintf("✨ Model switched to `%s`", msg.Model))
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchCancelMsg:
		m.conversation.AddAssistant("❌ Provider/model selection cancelled.")
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchErrorMsg:
		m.conversation.AddAssistant(fmt.Sprintf("❌ Error: %v", msg.Err))
		m.refreshViewport()
		m.loading = false
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
