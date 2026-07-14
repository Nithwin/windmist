package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func renderBanner(m Model) string {
	var b strings.Builder
logo := `██╗    ██╗██╗███╗   ██╗██████╗ ███╗   ███╗██╗███████╗████████╗
██║    ██║██║████╗  ██║██╔══██╗████╗ ████║██║██╔════╝╚══██╔══╝
██║ █╗ ██║██║██╔██╗ ██║██║  ██║██╔████╔██║██║███████╗   ██║
██║███╗██║██║██║╚██╗██║██║  ██║██║╚██╔╝██║██║╚════██║   ██║
╚███╔███╔╝██║██║ ╚████║██████╔╝██║ ╚═╝ ██║██║███████║   ██║
 ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝     ╚═╝╚══════╝   ╚═╝`

	b.WriteString(ui.TitleStyle.Render(logo))
	b.WriteString("\n")

	b.WriteString(ui.SubtitleStyle.Render("⚡ WindMist - AI Coding Assistant"))
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
