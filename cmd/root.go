package cmd

import (
	_ "github.com/Nithwin/WindMist/internal/providers/gemini"

	"github.com/Nithwin/WindMist/internal/chat"
	"github.com/spf13/cobra"
)

// rootCmd keeps the app focused on the simplest terminal chat flow.
var rootCmd = &cobra.Command{
	Use:   "windmist",
	Short: "WindMist - simple terminal AI assistant",
	Long:  "WindMist is a simple AI assistant for your terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return chat.Run()
	},
}

// Execute starts the root command.
func Execute() error {
	return rootCmd.Execute()
}
