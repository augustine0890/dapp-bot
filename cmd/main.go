package main

import (
	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/discord"
	"github.com/augustine0890/dapp-bot/pkg/logging"
)

func main() {
	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Error("Failed to read config file:", err)
	}

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
