package chat

import (
	"context"

	"github.com/Nithwin/WindMist/internal/ai"
)

// sendMessage starts streaming the AI response.
func (m Model) sendMessage(prompt string) {
	go func() {

		req := &ai.GenerateRequest{
			Messages: []ai.Message{
				{
					Role:    ai.RoleUser,
					Content: prompt,
				},
			},
		}

		err := m.provider.Stream(
			context.Background(),
			req,
			func(chunk string) {

				program.Send(StreamingMsg{
					Text: chunk,
				})

			},
		)

		if err != nil {
			program.Send(StreamingMsg{
				Err: err,
				Done: true,
			})
			return
		}

		program.Send(StreamingMsg{
			Done: true,
		})

	}()
}