package chat

import (
	"context"
	"strings"

	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Scroll conversation when command palette or file picker is closed.
	if !m.showCommands && !m.showFilePicker {
		switch msg.String() {

		case "ctrl+up", "shift+up":
			m.viewport.ScrollUp(1)
			return m, nil

		case "ctrl+down", "shift+down":
			m.viewport.ScrollDown(1)
			return m, nil

		case "up":
			if len(m.inputHistory) > 0 && m.historyIndex > 0 {
				m.historyIndex--
				m.input.SetValue(m.inputHistory[m.historyIndex])
				m.input.CursorEnd()
			}
			return m, nil

		case "down":
			if len(m.inputHistory) > 0 && m.historyIndex < len(m.inputHistory) {
				m.historyIndex++
				if m.historyIndex == len(m.inputHistory) {
					m.input.SetValue("")
				} else {
					m.input.SetValue(m.inputHistory[m.historyIndex])
					m.input.CursorEnd()
				}
			}
			return m, nil

		case "pgup":
			m.viewport.ScrollUp(m.viewport.Height / 2)
			return m, nil

		case "pgdown":
			m.viewport.ScrollDown(m.viewport.Height / 2)
			return m, nil

		case "home":
			m.viewport.GotoTop()
			return m, nil

		case "end":
			m.viewport.GotoBottom()
			return m, nil
		}
	}

	// Handle approval keys
	if m.waitingApproval {
		switch msg.String() {
		case "y", "Y":
			if m.approvalChan != nil {
				m.approvalChan <- true
			}
			m.waitingApproval = false
			m.refreshViewport()
			return m, nil
		case "n", "N":
			if m.approvalChan != nil {
				m.approvalChan <- false
			}
			m.waitingApproval = false
			m.refreshViewport()
			return m, nil
		case "ctrl+c", "esc":
			if m.approvalChan != nil {
				m.approvalChan <- false
			}
			m.waitingApproval = false
			if m.agent != nil {
				m.agent.Close()
			}
			return m, tea.Quit
		}
		// Block other inputs
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "esc":
		if m.loading && m.cancel != nil {
			m.cancel()
			m.loading = false
			m.streaming = false
			m.spinnerFrame = 0
			m.conversation.AddAssistant("\n\n*(Cancelled by user)*")
			m.refreshViewport()
			return m, nil
		}
		if m.agent != nil {
			m.agent.Close()
		}
		return m, tea.Quit
	}

	// Hide splash on first key press.
	if m.showSplash {
		m.showSplash = false

		// Preserve the first typed character.
		if len(msg.String()) == 1 {
			m.input.SetValue(msg.String())
			m.input.CursorEnd()
		}

		m.refreshViewport()
		return m, nil
	}

	// Navigate the command palette.
	if m.showCommands {
		switch msg.String() {

		case "up":
			if m.selectedCommand > 0 {
				m.selectedCommand--
			}
			return m, nil

		case "down":
			if m.selectedCommand < len(m.filteredCommands)-1 {
				m.selectedCommand++
			}
			return m, nil

		case "esc":
			m.showCommands = false
			m.filteredCommands = nil
			m.selectedCommand = 0
			return m, nil
		}
	}

	// Navigate the file picker.
	if m.showFilePicker {
		switch msg.String() {

		case "up":
			if m.selectedFile > 0 {
				m.selectedFile--
			}
			return m, nil

		case "down":
			if m.selectedFile < len(m.filteredFiles)-1 {
				m.selectedFile++
			}
			return m, nil

		case "esc":
			m.showFilePicker = false
			m.filteredFiles = nil
			m.selectedFile = 0
			return m, nil
		}
	}

	// Handle Inline Prompt escape
	if m.onPromptSubmit != nil && msg.String() == "esc" {
		m.onPromptSubmit = nil
		m.inlinePrompt = ""
		m.isPassword = false
		m.input.SetValue("")
		return m, nil
	}

	// Handle Inline Selector keys
	if m.showSelector {
		switch msg.String() {
		case "esc", "ctrl+c":
			m.showSelector = false
			if m.onCancel != nil {
				return m, m.onCancel()
			}
			return m, nil
		case "enter":
			m.showSelector = false
			if i, ok := m.selectorList.SelectedItem().(selector.Option); ok {
				if m.onSelect != nil {
					return m, m.onSelect(i)
				}
			}
			return m, nil
		}
	}

	switch msg.String() {

	case "enter":
		prompt := strings.TrimSpace(m.input.Value())

		// Execute inline prompt submit
		if m.onPromptSubmit != nil {
			m.input.SetValue("")
			cmd := m.onPromptSubmit(prompt)
			m.onPromptSubmit = nil
			m.inlinePrompt = ""
			m.isPassword = false
			return m, cmd
		}

		if prompt == "" {
			return m, nil
		}

		// Execute selected command from palette.
		if m.showCommands && len(m.filteredCommands) > 0 {
			cmd := m.filteredCommands[m.selectedCommand]

			m.showCommands = false
			m.filteredCommands = nil
			m.selectedCommand = 0
			m.input.SetValue("")

			return m, cmd.Execute(&m)
		}

		// Execute selected file from picker.
		if m.showFilePicker && len(m.filteredFiles) > 0 {
			file := m.filteredFiles[m.selectedFile]

			m.showFilePicker = false
			m.filteredFiles = nil
			m.selectedFile = 0

			// Replace the @query with the filename
			value := m.input.Value()
			words := strings.Fields(value)
			if len(words) > 0 {
				words[len(words)-1] = file + " "
				newValue := strings.Join(words, " ")
				m.input.SetValue(newValue)
				m.input.CursorEnd()
			}

			// Don't send the message yet
			return m, nil
		}

		// Execute typed slash command.
		if strings.HasPrefix(prompt, "/") {
			m.inputHistory = append(m.inputHistory, prompt)
			m.historyIndex = len(m.inputHistory)

			if command, ok := FindCommand(prompt); ok {
				m.input.SetValue("")
				return m, command.Execute(&m)
			}

			m.conversation.AddAssistant("Unknown command: " + prompt)
			m.input.SetValue("")
			return m, nil
		}

		// Normal AI message.
		// Block new AI messages while a request is in-flight.
		if m.loading {
			// Don't fire another request — just show feedback.
			m.conversation.AddAssistant("⏳ *Please wait — a request is still in progress. Press `Ctrl+C` to cancel it.*")
			m.refreshViewport()
			return m, nil
		}

		m.inputHistory = append(m.inputHistory, prompt)
		m.historyIndex = len(m.inputHistory)

		m.conversation.AddUser(prompt)
		m.refreshViewport()
		m.loading = true
		m.streaming = true

		m.input.SetValue("")

		// Create an empty assistant message.
		// Streaming chunks will be appended to this.
		m.conversation.AddAssistant("")
		m.refreshViewport()

		// Cancel any previous in-flight request before starting a new one.
		if m.cancel != nil {
			m.cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		m.cancel = cancel

		return m, m.sendMessageCmd(ctx, prompt)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Update slash command suggestions (check first line only).
	value := m.input.Value()
	firstLine := strings.SplitN(value, "\n", 2)[0]

	if strings.HasPrefix(firstLine, "/") {
		m.showCommands = true
		m.filteredCommands = FilterCommands(firstLine)
		m.showFilePicker = false
	} else {
		m.showCommands = false
		m.filteredCommands = nil
		m.selectedCommand = 0

		// Check for file picker trigger (@) anywhere in the text
		// Don't trigger if there's a trailing space
		if strings.HasSuffix(value, " ") || strings.HasSuffix(value, "\n") {
			m.showFilePicker = false
			m.filteredFiles = nil
			m.selectedFile = 0
		} else {
			words := strings.Fields(value)
			if len(words) > 0 && strings.HasPrefix(words[len(words)-1], "@") {
				m.showFilePicker = true
				query := words[len(words)-1][1:]
				m.filteredFiles = FilterFiles(m.workspaceFiles, query)
				if m.selectedFile >= len(m.filteredFiles) {
					m.selectedFile = 0
				}
			} else {
				m.showFilePicker = false
				m.filteredFiles = nil
				m.selectedFile = 0
			}
		}
	}
	m.updateViewportSize()

	return m, cmd
}
