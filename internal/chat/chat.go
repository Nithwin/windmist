package chat

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Nithwin/WindMist/internal/ai"
)

// sendMessage starts an AI request in the background.
func (m Model) sendMessage(prompt string) tea.Cmd {
	return func() tea.Msg {

		req := &ai.GenerateRequest{
			Messages: []ai.Message{
				{
					Role:    ai.RoleUser,
					Content: prompt,
				},
			},
		}

		resp, err := m.provider.Generate(context.Background(), req)
		if err != nil {
			return ResponseMsg{
				Err: err,
			}
		}

		return ResponseMsg{
			Text: resp.Text,
		}
	}
}