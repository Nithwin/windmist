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
	logo := ui.BaseStyle.
		Bold(true).
		Foreground(ui.BrandCyan).
		Render("🌀 WindMist v0.5")

	// ── right: status tags ────────────────────────────────────
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

	modelTag := ui.BaseStyle.Foreground(ui.Cyan).Bold(true).Render(model)
	tokenTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(fmt.Sprintf("📊 %d", tokens))
	costTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(fmt.Sprintf("💰 $%.3f", cost))
	modeTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(fmt.Sprintf("🔨 %s", strings.ToUpper(mode)))
	timeTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(fmt.Sprintf("⏱ %s", duration))
	themeTag := ui.BaseStyle.Foreground(ui.BrandCyan).Render(ui.CurrentThemeName)

	tags := []string{modelTag, tokenTag, costTag, modeTag, timeTag, themeTag}
	
	right := strings.Join(tags, ui.BaseStyle.Foreground(ui.Muted).Render(" │ "))

	// ── padded spacer fills remaining width ──────────────────────
	totalWidth := m.MaxContentWidth()
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
