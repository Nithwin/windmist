package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// formatTokenCount renders a human-friendly token count (e.g. "1.2k")
func formatTokenCount(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func renderHeader(m Model) string {
	model := "—"
	if provider, err := m.cfg.ActiveProvider(); err == nil {
		model = provider.Model
	}

	// ── left: brand name & path ────────────────────────────────────
	logo := ui.BaseStyle.
		Bold(true).
		Foreground(ui.BrandCyan).
		Render("🌀 WindMist v2.0")

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	
	cwdDisplay := filepath.Base(cwd)
	dirTag := ui.BaseStyle.Foreground(ui.Muted).Render(fmt.Sprintf(" [%s]", cwdDisplay))
	
	left := lipgloss.JoinHorizontal(lipgloss.Center, logo, dirTag)

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

	// Token display: show real-time streaming tokens when active,
	// otherwise show session total
	var tokenTag string
	if m.streaming && m.streamTokens.TotalTokens > 0 {
		// Show live streaming tokens with animated indicator
		frame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
		inTok := formatTokenCount(m.streamTokens.InputTokens)
		outTok := formatTokenCount(m.streamTokens.OutputTokens)
		tokenTag = ui.BaseStyle.Foreground(ui.Cyan).Bold(true).Render(
			fmt.Sprintf("%s %s↑ %s↓", frame, inTok, outTok),
		)
	} else if m.loading {
		// Loading but no token data yet
		frame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
		tokenTag = ui.BaseStyle.Foreground(ui.Cyan).Render(
			fmt.Sprintf("%s %s tok", frame, formatTokenCount(tokens)),
		)
	} else {
		tokenTag = ui.BaseStyle.Foreground(ui.MutedLight).Render(
			fmt.Sprintf("%s tok", formatTokenCount(tokens)),
		)
	}

	// Only show cost if it's > 0 (to avoid showing $0.000 for free APIs like Ollama/Groq)
	costStr := ""
	if cost > 0 {
		costStr = fmt.Sprintf("$%.3f", cost)
	}
	costTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(costStr)

	modeTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(strings.ToUpper(mode))
	timeTag := ui.BaseStyle.Foreground(ui.MutedLight).Render(duration)
	themeTag := ui.BaseStyle.Foreground(ui.BrandCyan).Render(ui.CurrentThemeName)

	tags := []string{modelTag, tokenTag}
	if costStr != "" {
		tags = append(tags, costTag)
	}
	tags = append(tags, modeTag, timeTag, themeTag)

	// Show queued message indicator
	if m.queuedMessage != "" {
		queueTag := ui.BaseStyle.Foreground(lipgloss.Color("220")).Bold(true).Render("📋 QUEUED")
		tags = append(tags, queueTag)
	}

	right := strings.Join(tags, ui.BaseStyle.Foreground(ui.Muted).Render(" │ "))

	// ── padded spacer fills remaining width ──────────────────────
	totalWidth := m.MaxContentWidth()
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)

	// Subtract 4 for left/right borders and padding (1+1+1+1)
	gap := totalWidth - 4 - leftLen - rightLen
	if gap < 1 {
		gap = 1
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		left,
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
