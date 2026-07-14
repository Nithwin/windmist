package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder

	if m.showSplash {
		b.WriteString(renderBanner(m))
	} else {
		b.WriteString(renderHeader(m))
		b.WriteString(m.viewport.View())

		b.WriteString(m.viewport.View())

		// Separator above input area
		b.WriteString(ui.DividerStyle.Render(strings.Repeat("─", 80)))
		b.WriteString("\n\n")

		// Show command palette ABOVE the input
		if m.showCommands {
			b.WriteString(renderCommandPalette(m))
			b.WriteString("\n")
		}
	}

	// Input row
	prompt := ui.PromptStyle.Render(" user")

	inputLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		prompt,
		lipgloss.NewStyle().
			Foreground(ui.Muted).
			Render("  ›  "),
		m.input.View(),
	)

	b.WriteString(inputLine)
	b.WriteString("\n")

	return b.String()
}
