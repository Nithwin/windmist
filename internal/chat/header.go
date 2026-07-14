package chat

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model) string {
	model := ""

	if provider, err := m.cfg.ActiveProvider(); err == nil {
		model = provider.Model
	}

	left := ui.TitleStyle.Render("⚡ WindMist")
	right := fmt.Sprintf("%s • %s", m.cfg.AI.Provider, model)

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(8).Render(""),
		ui.SuccessStyle.Render(right),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(0, 1).
		Width(80)

	return box.Render(header) + "\n\n"
}