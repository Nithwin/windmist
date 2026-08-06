package chat

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/remote"
	tea "github.com/charmbracelet/bubbletea"
)

// spinnerTickMsg is sent periodically to animate the loading spinner.
type spinnerTickMsg struct{}

// spinnerTickCmd returns a tea.Cmd that fires a tick after a short delay.
func spinnerTickCmd() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

// listenRemoteCmd waits for incoming commands from the remote hub.
func listenRemoteCmd() tea.Cmd {
	return func() tea.Msg {
		hub := remote.GetHub()
		if hub == nil {
			return nil
		}
		cmd := <-hub.Incoming
		return remoteCommandMsg{
			Type: cmd.Type,
			Args: cmd.Args,
		}
	}
}

// sendMessageCmd returns a tea.Cmd that starts the AI request in a goroutine
// and immediately begins the spinner tick loop.
func (m Model) sendMessageCmd(ctx context.Context, prompt string) tea.Cmd {
	return tea.Batch(
		// Start the spinner tick loop
		spinnerTickCmd(),
		// Fire the AI request in a goroutine
		func() tea.Msg {
			startTime := time.Now()
			initialMessages := m.getInitialMessages()
			res, err := m.agent.Run(ctx, initialMessages, prompt, func(s string) {
				program.Send(StreamingMsg{
					Text: s,
				})
			})

			// Auto-title the session if it's the first message (Moved here to avoid concurrent API limits on Free Tier)
			if m.session != nil && m.session.Title == "New Session" && m.store != nil {
				sessionID := m.session.ID
				provider := m.provider
				go func() {
					// Small delay to ensure the main stream request is fully closed
					time.Sleep(1 * time.Second)
					titleReq := &ai.GenerateRequest{
						System: "You are an AI that creates extremely short, 2-4 word titles for chat sessions based on the user's first prompt. Do not use punctuation. Do not use quotes. Keep it lowercase.",
						Messages: []ai.Message{
							{Role: ai.RoleUser, Content: prompt},
						},
						MaxTokens: 20,
					}
					resp, err := provider.Generate(context.Background(), titleReq)
					if err == nil && resp.Text != "" && program != nil {
						program.Send(sessionTitleMsg{
							SessionID: sessionID,
							Title:     resp.Text,
						})
					}
				}()
			}

			if err != nil {
				if hub := remote.GetHub(); hub != nil {
					hub.Broadcast <- "❌ Error: " + err.Error()
				}
				return StreamingMsg{
					Err:  err,
					Done: true,
				}
			}

			// Broadcast final text to remote hub
			if hub := remote.GetHub(); hub != nil {
				if res.Content != "" {
					hub.Broadcast <- res.Content
				} else {
					hub.Broadcast <- "✅ WindMist finished processing your request."
				}
			}

			duration := time.Since(startTime)

			return StreamingMsg{
				Done:     true,
				Usage:    res.Usage,
				Duration: duration,
				Turns:    res.Turns,
			}
		},
	)
}

func (m Model) getInitialMessages() []ai.Message {
	if m.store == nil || m.session == nil {
		return nil
	}

	storeMsgs, err := m.store.GetMessagesBySession(m.session.ID)
	if err != nil || len(storeMsgs) == 0 {
		return nil
	}

	// Truncate history to last 20 messages to prevent massive token usage on free tiers
	if len(storeMsgs) > 20 {
		storeMsgs = storeMsgs[len(storeMsgs)-20:]
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
