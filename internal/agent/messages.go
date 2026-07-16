package agent

import "github.com/Nithwin/WindMist/internal/ai"

// appendUser appends a user message to the conversation history.
func appendUser(messages []ai.Message, content string) []ai.Message {
	return append(messages, ai.Message{
		Role:    ai.RoleUser,
		Content: content,
	})
}

// appendAssistant appends an assistant message, including any tool calls requested by the model.
func appendAssistant(messages []ai.Message, content string, toolCalls []ai.ToolCall) []ai.Message {
	return append(messages, ai.Message{
		Role:      ai.RoleAssistant,
		Content:   content,
		ToolCalls: toolCalls,
	})
}

// appendToolResults appends tool execution results to the conversation history.
func appendToolResults(messages []ai.Message, results []ai.ToolResult) []ai.Message {
	if len(results) == 0 {
		return messages
	}
	return append(messages, ai.Message{
		Role:        ai.RoleTool,
		ToolResults: results,
	})
}

// pruneMessages shortens the conversation history to prevent exceeding model context limits.
// It always preserves the first message (original user prompt) and keeps the most recent
// maxKeep messages. maxKeep should be an even number (e.g., 8) to ensure that Assistant
// tool calls and Tool results are never separated.
func pruneMessages(messages []ai.Message, maxKeep int) []ai.Message {
	// If the history is already small enough, do nothing
	if len(messages) <= maxKeep+1 {
		return messages
	}

	pruned := make([]ai.Message, 0, maxKeep+1)

	// 1. Always keep the first message (index 0)
	pruned = append(pruned, messages[0])

	// 2. Keep the last maxKeep messages from the end of the slice
	startIdx := len(messages) - maxKeep
	pruned = append(pruned, messages[startIdx:]...)

	return pruned
}
