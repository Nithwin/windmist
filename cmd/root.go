package cmd

import (
	_ "github.com/Nithwin/WindMist/internal/providers/gemini"

	"github.com/Nithwin/WindMist/internal/chat"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "windmist",
	Short: "WindMist - AI Coding Assistant",
	Long:  "WindMist is an AI-powered coding assistant for your terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return chat.Run()
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}