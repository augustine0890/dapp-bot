package main

import (
	"fmt"

	"github.com/augustine0890/dapp-bot/pkg/config"
)

func main() {
	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %s\n", err)
	}

	fmt.Println(cfg)
}
