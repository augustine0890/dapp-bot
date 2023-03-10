package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/augustine0890/dapp-bot/internal/database"
	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/logging"
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
func NewDiscord(cfg *config.Config) (*Discord, error) {
	// Connect to MongoDB
	mongoClient, err := database.GetMongoClient(cfg.MongoURI, cfg.MongoDBName, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to create MongoDB client: %w", err)
	}

	// Create a new Discord session
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	intents := discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMembers
	session.Identify.Intents = intents

	// Create a new CommandHandler and register commands
	ch := NewCommandHandler()
	ch.RegisterCommand(HandlePing, "ping")
	ch.RegisterCommand(HandlePlayDapp, "dapp")
	ch.RegisterCommand(CheckPointCommand(cfg, mongoClient), "cp", "checkpoint")
	ch.RegisterCommand(RankCommand(cfg, mongoClient), "r", "rank")
	ch.RegisterCommand(MyRankCommand(cfg, mongoClient), "mr", "myrank")

	// Register the command handler function
	session.AddHandler(ch.HandleCommand)

	// Register the member join/leave handler function
	RegisterHandler(session, mongoClient, cfg, &discordgo.GuildMemberAdd{}, HandleMember)
	RegisterHandler(session, mongoClient, cfg, &discordgo.GuildMemberRemove{}, HandleMember)

	// Register additional event handlers here as needed

	session.AddHandler(HandleRemoveReaction)

	// Create a new Discord instance
	d := &Discord{
		session:       session,
		commandPrefix: "!",
		mongoClient:   mongoClient,
		reactionCh:    make(chan *discordgo.MessageReactionAdd),
	}

	// Register the ready command handler function
	// session.AddHandler(HandleReady)

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

	log.Println("Bot is now running. Press CTRL-C to exit.")

	// Wait for CTRL-C or SIGINT/SIGTERM
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Cleanly close down the Discord session
	err = d.session.Close()
	if err != nil {
		logging.Error("Error closing Discord session:", err)
	}

	return nil
}
