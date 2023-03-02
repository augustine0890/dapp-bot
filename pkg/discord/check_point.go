package discord

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/augustine0890/dapp-bot/internal/database"
	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/logging"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckPointCommand returns a command handler function for the !checkpoint command
func CheckPointCommand(cfg *config.Config, mongoClient *mongo.Client) CommandHandlerFunc {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		handleCheckPoint(s, m, args, cfg, mongoClient)
	}
}

func RankCommand(cfg *config.Config, mongoClient *mongo.Client) CommandHandlerFunc {
	return func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		handleRank(s, m, args, cfg, mongoClient)
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
			logging.Error("Error sending message", err)
		}
		return
	}

	userID := m.Author.ID
	filter := bson.M{"_id": userID}
	var user database.User
	err := usersColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logging.Error("Failed retrieving user points", err)
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

func handleRank(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cfg *config.Config, mongoClient *mongo.Client) {
	// Get the users collection from MongoDB
	// Retrieve the user's points from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersColl := database.GetUsersColl(mongoClient, cfg)

	// Find the top 10 users based on their points
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"points", -1}, {"updatedAt", 1}})
	findOptions.SetLimit(10)
	cursor, err := usersColl.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		msg := "Error fetching user data from the database"
		logging.Error(msg, err)
		return
	}
	defer cursor.Close(ctx)

	// Create the fields for each user in the top 10
	// Build the list of rank fields
	topRank := make([]*discordgo.MessageEmbedField, 0, 10)
	emojiRank := []string{"🥇", "🥈", "🥉", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
	i := 0
	for cursor.Next(ctx) {
		// Decode the user document
		var rankUser database.User
		err := cursor.Decode(&rankUser)
		if err != nil {
			logging.Error("Failed to get rank.", err)
			return
		}

		// Add the user to the list of rank fields
		topRank = append(topRank, &discordgo.MessageEmbedField{
			Name:  emojiRank[i] + " " + rankUser.UserName,
			Value: strconv.Itoa(rankUser.Points) + " 🧧",
		})
		i++
	}

	if err := cursor.Err(); err != nil {
		logging.Error("Failed to get rank.", err)
		return
	}

	// Open the image file
	imageFile, err := os.Open("./assets/images/winners.jpg")
	if err != nil {
		logging.Error("Failed to open image file:", err)
		return
	}
	defer imageFile.Close()
	// Get the file info
	fileInfo, err := imageFile.Stat()
	if err != nil {
		logging.Error("Failed to get image file info:", err)
		return
	}

	// Read the image file into a byte slice
	imageData := make([]byte, fileInfo.Size())
	_, err = imageFile.Read(imageData)
	if err != nil {
		logging.Error("Failed to read image file:", err)
		return
	}
	// Create a new discordgo file from the byte slice
	image := &discordgo.File{Name: fileInfo.Name(), Reader: bytes.NewReader(imageData)}
	// Create a new embed message
	embed := &discordgo.MessageEmbed{
		Title:       "🏆 The Cumulative Points TOP 10 Leaderboard 🏆",
		Description: "Congratulations! You made it! 🥳",
		Color:       0x00AAFF,
		Fields:      topRank,
		Image: &discordgo.MessageEmbedImage{
			URL:      "attachment://" + image.Name,
			Width:    400,
			ProxyURL: "",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Made by %s", s.State.User.Username),
			IconURL: s.State.User.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Send the message with the image
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: embed,
		Files: []*discordgo.File{image},
	})
	if err != nil {
		logging.Error("Error sending message to channel.", err)
		return
	}
}
