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

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]
		value := args[1]

		if err := config.Update(key, value); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configuration updated successfully.")
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}