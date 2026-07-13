package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage WindMist configuration",
	Long:  "Manage and inspect the WindMist configuration.",
}

func init() {
	rootCmd.AddCommand(configCmd)
}