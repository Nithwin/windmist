package chat

import (
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

// Run starts the WindMist Bubble Tea application.
func Run() error {
	// Silence standard logger to prevent third-party libraries (e.g. tgbotapi)
	// from printing directly to stderr and corrupting the TUI.
	log.SetOutput(io.Discard)

	model, err := New()
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	program = p

	_, err = p.Run()
	return err
}
