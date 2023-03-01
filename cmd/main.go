package main

import (
	"fmt"
	"log"

	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/discord"
)

func main() {
	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
	}

	// Create a new Discord bot instance
	dc, err := discord.NewDiscord(cfg)
	if err != nil {
		log.Fatalf("Failed to create Discord bot instance: %v\n", err)
	}

	// Connect to bot to the server
	err = dc.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Discord server: %v\n", err)
	}
}
