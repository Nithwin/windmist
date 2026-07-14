package chat

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the WindMist application.
type Model struct {
	cfg *config.Config

	provider ai.Provider

	input textinput.Model

	showSplash bool

	width  int
	height int
}

// New creates a new Bubble Tea model.
func New() Model {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	provider, err := ai.New(cfg)
	if err != nil {
		panic(err)
	}

	input := textinput.New()
	input.Prompt = ui.PromptStyle.Render("❯ ")
	input.Placeholder = "Ask anything..."
	input.Focus()
	input.CharLimit = 0
	input.Width = 80

	return Model{
		cfg:         cfg,
		provider:    provider,
		input:       input,
		showSplash:  true,
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
