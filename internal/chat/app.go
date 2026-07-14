package chat

import tea "github.com/charmbracelet/bubbletea"

var program *tea.Program

// Run starts the WindMist Bubble Tea application.
func Run() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
	)

	program = p

	_, err := p.Run()
	return err
}