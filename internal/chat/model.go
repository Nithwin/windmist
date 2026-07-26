package chat

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
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
	agent    *agent.Agent

	conversation Conversation

	input textarea.Model

	showSplash bool

	showCommands     bool
	filteredCommands []Command
	selectedCommand  int

	loading   bool
	streaming bool

	waitingApproval bool
	approvalCommand string
	approvalChan    chan bool

	viewport viewport.Model

	markdown *ui.MarkdownRenderer

	width  int
	height int
}

// New creates a new Bubble Tea model.
func New() (Model, error) {
	cfg, err := config.Load()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	provider, err := ai.New(cfg)
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	manager := tools.NewManager()
	defaults.RegisterAll(manager, func(cmd string) bool {
		if program == nil {
			return false
		}
		ch := make(chan bool)
		program.Send(ApprovalRequestMsg{
			Command:      cmd,
			ResponseChan: ch,
		})
		return <-ch
	})
	ag := agent.New(provider, manager, agent.Config{})

	renderer, err := ui.NewMarkdownRenderer()
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize markdown renderer: %w", err)
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
		agent:        ag,
		conversation: Conversation{},
		input:        ta,

		showSplash: true,

		showCommands:     false,
		filteredCommands: nil,
		selectedCommand:  0,

		loading:   false,
		streaming: false,

		waitingApproval: false,

		viewport: vp,

		markdown: renderer,
	}, nil
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// MaxContentWidth calculates the maximum width for the UI content based on the window size.
// We allow it to expand dynamically, capped at 120 columns for typography/readability.
func (m Model) MaxContentWidth() int {
	w := m.width - 4 // padding
	if w > 120 {
		return 120
	}
	if w < 40 {
		return 40
	}
	return w
}
