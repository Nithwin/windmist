package chat

import (
	"fmt"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui"
)

// PrintBanner displays the WindMist welcome banner.
func PrintBanner(cfg *config.Config) {
	logo := `
██╗    ██╗██╗███╗   ██╗██████╗ ███╗   ███╗██╗███████╗████████╗
██║    ██║██║████╗  ██║██╔══██╗████╗ ████║██║██╔════╝╚══██╔══╝
██║ █╗ ██║██║██╔██╗ ██║██║  ██║██╔████╔██║██║███████╗   ██║
██║███╗██║██║██║╚██╗██║██║  ██║██║╚██╔╝██║██║╚════██║   ██║
╚███╔███╔╝██║██║ ╚████║██████╔╝██║ ╚═╝ ██║██║███████║   ██║
 ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═════╝ ╚═╝     ╚═╝╚══════╝   ╚═╝
`

	fmt.Println(ui.TitleStyle.Render(logo))
	fmt.Println(ui.SubtitleStyle.Render("⚡ WindMist - AI Coding Assistant"))
	fmt.Println()

	fmt.Print(ui.LabelStyle.Render("Provider : "))
	fmt.Println(cfg.AI.Provider)

	if provider, err := cfg.ActiveProvider(); err == nil {
		fmt.Print(ui.LabelStyle.Render("Model    : "))
		fmt.Println(provider.Model)
	}

	fmt.Println()

	fmt.Println(ui.DividerStyle.Render("────────────────────────────────────────────────────────────"))

	fmt.Println(ui.SuccessStyle.Render("Type /help for commands"))
	fmt.Println(ui.SuccessStyle.Render("Type /exit to quit"))

	fmt.Println(ui.DividerStyle.Render("────────────────────────────────────────────────────────────"))
	fmt.Println()
}