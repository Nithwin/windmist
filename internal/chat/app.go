package chat

import tea "github.com/charmbracelet/bubbletea"

var program *tea.Program

// Run starts the WindMist Bubble Tea application.
func Run() error {
	model, err := New()
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	program = p

	_, err = p.Run()
	return err
}
