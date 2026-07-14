package chat

// refreshViewport updates the viewport with the latest conversation.
func (m *Model) refreshViewport() {
	m.viewport.SetContent(renderConversation(*m))
	m.viewport.GotoBottom()
}