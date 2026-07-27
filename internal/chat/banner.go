package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func renderBanner(m Model) string {
	var b strings.Builder

	wordmark := `██╗    ██╗██╗███╗   ██╗██████╗ ███╗   ███╗██╗███████╗████████╗
██║    ██║██║████╗  ██║██╔══██╗████╗ ████║██║██╔════╝╚══██╔══╝
██║ █╗ ██║██║██╔██╗ ██║██║  ██║██╔████╔██║██║███████╗   ██║   
██║███╗██║██║██║╚██╗██║██║  ██║██║╚██╔╝██║██║╚════██║   ██║   
╚███╔███╔╝██║██║ ╚████║██████╔╝██║ ╚═╝ ██║██║███████║   ██║   
 ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝     ╚═╝╚══════╝   ╚═╝`

	cyanStyle := ui.BaseStyle.Foreground(ui.BrandCyan)

	b.WriteString(cyanStyle.Bold(true).Render(wordmark))
	b.WriteString("\n")

	b.WriteString(ui.BaseStyle.Foreground(ui.MutedLight).Render("🌀 WindMist v2.0 — AI Coding Assistant"))
	b.WriteString("\n\n")

	b.WriteString(ui.LabelStyle.Render("Provider : "))
	b.WriteString(m.cfg.AI.Provider)
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render("Model    : "))
	if provider, err := m.cfg.ActiveProvider(); err == nil {
		b.WriteString(provider.Model)
	}
	b.WriteString("\n")

	b.WriteString(ui.LabelStyle.Render("Mode     : "))
	if m.session != nil {
		modeColor := ui.SuccessStyle
		if m.session.AgentMode == "plan" {
			modeColor = ui.BaseStyle.Foreground(ui.Amber)
		} else if m.session.AgentMode == "auto" {
			modeColor = ui.BaseStyle.Foreground(ui.Purple)
		}
		b.WriteString(modeColor.Render(strings.ToUpper(m.session.AgentMode)))
	} else {
		b.WriteString("BUILD")
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
