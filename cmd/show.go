package cmd

import (
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current configuration",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		model, err := cfg.ActiveModel()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("WindMist Configuration")
		fmt.Println("----------------------")
		fmt.Printf("Provider : %s\n", cfg.AI.Provider)
		fmt.Printf("Model    : %s\n", model)
		fmt.Printf("Theme    : %s\n", cfg.UI.Theme)
		fmt.Printf("Cache    : %t\n", cfg.Cache.Enabled)
	},
}

func init() {
	configCmd.AddCommand(showCmd)
}
