package chat

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/store"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the WindMist application.
type Model struct {
	cfg *config.Config

	provider ai.Provider
	agent    *agent.Agent
	store    *store.Store
	session  *store.Session

	conversation Conversation

	input        textarea.Model
	inputHistory []string
	historyIndex int

	showSplash bool

	showCommands     bool
	filteredCommands []Command
	selectedCommand  int

	showFilePicker bool
	workspaceFiles []string
	filteredFiles  []string
	selectedFile   int

	loading      bool
	streaming    bool
	responseTime time.Duration

	waitingApproval bool
	approvalCommand string
	approvalChan    chan bool

	viewport viewport.Model

	markdown *ui.MarkdownRenderer

	width  int
	height int

	cancel context.CancelFunc
}

// New creates a new Bubble Tea model.
func New() (Model, error) {
	cfg, err := config.Load()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	customDir := ""
	cfgDir, err := config.ConfigDir()
	if err == nil {
		customDir = filepath.Join(cfgDir, "themes")
	}
	_ = ui.LoadTheme(cfg.UI.Theme, customDir)

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
	}, cfg)
	dbStore, err := store.NewStore()
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize db store: %w", err)
	}

	// For now, create a new session on startup
	// Later we can implement logic to load an existing session
	// using the /session commands
	activeModel, _ := cfg.ActiveModel()
	sess := &store.Session{
		ID:          fmt.Sprintf("sess_%d", time.Now().Unix()),
		Title:       "New Session",
		ProjectPath: ".",
		Provider:    cfg.AI.Provider,
		Model:       activeModel,
		AgentMode:   "auto",
	}
	_ = dbStore.CreateSession(sess)

	ag := agent.New(provider, manager, agent.Config{
		Store:     dbStore,
		SessionID: sess.ID,
		Mode:      sess.AgentMode,
	})

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
	
	vp := viewport.New(0, 0)
	
	model := Model{
		cfg:          cfg,
		provider:     provider,
		agent:        ag,
		store:        dbStore,
		session:      sess,
		conversation: Conversation{},
		input:        ta,
		inputHistory: make([]string, 0),
		historyIndex: 0,

		showSplash: true,

		showCommands:     false,
		filteredCommands: nil,
		selectedCommand:  0,

		loading:   false,
		streaming: false,

		waitingApproval: false,

		viewport: vp,

		markdown: renderer,
	}

	model.UpdateInputStyles()

	model.updateViewportSize()

	// Async load files so startup is fast
	go func() {
		files := getWorkspaceFiles()
		// We'd ideally send a Msg, but this is fine for now
		model.workspaceFiles = files
	}()

	return model, nil

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

// UpdateInputStyles applies the current UI colors to the textarea input.
func (m *Model) UpdateInputStyles() {
	plain := ui.BaseStyle
	m.input.FocusedStyle.Base = plain.Foreground(ui.White)
	m.input.FocusedStyle.Text = plain.Foreground(ui.White)
	m.input.FocusedStyle.CursorLine = plain.Foreground(ui.White)
	m.input.FocusedStyle.Placeholder = plain.Foreground(ui.Muted)
	m.input.FocusedStyle.EndOfBuffer = plain.Foreground(ui.Muted)
	m.input.FocusedStyle.Prompt = plain
	m.input.BlurredStyle.Base = plain.Foreground(ui.MutedLight)
	m.input.BlurredStyle.Text = plain.Foreground(ui.MutedLight)
	m.input.BlurredStyle.Placeholder = plain.Foreground(ui.Muted)
	m.input.BlurredStyle.CursorLine = plain
	m.input.BlurredStyle.Prompt = plain
}
