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

func HandleReady(s *discordgo.Session, event *discordgo.Ready) {
	targetGuild := "1019782712799805440"

	// Get the ID of the first server that the bot is connected to.
	guilds := s.State.Guilds
	if len(guilds) == 0 {
		fmt.Println("Error: Bot is not connected to any servers.")
		return
	}

	guildID := guilds[0].ID
	if guildID != targetGuild {
		return
	}

	// Fetch all the members of the server with the specified ID.
	const batchSize = 1000
	var members []*discordgo.Member
	var lastMemberID string
	for {
		fetchedMembers, err := s.GuildMembers(guildID, lastMemberID, batchSize)
		if err != nil {
			fmt.Println("Error: Could not fetch members", err)
			return
		}
		members = append(members, fetchedMembers...)

		if len(fetchedMembers) < batchSize {
			break
		}
		lastMemberID = fetchedMembers[len(fetchedMembers)-1].User.ID
	}

	fmt.Printf("Found %d members in server %s:\n", len(members), guilds[0].Name)
	for _, member := range members {
		// Skip this member if they are a bot.
		if member.User.Bot {
			fmt.Println("This member is a bot.")
			continue
		}
		// fmt.Printf("%s#%s (ID: %s)\n", member.User.Username, member.User.Discriminator, member.User.ID)
		// fmt.Printf("Joined at: %s\n", member.JoinedAt.Format("2006-01-02 15:04:05"))
		// fmt.Printf("Roles: %s\n", strings.Join(member.Roles, ", "))
		// fmt.Println()
	}
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
