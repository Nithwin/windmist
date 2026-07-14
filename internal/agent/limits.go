package agent

const (
	// DefaultMaxTurns is the maximum number of reasoning iterations
	// the agent will perform before stopping.
	DefaultMaxTurns = 25

	// MaxToolCallsPerTurn is the maximum number of tool calls the
	// agent will execute from a single model response.
	MaxToolCallsPerTurn = 32
)