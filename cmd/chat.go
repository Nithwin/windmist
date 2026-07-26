package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/agent"
	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/defaults"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:  "chat <prompt>",
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		// Load the saved settings before starting the assistant.
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		// Create the configured AI provider.
		provider, err := ai.New(cfg)
		if err != nil {
			log.Fatal(err)
		}

		// Register only the basic tools used in this beginner version.
		manager := tools.NewManager()
		defaults.RegisterAll(manager, nil, cfg)

		// Run one prompt and print the final answer.
		ag := agent.New(provider, manager, agent.Config{})
		res, err := ag.Run(context.Background(), nil, args[0], nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Content)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
