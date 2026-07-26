package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/store"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
	"github.com/Nithwin/WindMist/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleEventMsg(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ApprovalRequestMsg:
		m.waitingApproval = true
		m.approvalCommand = msg.Command
		m.approvalChan = msg.ResponseChan
		m.refreshViewport()
		return m, nil

	case switchModeSuccessMsg:
		m.session.AgentMode = msg.Mode
		if m.store != nil {
			_ = m.store.UpdateSession(m.session)
		}
		
		// Update Agent config mode
		m.agent = agent.New(m.provider, m.agent.Manager(), agent.Config{
			Store:     m.store,
			SessionID: m.session.ID,
			Mode:      m.session.AgentMode,
		})

		m.conversation.AddAssistant(fmt.Sprintf("✨ Switched Agent Mode to **%s**", strings.ToUpper(msg.Mode)))
		m.refreshViewport()
		return m, nil

	case createNewSessionMsg:
		activeModel, _ := m.cfg.ActiveModel()
		sess := &store.Session{
			ID:          fmt.Sprintf("sess_%d", time.Now().Unix()),
			Title:       "New Session",
			ProjectPath: ".",
			Provider:    m.cfg.AI.Provider,
			Model:       activeModel,
			AgentMode:   "auto",
		}
		if m.store != nil {
			_ = m.store.CreateSession(sess)
		}

		m.session = sess
		m.agent = agent.New(m.provider, m.agent.Manager(), agent.Config{
			Store:     m.store,
			SessionID: sess.ID,
			Mode:      sess.AgentMode,
		})

		m.conversation.Clear()
		m.conversation.AddAssistant("✨ Started a new session.")
		m.refreshViewport()
		return m, nil

	case undoFileChangeMsg:
		if m.store == nil || m.session == nil {
			m.conversation.AddAssistant("❌ Persistence not enabled.")
			m.refreshViewport()
			return m, nil
		}

		changes, err := m.store.GetLastBatchForUndo(m.session.ID)
		if err != nil || len(changes) == 0 {
			m.conversation.AddAssistant("❌ No file changes found to undo.")
			m.refreshViewport()
			return m, nil
		}

		for _, change := range changes {
			if change.ChangeType == "create" {
				_ = os.Remove(change.FilePath)
			} else {
				_ = os.WriteFile(change.FilePath, []byte(change.BeforeContent), 0644)
			}
		}

		_ = m.store.SetBatchUndoneState(m.session.ID, changes[0].BatchID, true)

		m.conversation.AddAssistant(fmt.Sprintf("⏮️ **Undid %d file edit(s)**", len(changes)))
		m.refreshViewport()
		return m, nil

	case redoFileChangeMsg:
		if m.store == nil || m.session == nil {
			m.conversation.AddAssistant("❌ Persistence not enabled.")
			m.refreshViewport()
			return m, nil
		}

		changes, err := m.store.GetNextBatchForRedo(m.session.ID)
		if err != nil || len(changes) == 0 {
			m.conversation.AddAssistant("❌ No file changes found to redo.")
			m.refreshViewport()
			return m, nil
		}

		for _, change := range changes {
			if change.ChangeType == "delete" {
				_ = os.Remove(change.FilePath)
			} else {
				_ = os.WriteFile(change.FilePath, []byte(change.AfterContent), 0644)
			}
		}
		
		_ = m.store.SetBatchUndoneState(m.session.ID, changes[0].BatchID, false)

		m.conversation.AddAssistant(fmt.Sprintf("⏭️ **Redid %d file edit(s)**", len(changes)))
		m.refreshViewport()
		return m, nil

	case switchSessionSuccessMsg:
		sess, err := m.store.GetSession(msg.SessionID)
		if err != nil {
			m.conversation.AddAssistant(fmt.Sprintf("❌ Error loading session: %v", err))
			m.refreshViewport()
			return m, nil
		}

		m.session = sess
		m.agent = agent.New(m.provider, m.agent.Manager(), agent.Config{
			Store:     m.store,
			SessionID: sess.ID,
			Mode:      sess.AgentMode,
		})

		m.conversation.Clear()
		initialMessages := m.getInitialMessages()
		for _, msg := range initialMessages {
			if msg.Role == ai.RoleUser {
				m.conversation.AddUser(msg.Content)
			} else if msg.Role == ai.RoleAssistant {
				content := msg.Content
				if len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						content += fmt.Sprintf("\n*(Tool Call: %s)*", tc.Name)
					}
				}
				m.conversation.AddAssistant(content)
			} else if msg.Role == ai.RoleTool {
				content := ""
				for _, tr := range msg.ToolResults {
					content += fmt.Sprintf("\n*(Tool Result: %s)*", tr.Name)
				}
				m.conversation.AddAssistant(content)
			}
		}

		m.conversation.AddAssistant(fmt.Sprintf("✨ Loaded session: **%s**", sess.Title))
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchProviderSuccessMsg:
		m.cfg.SetProvider(msg.Provider)
		m.cfg.SetModel(msg.Provider, msg.Model)
		_ = config.Save(m.cfg)

		provider, err := ai.New(m.cfg)
		if err == nil {
			m.provider = provider
			manager := tools.NewManager()
			defaults.RegisterAll(manager, func(cmd string) bool {
				if program == nil {
					return false
				}
				ch := make(chan bool)
				program.Send(ApprovalRequestMsg{Command: cmd, ResponseChan: ch})
				return <-ch
			}, m.cfg)
			m.agent = agent.New(provider, manager, agent.Config{})
		}

		m.conversation.AddAssistant(fmt.Sprintf("✨ Provider switched to **%s** (model: `%s`)", msg.Provider, msg.Model))
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchModelSuccessMsg:
		m.cfg.SetModel(m.cfg.AI.Provider, msg.Model)
		_ = config.Save(m.cfg)

		provider, err := ai.New(m.cfg)
		if err == nil {
			m.provider = provider
			manager := tools.NewManager()
			defaults.RegisterAll(manager, func(cmd string) bool {
				if program == nil {
					return false
				}
				ch := make(chan bool)
				program.Send(ApprovalRequestMsg{Command: cmd, ResponseChan: ch})
				return <-ch
			}, m.cfg)
			m.agent = agent.New(provider, manager, agent.Config{})
		}

		m.conversation.AddAssistant(fmt.Sprintf("✨ Model switched to `%s`", msg.Model))
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchSubagentSuccessMsg:
		m.cfg.SubAgent.Provider = msg.Provider
		m.cfg.SubAgent.Model = msg.Model
		_ = config.Save(m.cfg)

		// Re-register tools with the new config so sub-agent picks it up
		manager := tools.NewManager()
		defaults.RegisterAll(manager, func(cmd string) bool {
			if program == nil {
				return false
			}
			ch := make(chan bool)
			program.Send(ApprovalRequestMsg{Command: cmd, ResponseChan: ch})
			return <-ch
		}, m.cfg)
		m.agent = agent.New(m.provider, manager, agent.Config{})

		if msg.Provider == "" {
			m.conversation.AddAssistant("✨ Sub-Agent reset to Auto (will use fast fallback or main model).")
		} else {
			m.conversation.AddAssistant(fmt.Sprintf("✨ Sub-Agent switched to **%s** (model: `%s`)", msg.Provider, msg.Model))
		}
		
		m.refreshViewport()
		m.loading = false
		return m, nil

	case mcpInstallSuccessMsg:
		m.conversation.AddAssistant(fmt.Sprintf("🔌 Successfully installed and configured MCP Server: **%s**.\nPlease restart WindMist (using `/exit`) to automatically discover the new tools from this server.", msg.Name))
		m.refreshViewport()
		return m, nil

	case switchThemeSuccessMsg:
		m.cfg.SetTheme(msg.Theme)
		_ = config.Save(m.cfg)

		customDir := ""
		cfgDir, err := config.ConfigDir()
		if err == nil {
			customDir = filepath.Join(cfgDir, "themes")
		}

		err = ui.LoadTheme(msg.Theme, customDir)
		if err != nil {
			m.conversation.AddAssistant(fmt.Sprintf("❌ Failed to load theme: %v", err))
		} else {
			m.UpdateInputStyles()
			m.conversation.AddAssistant(fmt.Sprintf("✨ Theme switched to **%s**", msg.Theme))
		}

		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchCancelMsg:
		m.conversation.AddAssistant("❌ Provider/model selection cancelled.")
		m.refreshViewport()
		m.loading = false
		return m, nil

	case switchErrorMsg:
		m.conversation.AddAssistant(fmt.Sprintf("❌ Error: %v", msg.Err))
		m.refreshViewport()
		m.loading = false
		return m, nil
	}

	return m, nil
}
