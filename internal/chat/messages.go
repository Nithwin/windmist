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

// switchSubagentSuccessMsg represents a successful subagent change.
type switchSubagentSuccessMsg struct {
	Provider string
	Model    string
}

// switchModeSuccessMsg represents a successful agent mode change.
type switchModeSuccessMsg struct {
	Mode string
}

// switchSessionSuccessMsg represents a successful session change.
type switchSessionSuccessMsg struct {
	SessionID string
}

// createNewSessionMsg signals to spin up a new session.
type createNewSessionMsg struct{}

// undoFileChangeMsg signals to undo the last file edit.
type undoFileChangeMsg struct{}

// redoFileChangeMsg signals to redo the last undone file edit.
type redoFileChangeMsg struct{}

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

// ApprovalRequestMsg is sent when a tool requires user approval.
type ApprovalRequestMsg struct {
	Command      string
	ResponseChan chan bool
}
