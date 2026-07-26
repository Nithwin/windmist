package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder

	if m.showSelector {
		return lipgloss.NewStyle().Margin(1, 2).Render(m.selectorList.View())
	}

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
		} else if m.showFilePicker {
			b.WriteString(renderFilePicker(m))
			b.WriteString("\n")
		}
	}

	if m.waitingApproval {
		approvalBox := ui.BaseStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("220")).
			Padding(1, 2).
			Render(
				ui.BaseStyle.Foreground(lipgloss.Color("220")).Bold(true).Render(fmt.Sprintf("⚠️ Agent wants to run: %s", m.approvalCommand)) + "\n\n" +
					ui.BaseStyle.Render("Allow execution? (y/N)"),
			)
		b.WriteString(approvalBox)
		b.WriteString("\n")
	} else {
		promptLabelText := " user"
		if m.inlinePrompt != "" {
			promptLabelText = " " + m.inlinePrompt
		}

		promptLabel := lipgloss.JoinHorizontal(
			lipgloss.Center,
			ui.PromptStyle.Render(promptLabelText),
			ui.BaseStyle.Foreground(ui.Muted).Render("  ›  "),
		)

		inputRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			promptLabel,
			m.input.View(),
		)

		b.WriteString(inputRow)
		b.WriteString("\n")
	}

	appStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(ui.Surface).
		Foreground(ui.White)

	return appStyle.Render(b.String())
}
