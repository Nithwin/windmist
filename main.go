package main

import (
	"log"

	"github.com/Nithwin/WindMist/cmd"
	"github.com/Nithwin/WindMist/internal/config"
)

func main() {
	// Load local settings before the CLI starts.
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Hand off to the Cobra command tree.
	cmd.Execute()
}
