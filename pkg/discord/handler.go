package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandlerFunc represents a function that handles a Discord command
type CommandHandlerFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

// CommandHandler represents a handler for Discord commands
type CommandHandler struct {
	commands map[string]CommandHandlerFunc
}

// NewCommandHandler creates a new CommandHandler instance
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: make(map[string]CommandHandlerFunc),
	}
}

// RegisterCommand registers a command with the CommandHandler
func (ch *CommandHandler) RegisterCommand(name string, handler CommandHandlerFunc) {
	ch.commands[name] = handler
}

// HandleCommand handles incoming messages and triggers the corresponding command handlers
func (ch *CommandHandler) HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Split the message content into command and arguments
	parts := strings.Split(strings.TrimSpace(m.Content), " ")
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "!") {
		// The message doesn't start with a command prefix
		return
	}
	// Look up the handler function for the command
	command := strings.TrimPrefix(parts[0], "!")
	handler, ok := ch.commands[command]
	if !ok {
		// Unknown command
		return
	}

	// Call the handler function
	args := parts[1:]
	handler(s, m, args)
}

// HandlePing handles the !ping command and sends a response message
func HandlePing(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		fmt.Println("Error handling !ping command:", err)
	}
}

// HandlePlayDapp handles the !dapp command and sends a response message
func HandlePlayDapp(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Let's play DappBot!")
	if err != nil {
		fmt.Println("Error handling !dapp command:", err)
	}
}

// HandleRank handles the !rank command and sends the user's rank card as a message
func HandleRank(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// TODO: Implement rank card generation logic
	// Send the generated rank card as a message
	_, err := s.ChannelMessageSend(m.ChannelID, "Your rank card:")
	if err != nil {
		fmt.Println("Error handling !rank command:", err)
	}
}
