package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

// Discord represents a Discord instance for the bot
type Discord struct {
	session       *discordgo.Session
	commandPrefix string
	mongoClient   *mongo.Client
	reactionCh    chan *discordgo.MessageReactionAdd
}

// NewDiscord creates a new Discord instance for the bot
func NewDiscord(token string, mongoClient *mongo.Client) (*Discord, error) {
	// Create a new Discord session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Create a new CommandHandler and register commands
	ch := NewCommandHandler()
	ch.RegisterCommand("ping", HandlePing)
	ch.RegisterCommand("dapp", HandlePlayDapp)
	ch.RegisterCommand("rank", HandleRank)

	// Register the command handler function
	session.AddHandler(ch.HandleCommand)

	// Create a new Discord instance
	d := &Discord{
		session:       session,
		commandPrefix: "!",
		mongoClient:   mongoClient,
		reactionCh:    make(chan *discordgo.MessageReactionAdd),
	}

	// Register the message reaction handler function
	// session.AddHandler(d.HandleReaction)

	return d, nil
}

// Connect connects the Discord instance to the server
func (d *Discord) Connect() error {
	err := d.session.Open()
	if err != nil {
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// Wait for CTRL-C or SIGINT/SIGTERM
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Cleanly close down the Discord session
	err = d.session.Close()
	if err != nil {
		fmt.Println("Error closing Discord session:", err)
	}

	return nil
}
