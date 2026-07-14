package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func renderConversation(m Model) string {
	var b strings.Builder

	// Always render the welcome state at the top of the conversation history
	hint := lipgloss.JoinVertical(
		lipgloss.Left,
		ui.AssistantLabelStyle.Render("⚡ WindMist is ready"),
		ui.MutedStyle.Render("Type a message below, or try:"),
		"",
		"  "+ui.LabelStyle.Render("/help")+"  "+ui.MutedLightStyle.Render("→  show all commands"),
		"  "+ui.LabelStyle.Render("/exit")+"  "+ui.MutedLightStyle.Render("→  quit"),
	)
	b.WriteString(hint)
	b.WriteString("\n\n")

	if len(m.conversation.Messages) == 0 {
		return b.String()
	}

	divider := ui.DividerStyle.Render(strings.Repeat("─", 76))
	b.WriteString(divider)
	b.WriteString("\n")

	for i, msg := range m.conversation.Messages {

		switch msg.Role {

		case "user":
			label := ui.UserLabelStyle.Render("  you")
			b.WriteString(label)
			b.WriteString("\n")
			content := ui.UserBubbleStyle.Render(msg.Content)
			b.WriteString(content)
			b.WriteString("\n")

		case "assistant":
			label := ui.AssistantLabelStyle.Render("⚡ WindMist")
			b.WriteString(label)
			b.WriteString("\n")
			contentStr := msg.Content
			if contentStr == "" && m.loading && i == len(m.conversation.Messages)-1 {
				contentStr = ui.MutedStyle.Render("Thinking...")
			} else {
				contentStr = ui.AssistantBubbleStyle.Render(contentStr)
			}
			b.WriteString(contentStr)
			b.WriteString("\n")
		}

		// subtle divider between exchanges (not after last msg)
		if i < len(m.conversation.Messages)-1 {
			b.WriteString(divider)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	return b.String()
}
