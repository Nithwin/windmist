package chat

import (
	"context"
)

// sendMessage starts running the agent request.
func (m Model) sendMessage(prompt string) {
	go func() {
		res, err := m.agent.Run(context.Background(), prompt)
		if err != nil {
			program.Send(StreamingMsg{
				Err:  err,
				Done: true,
			})
			return
		}

		program.Send(StreamingMsg{
			Text: res.Content,
			Done: true,
		})
	}()
}
