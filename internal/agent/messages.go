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
