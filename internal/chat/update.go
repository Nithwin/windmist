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

		switch msg.String() {

		case "enter":
			prompt := strings.TrimSpace(m.input.Value())

			if prompt == "" {
				return m, nil
			}

			// Is it a slash command?
			if strings.HasPrefix(prompt, "/") {
				if command, ok := FindCommand(prompt); ok {
					m.input.SetValue("")
					return m, command.Execute(&m)
				}

				m.conversation.AddAssistant("Unknown command: " + prompt)
				m.input.SetValue("")
				return m, nil
			}

			// Normal AI message
			if err := m.sendMessage(); err != nil {
				m.conversation.AddAssistant("Error: " + err.Error())
			}
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
