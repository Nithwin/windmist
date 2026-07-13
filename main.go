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

	// cfg, err := config.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Provider :", cfg.AI.Provider)
	// log.Println("Model    :", cfg.AI.Model)
	// log.Println("Theme    :", cfg.UI.Theme)
	cmd.Execute()
}
