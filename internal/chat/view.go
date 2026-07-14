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
		b.WriteString("\n")

		// Separator above input area
		b.WriteString(ui.DividerStyle.Render(strings.Repeat("─", 80)))
		b.WriteString("\n\n")

		// Show command palette ABOVE the input
		if m.showCommands {
			b.WriteString(renderCommandPalette(m))
			b.WriteString("\n")
		}
	}

	// Input row (label and textarea joined horizontally at Top so cursor is next to user ›)
	promptLabel := lipgloss.JoinHorizontal(
		lipgloss.Center,
		ui.PromptStyle.Render(" user"),
		lipgloss.NewStyle().Foreground(ui.Muted).Render("  ›  "),
	)

	inputRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		promptLabel,
		m.input.View(),
	)

	b.WriteString(inputRow)
	b.WriteString("\n")

	return b.String()
}
