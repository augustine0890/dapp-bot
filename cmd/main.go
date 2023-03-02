package main

import (
	"flag"
	"log"

	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/discord"
	"github.com/augustine0890/dapp-bot/pkg/logging"
)

func main() {
	// Flag will be stored in the stage variable at runtime
	stage := flag.String("stage", "prod", "The enviroment running")
	flag.Parse()

	// Load the application configuration
	cfg, err := config.LoadConfig(*stage)
	if err != nil {
		logging.Error("Failed to read config file:", err)
	}
	log.Printf("Config loaded and running with %s stage.", *stage)

	// Create a new Discord bot instance
	dc, err := discord.NewDiscord(cfg)
	if err != nil {
		logging.Fatal("Failed to create Discord bot instance:", err)
	}

	// Connect to bot to the server
	err = dc.Connect()
	if err != nil {
		logging.Fatal("Failed to connect to Discord server:", err)
	}
}
