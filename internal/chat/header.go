package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model) string {
	model := "—"
	if provider, err := m.cfg.ActiveProvider(); err == nil {
		model = provider.Model
	}

	// ── left: brand name ──────────────────────────────────────────
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.Purple).
		Render("🌀 WindMist v0.5")

	// ── right: provider badge ────────────────────────────────────
	providerTag := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.Cyan).
		Render(m.cfg.AI.Provider)

	modelTag := lipgloss.NewStyle().
		Foreground(ui.MutedLight).
		Render(model)

	right := fmt.Sprintf("%s %s %s",
		providerTag,
		lipgloss.NewStyle().Foreground(ui.Muted).Render("›"),
		modelTag,
	)

	// ── padded spacer fills remaining width ──────────────────────
	const totalWidth = 78
	leftLen := lipgloss.Width(logo)
	rightLen := lipgloss.Width(right)
	gap := totalWidth - leftLen - rightLen
	if gap < 1 {
		gap = 1
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		logo,
		strings.Repeat(" ", gap),
		right,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PurpleDark).
		Padding(0, 1).
		Width(totalWidth)

	return box.Render(row) + "\n\n"
}