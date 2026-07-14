package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Brand purple gradient endpoints
	Purple     = lipgloss.Color("#8B5CF6")
	PurpleDark = lipgloss.Color("#6D28D9")
	PurpleDim  = lipgloss.Color("#4C1D95")
	Cyan       = lipgloss.Color("#22D3EE")
	Green      = lipgloss.Color("#10B981")
	Amber      = lipgloss.Color("#F59E0B")
	Red        = lipgloss.Color("#EF4444")
	Muted      = lipgloss.Color("#6B7280")
	MutedLight = lipgloss.Color("#9CA3AF")
	Surface    = lipgloss.Color("#1E1B2E")
	White      = lipgloss.Color("#F8FAFC")

	// ── Typography ──────────────────────────────────────────────────
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Purple)

	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Green)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	MutedLightStyle = lipgloss.NewStyle().
			Foreground(MutedLight)

	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Amber)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Green)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red)

	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B3551"))

	// ── Chat bubbles ────────────────────────────────────────────────
	UserLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Amber)

	UserBubbleStyle = lipgloss.NewStyle().
			Foreground(White).
			PaddingLeft(2)

	AssistantLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Purple)

	AssistantBubbleStyle = lipgloss.NewStyle().
				Foreground(MutedLight).
				PaddingLeft(2)

	// ── Input area ──────────────────────────────────────────────────
	InputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B3551")).
			Padding(0, 1)

	InputBoxFocusStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Purple).
				Padding(0, 1)
)