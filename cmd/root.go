package cmd

import (
	"os"

	_ "github.com/Nithwin/WindMist/internal/providers/gemini"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "windmist",
	Short: "AI Software Engineer in your terminal",
	Long: `WindMist is an AI-powered developer assistant built for the terminal.

It helps you understand repositories, explain code, generate features,
review changes, automate development workflows, and interact with
multiple AI providers from a single CLI.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags will be added here later.
}
