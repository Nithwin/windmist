package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
)

// sendMessage starts running the agent request.
func (m Model) sendMessage(ctx context.Context, prompt string) {
	go func() {
		initialMessages := m.getInitialMessages()
		res, err := m.agent.Run(ctx, initialMessages, prompt, func(s string) {
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

func (m Model) getInitialMessages() []ai.Message {
	if m.store == nil || m.session == nil {
		return nil
	}

	storeMsgs, err := m.store.GetMessagesBySession(m.session.ID)
	if err != nil || len(storeMsgs) == 0 {
		return nil
	}

	var msgs []ai.Message
	for _, sm := range storeMsgs {
		msg := ai.Message{
			Role:    ai.Role(sm.Role),
			Content: sm.Content,
		}

		if sm.ToolCalls != "" {
			var calls []ai.ToolCall
			if err := json.Unmarshal([]byte(sm.ToolCalls), &calls); err == nil {
				msg.ToolCalls = calls
			}
		}

		if sm.ToolResults != "" {
			var res []ai.ToolResult
			if err := json.Unmarshal([]byte(sm.ToolResults), &res); err == nil {
				msg.ToolResults = res
			}
		}

		msgs = append(msgs, msg)
	}

	return msgs
}
