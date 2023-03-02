package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/augustine0890/dapp-bot/internal/database"
	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckPointCommand returns a command handler function for the !checkpoint command
func CheckPointCommand(cfg *config.Config, mongoClient *mongo.Client) CommandHandlerFunc {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		handleCheckPoint(s, m, args, cfg, mongoClient)
	}
}

// HandleCheckPoint handles the !checkpoint command, sending the user's points as an embed message
func handleCheckPoint(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cfg *config.Config, mongoClient *mongo.Client) {
	// Retrieve the user's points from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersColl := database.GetUsersColl(mongoClient, cfg)

	attendanceChannelID := cfg.AttendanceID
	if m.ChannelID != attendanceChannelID {
		message := fmt.Sprintf("<@%s> Please go to the <#%s> channel for Daily Attendance and Points Checking.", m.Author.ID, attendanceChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
		}
		return
	}

	userID := m.Author.ID
	filter := bson.M{"_id": userID}
	var user database.User
	err := usersColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("Error retrieving user points:", err)
		return
	}

	// Create an embed massage with the user's points
	embed := &discordgo.MessageEmbed{
		Title: "The Cumulative Points",
		Author: &discordgo.MessageEmbedAuthor{
			Name: m.Author.Username,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Given to %s", m.Author.Username),
			IconURL: m.Author.AvatarURL(""),
		},
		Color:     0x00aaff,
		Timestamp: user.UpdatedAt.Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Points",
				Value:  fmt.Sprintf("%d", user.Points),
				Inline: true,
			},
		},
	}
	// Send the embed message as a reply to the original message
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
