package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nithwin/WindMist/internal/ai"
)

func (m Model) sendMessage(ctx context.Context, prompt string) {
	// Auto-title the session if it's the first message
	if m.session != nil && m.session.Title == "New Session" && m.store != nil {
		go func() {
			titleReq := &ai.GenerateRequest{
				System: "You are an AI that creates extremely short, 2-4 word titles for chat sessions based on the user's first prompt. Do not use punctuation. Do not use quotes. Keep it lowercase.",
				Messages: []ai.Message{
					{Role: ai.RoleUser, Content: prompt},
				},
				MaxTokens: 20,
			}
			resp, err := m.provider.Generate(context.Background(), titleReq)
			if err == nil && resp.Text != "" {
				m.session.Title = resp.Text
				_ = m.store.UpdateSession(m.session)
			}
		}()
	}

	go func() {
		startTime := time.Now()
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

		duration := time.Since(startTime)
		
		// Agent loop completed
		program.Send(StreamingMsg{
			Text:     "\n\n(Finished in " + fmt.Sprintf("%d turns", res.Turns) + ")",
			Done:     true,
			Usage:    res.Usage,
			Duration: duration,
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
