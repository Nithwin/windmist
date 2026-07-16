package cmd

import (
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Select or switch the active AI model interactively",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		var value string
		if len(args) == 1 {
			value = args[0]
		} else {
			ollamaBaseURL := ""
			if pConfig, ok := cfg.Providers["ollama"]; ok {
				ollamaBaseURL = pConfig.BaseURL
			}
			opt, err := selector.Run(
				fmt.Sprintf("Select Model for %s", cfg.AI.Provider),
				"Choose the active model to use:",
				config.GetModelOptions(cfg.AI.Provider, ollamaBaseURL),
			)
			if err != nil {
				log.Fatal(err)
			}
			value = opt.Value
			if value == "__CUSTOM__" {
				customVal, err := selector.RunInput("Custom Model ID", "Enter exact model ID (e.g. gpt-4o)", "")
				if err != nil {
					log.Fatal(err)
				}
				value = customVal
			}
		}

		if err := cfg.SetModel(cfg.AI.Provider, value); err != nil {
			log.Fatal(err)
		}

		if err := config.Save(cfg); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("✔ Model set to %s (provider: %s)\n", value, cfg.AI.Provider)
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
