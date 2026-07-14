package chat

// ResponseMsg is sent when the AI finishes generating a response.
type ResponseMsg struct {
	Text string
	Err  error
}

// StreamingMsg represents a streamed chunk from the AI.
type StreamingMsg struct {
	Text string
}

// DoneMsg signals that streaming has completed.
type DoneMsg struct{}