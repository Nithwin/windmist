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
		ui.AssistantLabelStyle.Render("🌀 WindMist v0.5 is ready"),
		ui.MutedStyle.Render("Type a message below, or try:"),
		ui.BaseStyle.Render(""),
		ui.BaseStyle.Render("  "+ui.LabelStyle.Render("/help")+"  "+ui.MutedLightStyle.Render("→  show all commands")),
		ui.BaseStyle.Render("  "+ui.LabelStyle.Render("/exit")+"  "+ui.MutedLightStyle.Render("→  quit")),
	)
	b.WriteString(hint)
	b.WriteString("\n\n")

	if len(m.conversation.Messages) == 0 {
		return b.String()
	}

	maxWidth := m.MaxContentWidth()

	divider := ui.DividerStyle.Render(strings.Repeat("─", maxWidth))
	b.WriteString(divider)
	b.WriteString("\n")

	for i, msg := range m.conversation.Messages {

		switch msg.Role {

		case "user":
			label := ui.UserLabelStyle.Render("  you")
			b.WriteString(label)
			b.WriteString(ui.BaseStyle.Render("\n"))
			content := ui.UserBubbleStyle.Width(maxWidth).Render(msg.Content)
			b.WriteString(content)
			b.WriteString(ui.BaseStyle.Render("\n"))

		case "assistant":
			label := ui.AssistantLabelStyle.Render("🌀 WindMist v0.5")
			b.WriteString(label)
			b.WriteString(ui.BaseStyle.Render("\n"))
			contentStr := msg.Content
			if contentStr == "" && m.loading && i == len(m.conversation.Messages)-1 {
				contentStr = ui.MutedStyle.Render("Thinking...")
			} else {
				rendered := m.markdown.RenderWithWidth(contentStr, maxWidth)

				contentStr = ui.AssistantBubbleStyle.Render(rendered)
			}
			b.WriteString(contentStr)
			b.WriteString(ui.BaseStyle.Render("\n"))
		}

		// subtle divider between exchanges (not after last msg)
		if i < len(m.conversation.Messages)-1 {
			b.WriteString(divider)
			b.WriteString(ui.BaseStyle.Render("\n"))
		}
	}

	b.WriteString("\n")
	return b.String()
}
