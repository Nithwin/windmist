package selector

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Option represents a selectable item in the selector list.
type Option struct {
	Label       string
	Description string
	Value       string
}

// ErrCancelled is returned when the user cancels the selector (e.g. via Esc or Ctrl+C).
var ErrCancelled = fmt.Errorf("selection cancelled")

type model struct {
	title       string
	description string
	options     []Option
	cursor      int
	selected    *Option
	cancelled   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.options) - 1
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case "enter":
			if len(m.options) > 0 {
				m.selected = &m.options[m.cursor]
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.Purple).
		MarginBottom(1)
	b.WriteString(titleStyle.Render(m.title) + "\n")

	// Optional Description
	if m.description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(ui.MutedLight).
			MarginBottom(1)
		b.WriteString(descStyle.Render(m.description) + "\n\n")
	} else {
		b.WriteString("\n")
	}

	// Options
	for i, opt := range m.options {
		cursor := "  "
		if m.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(ui.Cyan).Bold(true).Render("❯ ")
		}

		labelStyle := lipgloss.NewStyle().Foreground(ui.White)
		if m.cursor == i {
			labelStyle = lipgloss.NewStyle().Foreground(ui.Cyan).Bold(true)
		}

		label := labelStyle.Render(opt.Label)

		var desc string
		if opt.Description != "" {
			descStyle := lipgloss.NewStyle().Foreground(ui.Muted)
			if m.cursor == i {
				descStyle = lipgloss.NewStyle().Foreground(ui.MutedLight)
			}
			desc = "  " + descStyle.Render(opt.Description)
		}

		b.WriteString(fmt.Sprintf("%s%s%s\n", cursor, label, desc))
	}

	// Footer instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(ui.Muted).
		MarginTop(1)
	b.WriteString("\n" + footerStyle.Render("↑/↓ navigate • enter select • esc/q cancel") + "\n")

	return b.String()
}

// Run displays an interactive arrow-key list and returns the selected Option.
func Run(title, description string, options []Option) (Option, error) {
	if len(options) == 0 {
		return Option{}, fmt.Errorf("no options provided")
	}

	p := tea.NewProgram(model{
		title:       title,
		description: description,
		options:     options,
	})

	finalModel, err := p.Run()
	if err != nil {
		return Option{}, fmt.Errorf("error running selector: %w", err)
	}

	m, ok := finalModel.(model)
	if !ok || m.cancelled || m.selected == nil {
		return Option{}, ErrCancelled
	}

	return *m.selected, nil
}
