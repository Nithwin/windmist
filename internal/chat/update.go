package chat

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/remote"
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles all user interactions and routes them to specific handlers.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(m.MaxContentWidth())
		m.refreshViewport()
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case StreamingMsg:
		return m.handleStreamMsg(msg)

	case ResponseMsg:
		if msg.Err != nil {
			m.conversation.AddAssistant(fmt.Sprintf("❌ Error: %v", msg.Err))
		} else {
			m.conversation.AddAssistant(msg.Text)
		}
		m.refreshViewport()
		return m, nil

	case spinnerTickMsg:
		if m.loading {
			m.spinnerFrame++
			m.refreshViewport()
			return m, spinnerTickCmd()
		}
		return m, nil

	case tea.MouseMsg:
		// Route mouse wheel events to the viewport for scrolling
		if !m.showSplash && !m.showSelector {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				m.viewport.ScrollUp(3)
				return m, nil
			case tea.MouseButtonWheelDown:
				m.viewport.ScrollDown(3)
				return m, nil
			}
		}
		return m, nil

	case WorkspaceFilesMsg:
		m.workspaceFiles = msg.Files
		return m, nil

	// All other custom events (Session, Agent Mode, Undo/Redo, Models)
	case ApprovalRequestMsg, switchModeSuccessMsg, createNewSessionMsg,
		undoFileChangeMsg, redoFileChangeMsg, switchSessionSuccessMsg,
		switchProviderSuccessMsg, switchModelSuccessMsg, switchSubagentSuccessMsg, switchThemeSuccessMsg, mcpInstallSuccessMsg, setAPIKeySuccessMsg, switchCancelMsg, switchErrorMsg, indexWorkspaceMsg, compactConversationMsg:

		var evtCmd tea.Cmd
		m, evtCmd = m.handleEventMsg(msg)
		return m, evtCmd

	case remoteCommandMsg:
		switch msg.Type {
		case "list_providers":
			// We can generate a string of providers
			var list string
			for p := range m.cfg.Providers {
				list += "- " + p + "\n"
			}
			if hub := remote.GetHub(); hub != nil {
				hub.Broadcast <- "Available Providers:\n" + list
			}
			return m, listenRemoteCmd()

		case "list_models":
			if hub := remote.GetHub(); hub != nil {
				pConfig, _ := m.cfg.ActiveProvider()
				// Send active model and custom models
				models := pConfig.Model + " (active)"
				for _, cm := range m.cfg.CustomModels[m.cfg.AI.Provider] {
					models += "\n- " + cm
				}
				hub.Broadcast <- "Models for " + m.cfg.AI.Provider + ":\n- " + models
			}
			return m, listenRemoteCmd()

		case "provider":
			if len(msg.Args) > 0 {
				newProvider := msg.Args[0]
				if err := m.cfg.SetProvider(newProvider); err == nil {
					_ = config.Save(m.cfg)
					if hub := remote.GetHub(); hub != nil {
						hub.Broadcast <- "✅ Switched provider to " + newProvider
					}
					// Fire the same command as TUI
					return m, tea.Batch(
						func() tea.Msg { return switchProviderSuccessMsg{Provider: newProvider} },
						listenRemoteCmd(),
					)
				}
				if hub := remote.GetHub(); hub != nil {
					hub.Broadcast <- "❌ Failed to switch provider: " + newProvider
				}
			}
			return m, listenRemoteCmd()

		case "model":
			if len(msg.Args) > 0 {
				newModel := msg.Args[0]
				if err := m.cfg.SetModel(m.cfg.AI.Provider, newModel); err == nil {
					_ = config.Save(m.cfg)
					if hub := remote.GetHub(); hub != nil {
						hub.Broadcast <- "✅ Switched model to " + newModel
					}
					return m, tea.Batch(
						func() tea.Msg { return switchModelSuccessMsg{Model: newModel} },
						listenRemoteCmd(),
					)
				}
				if hub := remote.GetHub(); hub != nil {
					hub.Broadcast <- "❌ Failed to switch model: " + newModel
				}
			}
			return m, listenRemoteCmd()

		case "ask":
			if len(msg.Args) > 0 {
				prompt := msg.Args[0]
				m.input.SetValue(prompt)
				// Simulate pressing enter
				return m, tea.Batch(func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEnter} }, listenRemoteCmd())
			}
			return m, listenRemoteCmd()
		}
		return m, listenRemoteCmd()

	case showInlineSelectorMsg:
		m.showSelector = true
		m.onSelect = msg.OnSelect
		m.onCancel = msg.OnCancel

		items := make([]list.Item, len(msg.Options))
		for i, opt := range msg.Options {
			items[i] = opt
		}

		d := list.NewDefaultDelegate()
		d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(ui.Cyan).BorderForeground(ui.Cyan)
		d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(ui.Cyan).BorderForeground(ui.Cyan)

		m.selectorList = list.New(items, d, 80, 20)
		m.selectorList.Title = msg.Title
		m.selectorList.SetShowStatusBar(false)
		m.selectorList.SetFilteringEnabled(true)
		m.selectorList.Styles.Title = lipgloss.NewStyle().Background(ui.Purple).Foreground(ui.White).Padding(0, 1)

		// Set size
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.selectorList.SetSize(m.width-h, m.height-v)

		return m, nil

	case showInlinePromptMsg:
		m.inlinePrompt = msg.Prompt
		m.isPassword = msg.IsPassword
		m.onPromptSubmit = msg.OnSubmit
		m.input.Reset()
		return m, nil

	// Update text input for other key events that don't match the main handler
	default:
		if m.showSelector {
			var listCmd tea.Cmd
			m.selectorList, listCmd = m.selectorList.Update(msg)
			return m, listCmd
		}

		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}
