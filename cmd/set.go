package cmd

import (
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/ui/selector"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "Update a configuration value (or select interactively)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		key := args[0]
		var value string
		if len(args) == 2 {
			value = args[1]
		}

		switch key {
		case "provider":
			if value == "" {
				opt, err := selector.Run(
					"Select AI Provider",
					"Choose which AI provider you want WindMist to use:",
					config.GetProviderOptions(),
				)
				if err != nil {
					log.Fatal(err)
				}
				value = opt.Value
				err = cfg.SetProvider(value)
				if err != nil {
					log.Fatal(err)
				}

				// Automatically prompt for model right after selecting provider
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
				err = cfg.SetModel(value, modelValue)
				if err != nil {
					log.Fatal(err)
				}
				if err := config.Save(cfg); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("✔ Provider set to %s (model: %s)\n", value, modelValue)
				return
			}
			err = cfg.SetProvider(value)

		case "model":
			if value == "" {
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
			err = cfg.SetModel(cfg.AI.Provider, value)

		case "api_key":
			if value == "" {
				customVal, err := selector.RunInput(fmt.Sprintf("API Key for %s", cfg.AI.Provider), "Enter API Key", "")
				if err != nil {
					log.Fatal(err)
				}
				value = customVal
			}
			err = cfg.SetAPIKey(cfg.AI.Provider, value)

		case "base_url":
			log.Fatal("base_url cannot be set or changed by the user for any provider")

		case "theme":
			if value == "" {
				opt, err := selector.Run("Select UI Theme", "Choose visual theme:", []selector.Option{
					{Label: "dark", Description: "Dark theme with purple & cyan accents", Value: "dark"},
					{Label: "light", Description: "Light theme", Value: "light"},
				})
				if err != nil {
					log.Fatal(err)
				}
				value = opt.Value
			}
			cfg.SetTheme(value)

		default:
			log.Fatalf("unknown configuration key: %s", key)
		}

		if err != nil {
			log.Fatal(err)
		}

		if err := config.Save(cfg); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configuration updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
