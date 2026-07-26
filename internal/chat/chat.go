package chat

import (
	"context"
	"fmt"
)

// sendMessage starts running the agent request.
func (m Model) sendMessage(ctx context.Context, prompt string) {
	go func() {
		res, err := m.agent.Run(ctx, prompt, func(s string) {
			program.Send(StreamingMsg{
				Text: s,
			})
		})

		if err != nil {
			program.Send(StreamingMsg{
				Err:  err,
				Done: true,
			})
			return
		}

		// Agent loop completed
		program.Send(StreamingMsg{
			Text: "\n\n(Finished in " + fmt.Sprintf("%d turns", res.Turns) + ")",
			Done: true,
		})
	}()
}
