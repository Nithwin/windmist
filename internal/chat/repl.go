package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui"
)

// Start launches the interactive chat session.
func Start() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	provider, err := ai.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	PrintBanner(cfg)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(ui.PromptStyle.Render("❯ You > "))

		if !scanner.Scan() {
			fmt.Println()
			return
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		switch input {
		case "/exit", "/quit":
			fmt.Println()
			fmt.Println(ui.SuccessStyle.Render("👋 Thanks for using WindMist. See you soon!"))
			fmt.Println()
			return
		}

		req := &ai.GenerateRequest{
			Messages: []ai.Message{
				{
					Role:    ai.RoleUser,
					Content: input,
				},
			},
		}

		resp, err := provider.Generate(context.Background(), req)
		if err != nil {
			fmt.Println(ui.ErrorStyle.Render("✗ " + err.Error()))
			continue
		}

		fmt.Println()
		fmt.Println(ui.SubtitleStyle.Render("WindMist"))
		fmt.Println(resp.Text)
		fmt.Println()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
