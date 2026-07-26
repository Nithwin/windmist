package agent

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
)

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

// estimateTokens roughly estimates the number of tokens in a string.
// A common heuristic is 1 token ≈ 4 characters.
func estimateTokens(s string) int {
	return len(s) / 4
}

// estimateMessageTokens calculates the approximate token size of a message.
func estimateMessageTokens(m ai.Message) int {
	tokens := estimateTokens(m.Content)
	for _, call := range m.ToolCalls {
		tokens += estimateTokens(call.Name) + estimateTokens(fmt.Sprintf("%v", call.Args))
	}
	for _, res := range m.ToolResults {
		tokens += estimateTokens(res.Name) + estimateTokens(res.Content)
	}
	return tokens
}

// pruneMessages uses a sliding window approach based on token estimation.
// It keeps the first message (original user instruction) and dynamically
// retains as many recent messages as possible without exceeding maxTokens.
func pruneMessages(messages []ai.Message, maxTokens int) []ai.Message {
	if len(messages) <= 1 {
		return messages
	}

	// Always keep the first message
	firstMsg := messages[0]
	firstTokens := estimateMessageTokens(firstMsg)

	budget := maxTokens - firstTokens
	if budget < 0 {
		budget = 0
	}

	var keep []ai.Message
	currentTokens := 0

	// Iterate backwards from the last message to the second message
	for i := len(messages) - 1; i > 0; i-- {
		msg := messages[i]
		tokens := estimateMessageTokens(msg)

		if currentTokens+tokens > budget {
			break
		}

		// Prepend to keep slice
		keep = append([]ai.Message{msg}, keep...)
		currentTokens += tokens
	}

	// Ensure no dangling ToolResults at the start of the `keep` slice.
	// A ToolResult must be preceded by an Assistant message.
	for len(keep) > 0 && keep[0].Role == ai.RoleTool {
		keep = keep[1:]
	}

	result := make([]ai.Message, 0, len(keep)+1)
	result = append(result, firstMsg)
	result = append(result, keep...)

	return result
}
