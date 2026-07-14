package chat

// updateViewportSize recalculates viewport width and height dynamically based on window size and open panels.
func (m *Model) updateViewportSize() {
	if m.height <= 0 || m.width <= 0 {
		return
	}

	m.viewport.Width = m.width

	// Fixed lines surrounding viewport when showSplash == false:
	// - renderHeader(m): 5 lines (box + 2 newlines)
	// - viewport trailing newline: 1 line
	// - divider above input: 3 lines (line + 2 newlines)
	// - inputLine + trailing newline: 2 lines
	// Total fixed lines = 11
	fixedLines := 11

	if m.showCommands && len(m.filteredCommands) > 0 {
		// command palette box (len + 4) + trailing newline (1) = len + 5 lines
		fixedLines += 5 + len(m.filteredCommands)
	}

	availableHeight := m.height - fixedLines
	if availableHeight < 3 {
		availableHeight = 3
	}

	m.viewport.Height = availableHeight
}

// refreshViewport updates the viewport with the latest conversation and ensures correct sizing.
func (m *Model) refreshViewport() {
	m.updateViewportSize()
	m.viewport.SetContent(renderConversation(*m))
	m.viewport.GotoBottom()
}