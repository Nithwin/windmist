package chat

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all user interactions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

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

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m, tea.Quit

		case "enter":
			if err := m.sendMessage(); err != nil {
				m.conversation.AddAssistant("Error: " + err.Error())
			}
			m.input.SetValue("")
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}
