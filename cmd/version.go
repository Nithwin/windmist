package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the WindMist version",
	Long:  "Display the current version of the WindMist CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WindMist %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
