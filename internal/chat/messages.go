package chat

// ResponseMsg is sent when the AI finishes generating a response.
type ResponseMsg struct {
	Text string
	Err  error
}

// StreamingMsg represents a streamed chunk from the AI.
type StreamingMsg struct {
	Text string
	Done bool
	Err  error
}

// DoneMsg signals that streaming has completed.
type DoneMsg struct{}

// switchProviderSuccessMsg represents a successful provider change.
type switchProviderSuccessMsg struct {
	Provider string
	Model    string
}

// switchModelSuccessMsg represents a successful model change.
type switchModelSuccessMsg struct {
	Model string
}

// switchCancelMsg represents a user cancellation of the menu.
type switchCancelMsg struct{}

// switchErrorMsg represents an error running the menu.
type switchErrorMsg struct {
	Err error
}
