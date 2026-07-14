package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func (m Model) View() string {
	if m.showSplash {
		return renderBanner(m)
	}

	var b strings.Builder

	b.WriteString(renderHeader(m))

	b.WriteString(ui.SubtitleStyle.Render("Welcome to WindMist!"))
	b.WriteString("\n")
	b.WriteString("Ask me anything or type ")
	b.WriteString(ui.LabelStyle.Render("/help"))
	b.WriteString(".")
	b.WriteString("\n\n")

	b.WriteString(ui.DividerStyle.Render("────────────────────────────────────────────────────────────"))
	b.WriteString("\n\n")

	b.WriteString(m.input.View())
	b.WriteString("\n")

	return b.String()
}