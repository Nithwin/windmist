package chat

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all user interactions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m, tea.Quit

		case "enter":
			// For now just clear the input.
			// Later we'll send this to Gemini.
			m.input.SetValue("")
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}