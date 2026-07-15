package chat

// ChatMessage represents a single message in the conversation.
type ChatMessage struct {
	Role    string
	Content string
}

// Conversation stores the current chat history.
type Conversation struct {
	Messages []ChatMessage
}

// AddUser adds a user message.
func (c *Conversation) AddUser(content string) {
	c.Messages = append(c.Messages, ChatMessage{
		Role:    "user",
		Content: content,
	})
}

// AddAssistant adds an assistant message.
func (c *Conversation) AddAssistant(content string) {
	c.Messages = append(c.Messages, ChatMessage{
		Role:    "assistant",
		Content: content,
	})
}

// Clear starts a new conversation.
func (c *Conversation) Clear() {
	c.Messages = nil
}
