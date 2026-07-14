package chat

import (
	"github.com/Nithwin/WindMist/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the WindMist application.
type Model struct {
	cfg   *config.Config
	input textinput.Model
}

// New creates a new Bubble Tea model.
func New() Model {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	input := textinput.New()
	input.Placeholder = "Ask WindMist anything..."
	input.Focus()
	input.CharLimit = 0
	input.Width = 80

	return Model{
		cfg:   cfg,
		input: input,
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
