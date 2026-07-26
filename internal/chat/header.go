package chat

import (
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
	logo := ui.BaseStyle.
		Bold(true).
		Foreground(ui.Purple).
		Render("🌀 WindMist v0.5")

	// ── right: provider badge ────────────────────────────────────
	providerTag := ui.BaseStyle.
		Bold(true).
		Foreground(ui.Cyan).
		Render(m.cfg.AI.Provider)

	modelTag := ui.BaseStyle.
		Foreground(ui.MutedLight).
		Render(model)

	themeTag := ui.BaseStyle.
		Foreground(ui.Purple).
		Render(ui.CurrentThemeName)

	right := providerTag + ui.BaseStyle.Render(" ") + ui.BaseStyle.Foreground(ui.Muted).Render("›") + ui.BaseStyle.Render(" ") + modelTag + ui.BaseStyle.Render(" ") + ui.BaseStyle.Foreground(ui.Muted).Render("›") + ui.BaseStyle.Render(" ") + themeTag

	// ── padded spacer fills remaining width ──────────────────────
	const totalWidth = 78
	leftLen := lipgloss.Width(logo)
	rightLen := lipgloss.Width(right)
	// Subtract 4 for left/right borders and padding (1+1+1+1)
	gap := totalWidth - 4 - leftLen - rightLen
	if gap < 1 {
		gap = 1
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		logo,
		ui.BaseStyle.Render(strings.Repeat(" ", gap)),
		right,
	)

	box := ui.BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PurpleDark).
		Padding(0, 1).
		Width(totalWidth)

	return box.Render(row) + "\n\n"
}
