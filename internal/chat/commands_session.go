package chat

import (
	"fmt"
	"os"

	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectSessionCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if program == nil {
			return switchErrorMsg{Err: fmt.Errorf("program instance not initialized")}
		}
		if m.store == nil {
			return switchErrorMsg{Err: fmt.Errorf("database not initialized")}
		}

		cwd, _ := os.Getwd()
		sessions, err := m.store.ListSessionsByProject(cwd)
		if err != nil {
			return switchErrorMsg{Err: fmt.Errorf("failed to fetch sessions: %w", err)}
		}

		if len(sessions) == 0 {
			return switchErrorMsg{Err: fmt.Errorf("no past sessions found in this project")}
		}

		var options []selector.Option
		for _, s := range sessions {
			desc := fmt.Sprintf("%s | Tokens: %d | Cost: $%.3f", s.UpdatedAt.Format("Jan 02 15:04"), s.TokenCount, s.CostEstimate)
			options = append(options, selector.Option{
				Label: s.Title,
				Desc:  desc,
				Value: s.ID,
			})
		}

		return showInlineSelectorMsg{
			Title:   "Select Session",
			Options: options,
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(opt selector.Option) tea.Cmd {
				return func() tea.Msg {
					return switchSessionSuccessMsg{
						SessionID: opt.Value,
					}
				}
			},
		}
	}
}
