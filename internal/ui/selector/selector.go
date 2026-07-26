package selector

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Option represents a selectable item in the selector list.
type Option struct {
	Label string
	Desc  string
	Value string
}

func (o Option) Title() string       { return o.Label }
func (o Option) Description() string { return o.Desc }
func (o Option) FilterValue() string { return o.Label + " " + o.Value }

// ErrCancelled is returned when the user cancels the selector (e.g. via Esc or Ctrl+C).
var ErrCancelled = fmt.Errorf("selection cancelled")

type model struct {
	list      list.Model
	selected  *Option
	cancelled bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(Option); ok {
				m.selected = &i
				return m, tea.Quit
			}
		case "esc":
			if !m.list.SettingFilter() {
				m.cancelled = true
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

// Run displays an interactive arrow-key list and returns the selected Option.
func Run(title, description string, options []Option) (Option, error) {
	if len(options) == 0 {
		return Option{}, fmt.Errorf("no options provided")
	}

	items := make([]list.Item, len(options))
	for i, opt := range options {
		items[i] = opt
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(ui.Cyan).BorderForeground(ui.Cyan)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(ui.Cyan).BorderForeground(ui.Cyan)

	l := list.New(items, d, 80, 20)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Background(ui.Purple).Foreground(ui.White).Padding(0, 1)

	p := tea.NewProgram(model{list: l}, tea.WithAltScreen())

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

// RunWithDefault displays an interactive arrow-key list and pre-selects the defaultValue.
func RunWithDefault(title, description string, options []Option, defaultValue string) (Option, error) {
	if len(options) == 0 {
		return Option{}, fmt.Errorf("no options provided")
	}

	items := make([]list.Item, len(options))
	selectedIndex := 0
	for i, opt := range options {
		items[i] = opt
		if opt.Value == defaultValue {
			selectedIndex = i
		}
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(ui.Cyan).BorderForeground(ui.Cyan)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(ui.Cyan).BorderForeground(ui.Cyan)

	l := list.New(items, d, 80, 20)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Background(ui.Purple).Foreground(ui.White).Padding(0, 1)
	l.Select(selectedIndex)

	p := tea.NewProgram(model{list: l}, tea.WithAltScreen())

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
