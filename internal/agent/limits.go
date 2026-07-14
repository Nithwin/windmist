package agent

import "errors"

// ErrMaxTurnsExceeded is returned when the agent exhausts its turn budget.
var ErrMaxTurnsExceeded = errors.New("agent reached maximum reasoning turns")

const (
	// DefaultMaxTurns is the maximum number of reasoning iterations
	// the agent will perform before stopping.
	DefaultMaxTurns = 25

	// MaxToolCallsPerTurn is the maximum number of tool calls the
	// agent will execute from a single model response.
	MaxToolCallsPerTurn = 32
)