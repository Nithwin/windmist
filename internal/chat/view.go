package chat

import (
	"fmt"
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
		b.WriteString(ui.DividerStyle.Render(strings.Repeat("─", m.MaxContentWidth()+4)))
		b.WriteString("\n\n")

		// Show command palette ABOVE the input
		if m.showCommands {
			b.WriteString(renderCommandPalette(m))
			b.WriteString("\n")
		}
	}

	if m.waitingApproval {
		approvalBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("220")).
			Padding(1, 2).
			Render(
				lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true).Render(fmt.Sprintf("⚠️ Agent wants to run: %s", m.approvalCommand)) + "\n\n" +
					lipgloss.NewStyle().Render("Allow execution? (y/N)"),
			)
		b.WriteString(approvalBox)
		b.WriteString("\n")
	} else {
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
	}

	return b.String()
}
