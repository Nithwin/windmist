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

var (
	flagModel    string
	flagProvider string
)

var chatCmd = &cobra.Command{
	Use:  "chat <prompt>",
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		if flagProvider != "" {
			if err := cfg.SetProvider(flagProvider); err != nil {
				log.Fatal(err)
			}
		}

		if flagModel != "" {
			if err := cfg.SetModel(cfg.AI.Provider, flagModel); err != nil {
				log.Fatal(err)
			}
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
	chatCmd.Flags().StringVarP(&flagModel, "model", "m", "", "AI model to use (e.g. gpt-4o, claude-3-5-sonnet-latest, qwen2.5:8b)")
	chatCmd.Flags().StringVarP(&flagProvider, "provider", "p", "", "AI provider to use (e.g. gemini, ollama, groq, openai, anthropic)")
	rootCmd.AddCommand(chatCmd)
}
