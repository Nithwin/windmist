package chat

import (
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/remote"
	"github.com/Nithwin/WindMist/internal/remote/telegram"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectRemoteCmd(m *Model) tea.Cmd {
	options := []selector.Option{
		{
			Label: "Configure Telegram Bot",
			Desc:  "Set your bot token and user ID",
			Value: "configure_telegram",
		},
	}

	if m.cfg.Remote.Telegram.BotToken != "" {
		if remote.GetHub() != nil {
			options = append([]selector.Option{
				{
					Label: "🔴 Stop Telegram Bot",
					Desc:  "Stop receiving messages on Telegram",
					Value: "stop_telegram",
				},
			}, options...)
		} else {
			options = append([]selector.Option{
				{
					Label: "🟢 Start Telegram Bot",
					Desc:  "Connect to Telegram and receive updates",
					Value: "start_telegram",
				},
			}, options...)
		}
	}

	return func() tea.Msg {
		return showInlineSelectorMsg{
			Title:   "Remote Control",
			Options: options,
			OnSelect: func(opt selector.Option) tea.Cmd {
				val := opt.FilterValue()

				switch val {
				case "configure_telegram":
					m.conversation.AddAssistant("🤖 **Telegram Bot Setup Guide:**\n\n1. Open Telegram and search for **@BotFather**.\n2. Send the `/newbot` command and follow the instructions to create your bot.\n3. Copy the **HTTP API Token** provided by BotFather.\n4. (Optional) To restrict access to yourself, message **@userinfobot** to get your numeric User ID.\n\nNow, let's configure your bot!")
					return func() tea.Msg {
						return showInlinePromptMsg{
							Prompt:     "Enter Telegram Bot Token (from @BotFather): ",
							IsPassword: true,
							OnSubmit: func(token string) tea.Cmd {
								if token == "" {
									m.conversation.AddAssistant("❌ Bot token cannot be empty.")
									return nil
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
												m.conversation.AddAssistant("❌ Failed to save config: " + err.Error())
												return nil
											}
											m.conversation.AddAssistant("✅ Telegram Bot configured successfully! You can now Start it using /remote.")
											return nil
										},
									}
								}
							},
						}
					}

				case "start_telegram":
					if remote.GetHub() == nil {
						remote.InitHub(&m.cfg.Remote)
					}

					tController, err := telegram.New(m.cfg.Remote.Telegram)
					if err != nil {
						m.conversation.AddAssistant("❌ Failed to init Telegram bot: " + err.Error())
						return nil
					}

					err = remote.GetHub().Register(tController)
					if err != nil {
						m.conversation.AddAssistant("❌ Failed to start Telegram bot: " + err.Error())
						return nil
					} else {
						m.conversation.AddAssistant("🟢 Telegram Bot started successfully!")
						return listenRemoteCmd()
					}

				case "stop_telegram":
					if remote.GetHub() != nil {
						_ = remote.GetHub().Unregister("telegram")
						m.conversation.AddAssistant("🔴 Telegram Bot stopped.")
					}
					return nil
				}
				return nil
			},
			OnCancel: func() tea.Cmd {
				return nil
			},
		}
	}
}
