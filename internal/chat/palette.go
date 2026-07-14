package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func renderCommandPalette(m Model) string {
	if !m.showCommands || len(m.filteredCommands) == 0 {
		return ""
	}

	var rows []string

	title := ui.TitleStyle.Render("Commands")

	rows = append(rows, title)
	rows = append(rows, ui.DividerStyle.Render(strings.Repeat("─", 58)))

	for i, cmd := range m.filteredCommands {

		prefix := " "

		if i == m.selectedCommand {
			prefix = "▶"
		}

		row := fmt.Sprintf(
			"%s %-12s %s",
			prefix,
			ui.LabelStyle.Render(cmd.Name),
			ui.MutedStyle.Render(cmd.Description),
		)

		rows = append(rows, row)
	}

	content := strings.Join(rows, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PurpleDark).
		Padding(0, 1).
		Width(76)

	return box.Render(content)
}