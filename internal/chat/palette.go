package chat

import (
	"fmt"
	"path/filepath"
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

	box := ui.BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PurpleDark).
		Padding(0, 1).
		Width(76)

	return box.Render(content)
}

func renderFilePicker(m Model) string {
	if !m.showFilePicker || len(m.filteredFiles) == 0 {
		return ""
	}

	var rows []string

	title := ui.TitleStyle.Render("Attach File")
	rows = append(rows, title)
	rows = append(rows, ui.DividerStyle.Render(strings.Repeat("─", 58)))

	for i, file := range m.filteredFiles {
		prefix := " "
		if i == m.selectedFile {
			prefix = "▶"
		}

		// Highlight filename vs path
		dir := filepath.Dir(file)
		name := filepath.Base(file)

		displayPath := ""
		if dir != "." {
			displayPath = dir + "/"
		}

		row := fmt.Sprintf(
			"%s %s%s",
			prefix,
			ui.MutedStyle.Render(displayPath),
			ui.LabelStyle.Render(name),
		)
		rows = append(rows, row)
	}

	content := strings.Join(rows, "\n")
	box := ui.BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Cyan).
		Padding(0, 1).
		Width(76)

	return box.Render(content)
}
