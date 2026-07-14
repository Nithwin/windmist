package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func renderConversation(m Model) string {
	var b strings.Builder

	if len(m.conversation.Messages) == 0 {
		b.WriteString(ui.SubtitleStyle.Render("⚡ WindMist"))
		b.WriteString("\n")
		b.WriteString("Welcome! Ask me anything or type ")
		b.WriteString(ui.LabelStyle.Render("/help"))
		b.WriteString(".")
		b.WriteString("\n\n")

		return b.String()
	}

	for _, msg := range m.conversation.Messages {
		switch msg.Role {

		case "user":
			b.WriteString(ui.PromptStyle.Render("❯ You"))
			b.WriteString("\n")
			b.WriteString(msg.Content)
			b.WriteString("\n\n")

		case "assistant":
			b.WriteString(ui.TitleStyle.Render("⚡ WindMist"))
			b.WriteString("\n")
			b.WriteString(msg.Content)
			b.WriteString("\n\n")
		}
	}

	return b.String()
}