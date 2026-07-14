package chat

import tea "github.com/charmbracelet/bubbletea"

// Run starts the WindMist Bubble Tea application.
func Run() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}