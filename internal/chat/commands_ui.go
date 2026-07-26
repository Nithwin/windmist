package chat

import (
	"github.com/Nithwin/WindMist/internal/ui"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectThemeCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		themes := ui.AvailableThemes()
		var options []selector.Option
		for _, t := range themes {
			options = append(options, selector.Option{
				Label: t,
				Value: t,
			})
		}

		return showInlineSelectorMsg{
			Title:   "Select Theme",
			Options: options,
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(opt selector.Option) tea.Cmd {
				return func() tea.Msg {
					return switchThemeSuccessMsg{
						Theme: opt.Value,
					}
				}
			},
		}
	}
}
