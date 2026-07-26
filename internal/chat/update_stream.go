package chat

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleStreamMsg(msg StreamingMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.loading = false
		m.streaming = false
		m.spinnerFrame = 0
		m.streamTokens = m.streamTokens // preserve for display

		if len(m.conversation.Messages) > 0 {
			m.conversation.Messages[len(m.conversation.Messages)-1].Content =
				"Error: " + msg.Err.Error()

			m.refreshViewport()
		}

		// Check for queued message even on error
		if m.queuedMessage != "" {
			return m, m.processQueuedMessage()
		}

		return m, nil
	}

	if len(m.conversation.Messages) > 0 {
		last := &m.conversation.Messages[len(m.conversation.Messages)-1]

		if last.Role == "assistant" {
			last.Content += msg.Text
			m.refreshViewport()
		}
	}

	// Update real-time token counter from streaming usage data
	if msg.Usage.TotalTokens > 0 {
		m.streamTokens = msg.Usage
	}

	if msg.Done {
		m.loading = false
		m.streaming = false
		m.spinnerFrame = 0
		m.responseTime = msg.Duration

		if m.session != nil {
			m.session.TokenCount += msg.Usage.TotalTokens
			// Rough cost estimation logic could go here or in a separate function
			// m.session.CostEstimate += calculateCost(...)

			// Save to DB
			if m.store != nil {
				_ = m.store.UpdateSession(m.session)
			}
		}

		// Reset stream tokens for next request
		m.streamTokens = msg.Usage

		// Auto-send queued message if one exists
		if m.queuedMessage != "" {
			return m, m.processQueuedMessage()
		}
	}

	return m, nil
}

// processQueuedMessage takes the queued message and sends it as a new AI request.
func (m *Model) processQueuedMessage() tea.Cmd {
	prompt := m.queuedMessage
	m.queuedMessage = ""

	m.inputHistory = append(m.inputHistory, prompt)
	m.historyIndex = len(m.inputHistory)

	m.conversation.AddUser(prompt)
	m.refreshViewport()
	m.loading = true
	m.streaming = true
	m.streamTokens = m.streamTokens // Reset for new request

	// Create an empty assistant message for streaming
	m.conversation.AddAssistant("")
	m.refreshViewport()

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	return m.sendMessageCmd(ctx, prompt)
}
