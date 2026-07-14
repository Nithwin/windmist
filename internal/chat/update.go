package chat

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all user interactions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
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

			return m, nil
		}

		// Update slash command suggestions.
		value := m.input.Value()

		if strings.HasPrefix(value, "/") {
			m.showCommands = true
			m.filteredCommands = FilterCommands(value)
		} else {
			m.showCommands = false
			m.filteredCommands = nil
			m.selectedCommand = 0
		}

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

			m.loading = true

			m.input.SetValue("")

			// Create an empty assistant message.
			// Streaming chunks will be appended to this.
			m.conversation.AddAssistant("")

			m.sendMessage(prompt)

			return m, nil
		}

	case StreamingMsg:

		if msg.Err != nil {
			m.loading = false

			if len(m.conversation.Messages) > 0 {
				m.conversation.Messages[len(m.conversation.Messages)-1].Content =
					"Error: " + msg.Err.Error()
			}

			return m, nil
		}

		if len(m.conversation.Messages) > 0 {
			last := &m.conversation.Messages[len(m.conversation.Messages)-1]

			if last.Role == "assistant" {
				last.Content += msg.Text
			}
		}

		if msg.Done {
			m.loading = false
		}

		return m, nil
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
