package chat

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleStreamMsg(msg StreamingMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.loading = false

		if len(m.conversation.Messages) > 0 {
			m.conversation.Messages[len(m.conversation.Messages)-1].Content =
				"Error: " + msg.Err.Error()

			m.refreshViewport()
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

	if msg.Done {
		m.loading = false
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
	}

	return m, nil
}
