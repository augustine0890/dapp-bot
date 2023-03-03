package discord

import (
	"context"
	"fmt"
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

		// Send a DM to the new member with the long text
		dm, err := s.UserChannelCreate(userID)
		if err != nil {
			logging.Error("Failed to create DM channel", err)
			return
		}
		message := fmt.Sprintf(`Dear %s,

		:partying_face: Welcome to the PlayDapp Discord family!! Feel free to chat with other PLAyers here and enjoy the Web3-based games! :hugging_face:

		**About Marketplace** :circus_tent:
		Marketplace is an exchange for Web3 game items. You can buy and sell these game items in the Web3 marketplace.
		FAQ for Marketplace: https://market.playdapp.com/faq/marketplace

		**About Item Manager** :joystick:
		The Item Manager is a service that enables you to convert game items in AWTG to be kept or archived, this process also allows the item to be used or traded on the PlayDapp Web3 Marketplace. FAQ for Item Manager: https://itemmanager.playdapp.com/faq

		**About Along with the Gods** :dragon_face:
		Along with the Gods (AWTG) is a play-to-earn mobile RPG and part of PlayDapp’s multi-homing game strategy, in which players can easily move across and use various platforms.

		**About Tournament** :space_invader:
		PlayDapp Tournaments is a PvP P2E Hypercasual Game platform. After the tournament ends, rewards are awarded based on the player’s ranking on the leaderboard. We plan on adding 3-4 games every quarter for a total of 40 games!
		FAQ for Tournament: https://tournament.playdapp.com/faq

		**Can’t receive your rewards?** :gift:
		Please submit a CS ticket here: https://playdapp.atlassian.net/servicedesk/customer/portals

		**Any bugs / technical problems?** :sob:
		Please submit a ticket here: https://dashboard-api.gamepot.ntruss.com/v2/cs/request?projectId=6cec6ce7-436c-4dd7-a558-b2ef6ee767dc&language=en&mode=

		**Social Media** :computer:
		*Twitter*: https://twitter.com/playdapp_io
		*Medium*: https://medium.com/playdappgames`, username)

		_, err = s.ChannelMessageSend(dm.ID, message)
		if err != nil {
			logging.Warn("Failed to send DM message", err)
		}
	}
}
