package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Nithwin/WindMist/internal/ai"
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:  "chat <prompt>",
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		provider, err := ai.New(cfg)
		if err != nil {
			log.Fatal(err)
		}

		req := &ai.GenerateRequest{
			Messages: []ai.Message{
				{
					Role:    ai.RoleUser,
					Content: args[0],
				},
			},
		}

		resp, err := provider.Generate(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp.Text)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
