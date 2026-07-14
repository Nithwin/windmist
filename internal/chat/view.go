package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(renderBanner(m))

	// Welcome
	b.WriteString(ui.SubtitleStyle.Render("Welcome to WindMist!"))
	b.WriteString("\n")
	b.WriteString("Type /help to see available commands.")
	b.WriteString("\n\n")
	
	// Input
	b.WriteString(m.input.View())
	b.WriteString("\n")

	return b.String()
}
