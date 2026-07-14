package chat

import (
	"context"

	"github.com/Nithwin/WindMist/internal/ai"
)

func (m *Model) sendMessage() error {
	prompt := m.input.Value()

	if prompt == "" {
		return nil
	}

	m.conversation.AddUser(prompt)

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
		return err
	}

	m.conversation.AddAssistant(resp.Text)

	m.input.SetValue("")

	return nil
}