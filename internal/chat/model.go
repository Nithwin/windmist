package chat

import (
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the WindMist application.
type Model struct {
	cfg *config.Config

	provider ai.Provider

	conversation Conversation

	input textarea.Model

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

	ta := textarea.New()
	ta.Placeholder = "Message WindMist... (Enter to send, Shift+Enter for new line)"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(76)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.Prompt = ""

	// Clean minimal style — no borders, transparent background
	plain := lipgloss.NewStyle()
	ta.FocusedStyle.Base = plain.Foreground(ui.White)
	ta.FocusedStyle.CursorLine = plain.Foreground(ui.White)
	ta.FocusedStyle.Placeholder = plain.Foreground(ui.Muted)
	ta.FocusedStyle.EndOfBuffer = plain.Foreground(ui.Muted)
	ta.BlurredStyle.Base = plain.Foreground(ui.MutedLight)
	ta.BlurredStyle.Placeholder = plain.Foreground(ui.Muted)
	ta.BlurredStyle.CursorLine = plain

	vp := viewport.New(0, 0)

	return Model{
		cfg:          cfg,
		provider:     provider,
		conversation: Conversation{},
		input:        ta,

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
	return textarea.Blink
}
