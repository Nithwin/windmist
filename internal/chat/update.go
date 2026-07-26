package chat

import (
	tea "github.com/charmbracelet/bubbletea"
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

	// All other custom events (Session, Agent Mode, Undo/Redo, Models)
	case ApprovalRequestMsg, switchModeSuccessMsg, createNewSessionMsg,
		undoFileChangeMsg, redoFileChangeMsg, switchSessionSuccessMsg,
		switchProviderSuccessMsg, switchModelSuccessMsg, switchSubagentSuccessMsg, switchThemeSuccessMsg, switchCancelMsg, switchErrorMsg:
		
		var evtCmd tea.Cmd
		m, evtCmd = m.handleEventMsg(msg)
		return m, evtCmd
	}

	// Update text input for other key events that don't match the main handler
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
