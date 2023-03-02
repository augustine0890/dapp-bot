package discord

import (
	"context"
	"time"

	"github.com/augustine0890/dapp-bot/internal/database"
	"github.com/augustine0890/dapp-bot/pkg/config"
	"github.com/augustine0890/dapp-bot/pkg/logging"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleMember(s *discordgo.Session, e interface{}, mongoClient *mongo.Client, cfg *config.Config) {
	var userID string
	var username string
	var joinedDate time.Time
	var leave bool

	switch event := e.(type) {
	case *discordgo.GuildMemberAdd:
		userID = event.Member.User.ID
		username = event.Member.User.Username
		joinedDate = event.Member.JoinedAt
		leave = false
	case *discordgo.GuildMemberRemove:
		userID = event.User.ID
		username = event.User.Username
		leave = true
	default:
		return
	}

	usersColl := database.GetUsersColl(mongoClient, cfg)
	ctx := context.TODO()

	if leave {
		// Delete user from users collection
		filter := bson.M{"_id": userID}
		_, err := usersColl.DeleteOne(ctx, filter)
		if err != nil {
			logging.Warn("Failed to delete user document", err)
		}

		// Delete user's activities from activities collection
		activitiesColl := database.GetActivitiesColl(mongoClient, cfg)
		_, err = activitiesColl.DeleteMany(ctx, bson.M{"user": userID})
		if err != nil {
			logging.Warn("Failed to delete activity documents for user", err)
		}

	} else {
		user := &database.User{
			ID:         userID,
			UserName:   username,
			Points:     0,
			JoinedDate: joinedDate,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		_, err := usersColl.InsertOne(ctx, user)
		if err != nil {
			logging.Warn("Failed to insert user document", err)
		}
	}
}
