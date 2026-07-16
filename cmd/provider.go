package cmd

import (
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider [name]",
	Short: "Select or switch the active AI provider interactively",
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
			opt, err := selector.Run(
				"Select AI Provider",
				"Choose which AI provider you want WindMist to use:",
				config.GetProviderOptions(),
			)
			if err != nil {
				log.Fatal(err)
			}
			value = opt.Value
		}

		if err := cfg.SetProvider(value); err != nil {
			log.Fatal(err)
		}

		ollamaBaseURL := ""
		if pConfig, ok := cfg.Providers["ollama"]; ok {
			ollamaBaseURL = pConfig.BaseURL
		}
		modelOpt, err := selector.Run(
			fmt.Sprintf("Select Model for %s", value),
			"Choose the active model for this provider:",
			config.GetModelOptions(value, ollamaBaseURL),
		)
		if err != nil {
			log.Fatal(err)
		}
		modelValue := modelOpt.Value
		if modelValue == "__CUSTOM__" {
			customVal, err := selector.RunInput("Custom Model ID", "Enter exact model ID (e.g. gpt-4o)", "")
			if err != nil {
				log.Fatal(err)
			}
			modelValue = customVal
		}

		if err := cfg.SetModel(value, modelValue); err != nil {
			log.Fatal(err)
		}

		if err := config.Save(cfg); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("✔ Provider set to %s (model: %s)\n", value, modelValue)
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)
}
