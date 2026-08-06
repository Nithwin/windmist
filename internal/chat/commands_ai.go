package chat

import (
	"fmt"
	"strings"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	tea "github.com/charmbracelet/bubbletea"
)

func selectProviderCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		return showInlineSelectorMsg{
			Title:   "Select AI Provider",
			Options: config.GetProviderOptions(),
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(providerOpt selector.Option) tea.Cmd {
				return func() tea.Msg {
					ollamaBaseURL := ""
					if pConfig, ok := m.cfg.Providers[providerOpt.Value]; ok {
						ollamaBaseURL = pConfig.BaseURL
					}

					return showInlineSelectorMsg{
						Title:   fmt.Sprintf("Select Model for %s", providerOpt.Value),
						Options: m.cfg.GetModelOptions(providerOpt.Value, ollamaBaseURL),
						OnCancel: func() tea.Cmd {
							return func() tea.Msg { return switchCancelMsg{} }
						},
						OnSelect: func(modelOpt selector.Option) tea.Cmd {
							return func() tea.Msg {
								if modelOpt.Value == "__CUSTOM__" {
									return showInlinePromptMsg{
										Prompt: "Enter exact model ID (e.g. gpt-4o):",
										OnSubmit: func(customVal string) tea.Cmd {
											return func() tea.Msg {
												customVal = strings.TrimSpace(customVal)
												if customVal == "" {
													return switchCancelMsg{}
												}
												m.cfg.AddCustomModel(providerOpt.Value, customVal)
												_ = config.Save(m.cfg)
												return switchProviderSuccessMsg{
													Provider: providerOpt.Value,
													Model:    customVal,
												}
											}
										},
									}
								}
								return switchProviderSuccessMsg{
									Provider: providerOpt.Value,
									Model:    modelOpt.Value,
								}
							}
						},
					}
				}
			},
		}
	}
}

func selectModelCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		ollamaBaseURL := ""
		if pConfig, ok := m.cfg.Providers[m.cfg.AI.Provider]; ok {
			ollamaBaseURL = pConfig.BaseURL
		}

		return showInlineSelectorMsg{
			Title:   fmt.Sprintf("Select Model for %s", m.cfg.AI.Provider),
			Options: m.cfg.GetModelOptions(m.cfg.AI.Provider, ollamaBaseURL),
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(modelOpt selector.Option) tea.Cmd {
				return func() tea.Msg {
					if modelOpt.Value == "__CUSTOM__" {
						return showInlinePromptMsg{
							Prompt: "Enter exact model ID (e.g. gpt-4o):",
							OnSubmit: func(customVal string) tea.Cmd {
								return func() tea.Msg {
									customVal = strings.TrimSpace(customVal)
									if customVal == "" {
										return switchCancelMsg{}
									}
									m.cfg.AddCustomModel(m.cfg.AI.Provider, customVal)
									_ = config.Save(m.cfg)
									return switchModelSuccessMsg{
										Model: customVal,
									}
								}
							},
						}
					}
					return switchModelSuccessMsg{
						Model: modelOpt.Value,
					}
				}
			},
		}
	}
}

func selectModeCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		options := []selector.Option{
			{Label: "Auto", Desc: "Routes to Chat / Plan / Build (local rules first, saves tokens)", Value: "auto"},
			{Label: "Chat", Desc: "Lightweight replies — no tools, minimal prompt (best for free tier)", Value: "chat"},
			{Label: "Build", Desc: "Full autonomy with read/write access", Value: "build"},
			{Label: "Plan", Desc: "Read-only analysis and architecture planning", Value: "plan"},
		}

		return showInlineSelectorMsg{
			Title:   "Select Agent Mode",
			Options: options,
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(opt selector.Option) tea.Cmd {
				return func() tea.Msg {
					return switchModeSuccessMsg{Mode: opt.Value}
				}
			},
		}
	}
}

func selectSubagentCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		return showInlineSelectorMsg{
			Title:   "Select Sub-Agent Provider",
			Options: append([]selector.Option{{Label: "Auto (Use Main Config)", Value: "auto"}}, config.GetProviderOptions()...),
			OnCancel: func() tea.Cmd {
				return func() tea.Msg { return switchCancelMsg{} }
			},
			OnSelect: func(providerOpt selector.Option) tea.Cmd {
				return func() tea.Msg {
					if providerOpt.Value == "auto" {
						return switchSubagentSuccessMsg{
							Provider: "",
							Model:    "",
						}
					}

					ollamaBaseURL := ""
					if pConfig, ok := m.cfg.Providers[providerOpt.Value]; ok {
						ollamaBaseURL = pConfig.BaseURL
					}

					return showInlineSelectorMsg{
						Title:   fmt.Sprintf("Select Sub-Agent Model for %s", providerOpt.Value),
						Options: append([]selector.Option{{Label: "Auto (Fast Default)", Value: "auto"}}, m.cfg.GetModelOptions(providerOpt.Value, ollamaBaseURL)...),
						OnCancel: func() tea.Cmd {
							return func() tea.Msg { return switchCancelMsg{} }
						},
						OnSelect: func(modelOpt selector.Option) tea.Cmd {
							return func() tea.Msg {
								if modelOpt.Value == "auto" {
									return switchSubagentSuccessMsg{
										Provider: providerOpt.Value,
										Model:    "",
									}
								} else if modelOpt.Value == "__CUSTOM__" {
									return showInlinePromptMsg{
										Prompt: "Enter exact model ID (e.g. gpt-4o-mini):",
										OnSubmit: func(customVal string) tea.Cmd {
											return func() tea.Msg {
												customVal = strings.TrimSpace(customVal)
												if customVal == "" {
													return switchCancelMsg{}
												}
												m.cfg.AddCustomModel(providerOpt.Value, customVal)
												_ = config.Save(m.cfg)
												return switchSubagentSuccessMsg{
													Provider: providerOpt.Value,
													Model:    customVal,
												}
											}
										},
									}
								}
								return switchSubagentSuccessMsg{
									Provider: providerOpt.Value,
									Model:    modelOpt.Value,
								}
							}
						},
					}
				}
			},
		}
	}
}

func setAPIKeyCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		provider := m.cfg.AI.Provider
		if provider == "" {
			provider = "default"
		}

		return showInlinePromptMsg{
			Prompt:     fmt.Sprintf("🔑 Enter API Key for [%s]:", provider),
			IsPassword: true,
			OnSubmit: func(val string) tea.Cmd {
				return func() tea.Msg {
					val = strings.TrimSpace(val)
					if val == "" {
						return switchCancelMsg{}
					}

					if err := m.cfg.SetAPIKey(provider, val); err != nil {
						return switchErrorMsg{Err: fmt.Errorf("failed to set api key: %w", err)}
					}

					if err := config.Save(m.cfg); err != nil {
						return switchErrorMsg{Err: fmt.Errorf("failed to save config: %w", err)}
					}

					return setAPIKeySuccessMsg{
						Provider: provider,
					}
				}
			},
		}
	}
}
