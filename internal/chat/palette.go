package chat

import (
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
)

func renderCommandPalette(m Model) string {
	if !m.showCommands || len(m.filteredCommands) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")

	b.WriteString(ui.DividerStyle.Render("┌──────────────────────────────────────────────────────────────┐"))
	b.WriteString("\n")

	for i, cmd := range m.filteredCommands {

		prefix := "  "

		if i == m.selectedCommand {
			prefix = "▶ "
		}

		b.WriteString(prefix)
		b.WriteString(ui.LabelStyle.Render(cmd.Name))
		b.WriteString("    ")
		b.WriteString(ui.MutedStyle.Render(cmd.Description))
		b.WriteString("\n")
	}

	b.WriteString(ui.DividerStyle.Render("└──────────────────────────────────────────────────────────────┘"))
	b.WriteString("\n")

	return b.String()
}