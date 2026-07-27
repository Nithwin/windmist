package chat

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/rag"
	"github.com/Nithwin/WindMist/internal/remote"
	"github.com/Nithwin/WindMist/internal/remote/telegram"
	"github.com/Nithwin/WindMist/internal/store"
	"github.com/Nithwin/WindMist/internal/tools"
	toolagent "github.com/Nithwin/WindMist/internal/tools/agent"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	"github.com/charmbracelet/bubbles/list"
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
	spinnerFrame int
	responseTime time.Duration

	// Input queuing: let user type next message while loading
	queuedMessage string

	// Streaming token counter: updated in real-time
	streamTokens ai.Usage

	waitingApproval bool
	approvalCommand string
	approvalChan    chan bool

	viewport viewport.Model

	markdown *ui.MarkdownRenderer

	// RAG components
	ragIndexer *rag.Indexer

	// Memory components
	summarizer *agent.Summarizer

	// Inline Selector state
	showSelector bool
	selectorList list.Model
	onSelect     func(selector.Option) tea.Cmd
	onCancel     func() tea.Cmd

	// Inline Prompt state
	inlinePrompt   string
	onPromptSubmit func(string) tea.Cmd
	isPassword     bool

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

	// Initialize RAG System
	home, _ := os.UserHomeDir()
	ragDbPath := filepath.Join(home, ".windmist", "rag.db")
	ragStore, err := rag.NewDocumentStore(ragDbPath)
	if err != nil {
		return Model{}, fmt.Errorf("failed to initialize RAG store: %w", err)
	}
	ragEmbedder := rag.NewTFIDFEmbedder(512)
	ragSearcher := rag.NewSearcher(ragStore, ragEmbedder)
	ragIndexer := rag.NewIndexer(ragStore, ragEmbedder)

	// Rebuild vocabulary on startup if we have indexed chunks
	if chunks, err := ragStore.GetAllChunks(); err == nil && len(chunks) > 0 {
		docs := make([]string, len(chunks))
		for i, c := range chunks {
			docs[i] = c.Content
		}
		ragEmbedder.BuildVocabulary(docs)
	}

	// Register Semantic Search Tool
	manager.Register(toolagent.NewSemanticSearchTool(ragSearcher))

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
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 3

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

		ragIndexer: ragIndexer,
		summarizer: agent.NewSummarizer(provider, 8000),
	}

	model.UpdateInputStyles()

	model.updateViewportSize()

	return model, nil
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, func() tea.Msg {
		return WorkspaceFilesMsg{Files: getWorkspaceFiles()}
	})

	if m.cfg.Remote.Telegram.Enabled && m.cfg.Remote.Telegram.BotToken != "" {
		if remote.GetHub() == nil {
			remote.InitHub(&m.cfg.Remote)
		}

		tController, err := telegram.New(m.cfg.Remote.Telegram)
		if err == nil {
			err = remote.GetHub().Register(tController)
			if err == nil {
				cmds = append(cmds, listenRemoteCmd())
			}
		}
	}

	return tea.Batch(cmds...)
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
