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

		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		provider, err := ai.New(cfg)
		if err != nil {
			log.Fatal(err)
		}

		manager := tools.NewManager()
		defaults.RegisterAll(manager)

		ag := agent.New(provider, manager, agent.Config{})

		res, err := ag.Run(context.Background(), args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Content)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
