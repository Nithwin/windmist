package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func renderBanner(m Model) string {
	var b strings.Builder

	wordmark := `██╗    ██╗██╗███╗   ██╗██████╗ ███╗   ███╗██╗███████╗████████╗
██║    ██║██║████╗  ██║██╔══██╗████╗ ████║██║██╔════╝╚══██╔══╝
██║ █╗ ██║██║██╔██╗ ██║██║  ██║██╔████╔██║██║███████╗   ██║   
██║███╗██║██║██║╚██╗██║██║  ██║██║╚██╔╝██║██║╚════██║   ██║   
╚███╔███╔╝██║██║ ╚████║██████╔╝██║ ╚═╝ ██║██║███████║   ██║   
 ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝     ╚═╝╚══════╝   ╚═╝`

	cyanStyle := lipgloss.NewStyle().Foreground(ui.Cyan)

	b.WriteString(cyanStyle.Bold(true).Render(wordmark))
	b.WriteString("\n")

	b.WriteString(lipgloss.NewStyle().Foreground(ui.MutedLight).Render("🌀 WindMist v0.5 — AI Coding Assistant"))
	b.WriteString("\n\n")

	b.WriteString(ui.LabelStyle.Render("Provider : "))
	b.WriteString(m.cfg.AI.Provider)
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render("Model    : "))
	if provider, err := m.cfg.ActiveProvider(); err == nil {
		b.WriteString(provider.Model)
	}

	b.WriteString("\n\n")

	b.WriteString(ui.DividerStyle.Render("────────────────────────────────────────────────────────────"))
	b.WriteString("\n")

	b.WriteString(ui.SuccessStyle.Render("Type /help for commands"))
	b.WriteString("\n")

	b.WriteString(ui.SuccessStyle.Render("Type /exit to quit"))
	b.WriteString("\n")

	b.WriteString(ui.DividerStyle.Render("────────────────────────────────────────────────────────────"))
	b.WriteString("\n\n")

	return b.String()
}
