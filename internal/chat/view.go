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
		// Input row (label and textarea joined horizontally at Top so cursor is next to user ›)
		promptLabel := lipgloss.JoinHorizontal(
			lipgloss.Center,
			ui.PromptStyle.Render(" user"),
			ui.BaseStyle.Foreground(ui.Muted).Render("  ›  "),
		)

		inputRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			promptLabel,
			m.input.View(),
		)

		b.WriteString(inputRow)
		b.WriteString("\n")
		b.WriteString(renderStatusBar(m))
	}

	appStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(ui.Surface).
		Foreground(ui.White)

	return appStyle.Render(b.String())
}

func renderStatusBar(m Model) string {
	modelName := "—"
	if provider, err := m.cfg.ActiveProvider(); err == nil {
		modelName = provider.Model
	}

	tokens := 0
	cost := 0.0
	mode := "build"
	
	if m.session != nil {
		tokens = m.session.TokenCount
		cost = m.session.CostEstimate
		mode = m.session.AgentMode
	}
	
	if mode == "" {
		mode = "build"
	}

	duration := fmt.Sprintf("%.1fs", m.responseTime.Seconds())
	if m.responseTime == 0 {
		duration = "—"
	}

	modelTag := lipgloss.JoinHorizontal(lipgloss.Left, "🤖 ", modelName)
	tokenTag := lipgloss.JoinHorizontal(lipgloss.Left, "📊 ", fmt.Sprintf("%d tokens", tokens))
	costTag := lipgloss.JoinHorizontal(lipgloss.Left, "💰 ", fmt.Sprintf("$%.3f", cost))
	modeTag := lipgloss.JoinHorizontal(lipgloss.Left, "🔨 ", strings.ToUpper(mode))
	timeTag := lipgloss.JoinHorizontal(lipgloss.Left, "⏱ ", duration)

	tags := []string{modelTag, tokenTag, costTag, modeTag, timeTag}
	
	var styledTags []string
	for _, tag := range tags {
		styledTags = append(styledTags, tag)
	}

	content := strings.Join(styledTags, ui.BaseStyle.Foreground(ui.Muted).Render(" │ "))
	
	box := ui.BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Border).
		Padding(0, 1).
		Width(m.MaxContentWidth()).
		Align(lipgloss.Center).
		Foreground(ui.MutedLight)

	return box.Render(content)
}
