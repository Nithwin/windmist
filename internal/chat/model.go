package chat

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the WindMist application.
type Model struct {
	cfg *config.Config

	provider ai.Provider

	conversation Conversation

	input textinput.Model

	showSplash bool

	showCommands     bool
	filteredCommands []Command
	selectedCommand  int

	loading   bool
	streaming bool

	viewport viewport.Model

	markdown *ui.MarkdownRenderer

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

	renderer, err := ui.NewMarkdownRenderer()
	if err != nil {
		panic(err)
	}

	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = "message WindMist..."
	input.Focus()
	input.CharLimit = 0
	input.Width = 60
	vp := viewport.New(0, 0)

	return Model{
		cfg:          cfg,
		provider:     provider,
		conversation: Conversation{},
		input:        input,

		showSplash: true,

		showCommands:     false,
		filteredCommands: nil,
		selectedCommand:  0,

		loading:   false,
		streaming: false,

		viewport: vp,

		markdown: renderer,
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
