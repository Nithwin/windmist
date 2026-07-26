package agent

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/pkoukk/tiktoken-go"
)

var tokenizer *tiktoken.Tiktoken

func init() {
	var err error
	tokenizer, err = tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		fmt.Printf("Warning: failed to load tiktoken: %v\n", err)
	}
}

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

// countTokens accurately counts the number of tokens in a string using tiktoken.
func countTokens(s string) int {
	if tokenizer != nil {
		return len(tokenizer.Encode(s, nil, nil))
	}
	// Fallback heuristic if tokenizer failed to load
	return len(s) / 4
}

// countMessageTokens calculates the token size of a message.
func countMessageTokens(m ai.Message) int {
	tokens := countTokens(m.Content)
	for _, call := range m.ToolCalls {
		tokens += countTokens(call.Name) + countTokens(fmt.Sprintf("%v", call.Args))
	}
	for _, res := range m.ToolResults {
		tokens += countTokens(res.Name) + countTokens(res.Content)
	}
	// Add base padding per message
	return tokens + 4
}

// MemoryStrategy defines an interface for context window management.
type MemoryStrategy interface {
	Prune(messages []ai.Message, maxTokens int) []ai.Message
}

// SlidingWindowMemory implements a basic sliding window token pruner.
type SlidingWindowMemory struct{}

// Prune uses a sliding window approach based on exact token estimation.
// It keeps the first message (original user instruction) and dynamically
// retains as many recent messages as possible without exceeding maxTokens.
func (s SlidingWindowMemory) Prune(messages []ai.Message, maxTokens int) []ai.Message {
	if len(messages) <= 1 {
		return messages
	}

	// Always keep the first message
	firstMsg := messages[0]
	firstTokens := countMessageTokens(firstMsg)

	budget := maxTokens - firstTokens
	if budget < 0 {
		budget = 0
	}

	var keep []ai.Message
	currentTokens := 0

	// Iterate backwards from the last message to the second message
	for i := len(messages) - 1; i > 0; i-- {
		msg := messages[i]
		tokens := countMessageTokens(msg)

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
