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
	Border     = lipgloss.Color("#3B3551")
	Selection  = lipgloss.Color("#3B3551")

	// Brand colors that NEVER change with themes
	BrandCyan = lipgloss.Color("#00E5FF")

	// ── Typography ──────────────────────────────────────────────────
	BaseStyle       lipgloss.Style
	TitleStyle      lipgloss.Style
	SubtitleStyle   lipgloss.Style
	LabelStyle      lipgloss.Style
	MutedStyle      lipgloss.Style
	MutedLightStyle lipgloss.Style
	PromptStyle     lipgloss.Style
	SuccessStyle    lipgloss.Style
	ErrorStyle      lipgloss.Style
	DividerStyle    lipgloss.Style

	// ── Chat bubbles ────────────────────────────────────────────────
	UserLabelStyle       lipgloss.Style
	UserBubbleStyle      lipgloss.Style
	AssistantLabelStyle  lipgloss.Style
	AssistantBubbleStyle lipgloss.Style

	// ── Input area ──────────────────────────────────────────────────
	InputBoxStyle      lipgloss.Style
	InputBoxFocusStyle lipgloss.Style
)

func init() {
	UpdateStyles()
}

// UpdateStyles re-evaluates all lipgloss styles based on the current color variables.
func UpdateStyles() {
	BaseStyle = lipgloss.NewStyle().Background(Surface)

	TitleStyle = BaseStyle.
		Bold(true).
		Foreground(Purple)

	SubtitleStyle = BaseStyle.
		Bold(true).
		Foreground(Green)

	LabelStyle = BaseStyle.
		Bold(true).
		Foreground(Cyan)

	MutedStyle = BaseStyle.
		Foreground(Muted)

	MutedLightStyle = BaseStyle.
		Foreground(MutedLight)

	PromptStyle = BaseStyle.
		Bold(true).
		Foreground(Amber)

	SuccessStyle = BaseStyle.
		Foreground(Green)

	ErrorStyle = BaseStyle.
		Bold(true).
		Foreground(Red)

	DividerStyle = BaseStyle.
		Foreground(Border)

	UserLabelStyle = BaseStyle.
		Bold(true).
		Foreground(Amber)

	UserBubbleStyle = BaseStyle.
		Foreground(White).
		PaddingLeft(2)

	AssistantLabelStyle = BaseStyle.
		Bold(true).
		Foreground(BrandCyan)

	AssistantBubbleStyle = BaseStyle.
		Foreground(White).
		PaddingLeft(2)

	InputBoxStyle = BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Padding(0, 1)

	InputBoxFocusStyle = BaseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Purple).
		Padding(0, 1)
}
