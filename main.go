package main

import (
	"log"

	"github.com/Nithwin/WindMist/cmd"
	"github.com/Nithwin/WindMist/internal/config"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}
	cmd.Execute()
}
