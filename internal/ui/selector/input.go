package selector

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputModel struct {
	title       string
	placeholder string
	textInput   textinput.Model
	cancelled   bool
	submitted   bool
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if strings.TrimSpace(m.textInput.Value()) != "" {
				m.submitted = true
				return m, tea.Quit
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.Purple).MarginBottom(1)
	b.WriteString(titleStyle.Render(m.title) + "\n\n")

	b.WriteString(m.textInput.View() + "\n\n")

	footerStyle := lipgloss.NewStyle().Foreground(ui.Muted)
	b.WriteString(footerStyle.Render("enter submit • esc cancel") + "\n")

	return b.String()
}

// RunInput prompts the user to type a custom string value.
func RunInput(title, placeholder, initialValue string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	if initialValue != "" {
		ti.SetValue(initialValue)
	}

	p := tea.NewProgram(inputModel{
		title:       title,
		placeholder: placeholder,
		textInput:   ti,
	})

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running input: %w", err)
	}

	m, ok := finalModel.(inputModel)
	if !ok || m.cancelled || !m.submitted {
		return "", ErrCancelled
	}

	return strings.TrimSpace(m.textInput.Value()), nil
}
