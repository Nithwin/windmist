package cmd

import (
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Update a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		key := args[0]
		value := args[1]

		switch key {
		case "provider":
			err = cfg.SetProvider(value)

		case "model":
			err = cfg.SetModel(cfg.AI.Provider, value)

		case "api_key":
			err = cfg.SetAPIKey(cfg.AI.Provider, value)

		case "base_url":
			err = cfg.SetBaseURL(cfg.AI.Provider, value)

		case "theme":
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