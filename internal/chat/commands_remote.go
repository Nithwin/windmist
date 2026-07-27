package chat

import (
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/remote"
	"github.com/Nithwin/WindMist/internal/remote/telegram"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectRemoteCmd(m *Model) tea.Cmd {
	var options []selector.Option

	if m.cfg.Remote.Telegram.BotToken == "" {
		options = append(options, selector.Option{
			Label: "🤖 Set up Telegram Integration",
			Desc:  "Connect a new Telegram bot to chat with WindMist remotely",
			Value: "configure_telegram",
		})
	} else {
		if remote.GetHub() != nil && remote.GetHub().HasController("telegram") {
			options = append(options, selector.Option{
				Label: "🔴 Stop Telegram Bot",
				Desc:  "Stop receiving messages on Telegram",
				Value: "stop_telegram",
			})
		} else {
			options = append(options, selector.Option{
				Label: "🟢 Start Telegram Bot",
				Desc:  "Connect to Telegram and receive updates",
				Value: "start_telegram",
			})
		}

		options = append(options, selector.Option{
			Label: "⚙️ Reconfigure Telegram Bot",
			Desc:  "Update your bot token and user ID",
			Value: "configure_telegram",
		})

		options = append(options, selector.Option{
			Label: "🗑️ Remove Telegram Configuration",
			Desc:  "Delete your bot token and stop the integration",
			Value: "delete_telegram",
		})
	}

	return func() tea.Msg {
		return showInlineSelectorMsg{
			Title:   "🌐 Remote Control & Integrations",
			Options: options,
			OnSelect: func(opt selector.Option) tea.Cmd {
				val := opt.Value

				switch val {
				case "configure_telegram":
					guideMsg := "🤖 **Telegram Bot Setup Guide:**\n\n1. Open Telegram and search for **@BotFather**.\n2. Send the `/newbot` command and follow the instructions to create your bot.\n3. Copy the **HTTP API Token** provided by BotFather.\n4. (Optional) To restrict access to yourself, message **@userinfobot** to get your numeric User ID.\n\nNow, let's configure your bot!"
					
					return tea.Batch(
						func() tea.Msg {
							return ResponseMsg{Text: guideMsg}
						},
						func() tea.Msg {
							return showInlinePromptMsg{
								Prompt:     "Enter Telegram Bot Token (from @BotFather): ",
								IsPassword: true,
								OnSubmit: func(token string) tea.Cmd {
									if token == "" {
										return func() tea.Msg { return ResponseMsg{Text: "❌ Bot token cannot be empty."} }
									}
									m.cfg.Remote.Telegram.BotToken = token

									return func() tea.Msg {
										return showInlinePromptMsg{
											Prompt:     "Enter your Telegram User ID (numeric, optional): ",
											IsPassword: false,
											OnSubmit: func(userID string) tea.Cmd {
												m.cfg.Remote.Telegram.AllowedID = userID
												m.cfg.Remote.Telegram.Enabled = true
												if err := config.Save(m.cfg); err != nil {
													return func() tea.Msg { return ResponseMsg{Text: "❌ Failed to save config: " + err.Error()} }
												}
												return func() tea.Msg { return ResponseMsg{Text: "✅ Telegram Bot configured successfully! You can now Start it using /remote."} }
											},
										}
									}
								},
							}
						},
					)

				case "start_telegram":
					if remote.GetHub() == nil {
						remote.InitHub(&m.cfg.Remote)
					}

					tController, err := telegram.New(m.cfg.Remote.Telegram)
					if err != nil {
						return func() tea.Msg { return ResponseMsg{Text: "❌ Failed to init Telegram bot: " + err.Error()} }
					}

					err = remote.GetHub().Register(tController)
					if err != nil {
						return func() tea.Msg { return ResponseMsg{Text: "❌ Failed to start Telegram bot: " + err.Error()} }
					} else {
						return tea.Batch(
							func() tea.Msg { return ResponseMsg{Text: "🟢 Telegram Bot started successfully!"} },
							listenRemoteCmd(),
						)
					}

				case "stop_telegram":
					if remote.GetHub() != nil {
						_ = remote.GetHub().Unregister("telegram")
						return func() tea.Msg { return ResponseMsg{Text: "🔴 Telegram Bot stopped."} }
					}
					return nil

				case "delete_telegram":
					if remote.GetHub() != nil {
						_ = remote.GetHub().Unregister("telegram")
					}
					m.cfg.Remote.Telegram.BotToken = ""
					m.cfg.Remote.Telegram.AllowedID = ""
					m.cfg.Remote.Telegram.Enabled = false
					if err := config.Save(m.cfg); err != nil {
						return func() tea.Msg { return ResponseMsg{Text: "❌ Failed to save config: " + err.Error()} }
					}
					return func() tea.Msg { return ResponseMsg{Text: "🗑️ Telegram configuration has been deleted."} }
				}
				return nil
			},
			OnCancel: func() tea.Cmd {
				return nil
			},
		}
	}
}
