package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the WindMist version",
	Long:  "Display the current version of the WindMist CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		if Commit != "none" {
			fmt.Printf("WindMist %s (commit %s, built %s)\n", Version, Commit, Date)
		} else {
			fmt.Printf("WindMist %s\n", Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
